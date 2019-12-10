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

func (ps *PostService) SelectPostById(id int) (postSl []input.Post, er error) {
	sqlQuery := `SELECT p.author, p.created, p.forum, p.id, p.isEdited, p.message, p.parent, p.thread
	FROM public.post as p
	where p.id = $1`
	postSlice := make([]input.Post, 0)
	rows, err := ps.db.Query(sqlQuery, id)
	if err != nil {
		return postSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		post := input.Post{}
		err := rows.Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		if err != nil {
			return postSlice, err
		}
		postSlice = append(postSlice, post)
	}
	return postSlice, nil
}
