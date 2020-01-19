package forum

import (
	"fmt"
	"github.com/jackc/pgx"
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

func (ps *PostService) CreatePosts(threadId int, postForum, created string, posts []Post) (ids []int, err error) {
	columns := 6
	placeholders := make([]string, 0, len(posts))
	args := make([]interface{}, 0, len(posts)*columns)
	for i, post := range posts {
		args = append(args, threadId, postForum, post.Parent, post.Author, post.Message, created)
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d)",
			i*columns+1, i*columns+2, i*columns+3, i*columns+4, i*columns+5, i*columns+6,
		))
	}
	query := fmt.Sprintf(
		"insert into public.post (thread, forum, parent, author, message, created) values %s",
		strings.Join(placeholders, ","),
	)
	query = query + " RETURNING id"
	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		scanId := 0
		//var timetz time.Time
		err = rows.Scan(&scanId)
		if err != nil {
			return
		}
		//scanPost.Created = timetz.Format(time.RFC3339Nano)
		ids = append(ids, scanId)
	}
	return
}

/*func (ps *PostService) CheckPosts(threadId int, posts []Post) (err error) {
	_, err := h.UserService.FindUserByNickName(newPosts[i].Author)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
	}
	if newPosts[i].Parent != 0 {
		err = h.PostService.FindPostById(newPosts[i].Parent, newPosts[i].Thread)
		if err != nil {
			return ctx.JSON(http.StatusConflict, forum.ErrorMessage{Message: "Can't find post"})
		}
	}
}*/
