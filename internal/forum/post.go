package forum

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type PostService struct {
	db *pgx.ConnPool
}

func NewPostService(db *pgx.ConnPool) *PostService {
	return &PostService{db: db}
}

func (ps *PostService) SelectPostById(id int) (post Post, err error) {
	sqlQuery := `SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread
	FROM public.post as p
	where p.id = $1`
	err = ps.db.QueryRow(sqlQuery, id).Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	return
}

func (ps *PostService) FindPostById(id int, thread int) (err error) {
	sqlQuery := `SELECT p.id
	FROM public.post as p
	where p.id = $1 AND p.thread = $2`
	var postId int64
	err = ps.db.QueryRow(sqlQuery, id, thread).Scan(&postId)
	return
}

func (ps *PostService) InsertPost(post Post) (lastId int, err error) {
	sqlQuery := `INSERT INTO public.post (author, created, forum, message, parent, thread)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`
	err = ps.db.QueryRow(sqlQuery, post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread).Scan(&lastId)
	return
}

func (ps *PostService) UpdatePostMessage(newMessage string, id int) (countUpdateString int64, err error) {
	sqlQuery := `UPDATE public.post SET message = $1,
                       is_edited = true
	where post.id = $2`
	result, err := ps.db.Exec(sqlQuery, newMessage, id)
	if err != nil {
		return
	}
	countUpdateString = result.RowsAffected()
	return
}

func (ps *PostService) CreatePosts(thread Thread, forumId int, created string, posts []Post) (post []Post, err error) {
	tx, err := ps.db.Begin()
	if err != nil {
		return nil, err
	}

	//now := time.Now()

	sqlStr := "INSERT INTO post(id, parent, thread, forum, author, created, message, path) VALUES "
	vals := []interface{}{}
	for _, post := range posts {
		var authorId int
		err = ps.db.QueryRow(`SELECT id FROM public."user" WHERE LOWER(nick_name) = LOWER($1)`,
			post.Author,
		).Scan(&authorId)
		if err != nil {
			_ = tx.Rollback()
			return nil, errors.New("404")
		}
		sqlQuery := `
		INSERT INTO public.forum_user (forum_id, user_id)
		VALUES ($1,$2)`
		_, err = ps.db.Exec(sqlQuery, forumId, authorId)

		if post.Parent == 0 {
			sqlStr += "(nextval('post_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"ARRAY[currval(pg_get_serial_sequence('post', 'id'))::bigint]),"
			vals = append(vals, post.Parent, thread.Id, thread.Forum, post.Author, created, post.Message)
		} else {
			var parentThreadId int32
			err = ps.db.QueryRow("SELECT post.thread FROM public.post WHERE post.id = $1",
				post.Parent,
			).Scan(&parentThreadId)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			if parentThreadId != int32(thread.Id) {
				_ = tx.Rollback()
				return nil, errors.New("Parent post was created in another thread")
			}

			sqlStr += " (nextval('post_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"(SELECT post.path FROM public.post WHERE post.id = ? AND post.thread = ?) || " +
				"currval(pg_get_serial_sequence('post', 'id'))::bigint),"

			vals = append(vals, post.Parent, thread.Id, thread.Forum, post.Author, created, post.Message, post.Parent, thread.Id)
		}

	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr += " RETURNING  id, parent, thread, forum, author, created, message, is_edited "

	sqlStr = ReplaceSQL(sqlStr, "?")
	if len(posts) > 0 {
		rows, err := tx.Query(sqlStr, vals...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(
				&(posts)[i].Id,
				&(posts)[i].Parent,
				&(posts)[i].Thread,
				&(posts)[i].Forum,
				&(posts)[i].Author,
				&(posts)[i].Created,
				&(posts)[i].Message,
				&(posts)[i].IsEdited,
			)
			i += 1

			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}


/*func (ps *PostService) CheckPosts(threadId int, posts []Post) (err error) {
	_, err := h.UserService.FindUserByNickName(newPosts[i].Author)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can'ps find user"})
	}
	if newPosts[i].Parent != 0 {
		err = h.PostService.FindPostById(newPosts[i].Parent, newPosts[i].Thread)
		if err != nil {
			return ctx.JSON(http.StatusConflict, forum.ErrorMessage{Message: "Can'ps find post"})
		}
	}
}*/
