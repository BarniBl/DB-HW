package forum

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/input"
)

type ForumService struct {
	db *sql.DB
}

func NewForumService(db *sql.DB) *ForumService {
	return &ForumService{db: db}
}

func (fs *ForumService) SelectForumBySlug(slug string) (forumSl []input.Forum, er error) {
	sqlQuery := `SELECT f.slug, f.title, f.user
	FROM public.forum as f 
	where f.slug = $1`
	forumSlice := make([]input.Forum, 0)
	rows, err := fs.db.Query(sqlQuery, slug)
	if err != nil {
		return forumSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		forum := input.Forum{}
		err := rows.Scan(&forum.Slug, &forum.Title, &forum.User)
		if err != nil {
			return forumSlice, err
		}
		forumSlice = append(forumSlice, forum)
	}
	return forumSlice, nil
}

func (fs *ForumService) InsertForum(forum input.Forum) error {
	sqlQuery := `INSERT INTO public.forum (slug, title, "user")
	VALUES ($1,$2,$3)`
	_, err := fs.db.Exec(sqlQuery, forum.Slug, forum.Title, forum.User)
	if err != nil {
		return err
	}
	return nil
}
