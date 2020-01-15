package forum

import (
	"database/sql"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
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
	countUpdateString, err = result.RowsAffected()
	return
}
