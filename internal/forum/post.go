package forum

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/input"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{db: db}
}

func (ps *PostService) SelectPostById(id int) (post input.Post, err error) {
	sqlQuery := `SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread
	FROM public.post as p
	where p.id = $1`
	err = ps.db.QueryRow(sqlQuery, id).Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	return
}

func (ps *PostService) UpdatePostMessage(newMessage string, id int) (countUpdateString int64, err error) {
	sqlQuery := `UPDATE public.post SET message = $1
	where post.id = $2`
	result, err := ps.db.Exec(sqlQuery, id)
	if err != nil {
		return
	}
	countUpdateString, err = result.RowsAffected()
	return
}
