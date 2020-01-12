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

func (fs *ForumService) SelectFullForumBySlug(slug string) (forum Forum, err error) {
	sqlQuery := `SELECT f.slug, f.title, f.user
	FROM public.forum as f
	where f.slug = $1`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Slug, &forum.Title, &forum.User)
	if err != nil {
		return
	}
	sqlQuery = `SELECT count(*)
	FROM public.thread as t
	where t.forum = $1`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Threads)
	if err != nil {
		return
	}
	sqlQuery = `
	SELECT count(*)
	FROM public.post as p
	where p.forum = $1`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Posts)
	return
}

func (fs *ForumService) SelectForumBySlug(slug string) (forum Forum, err error) {
	sqlQuery := `
	SELECT f.slug, f.title, f.user
	FROM public.forum as f
	where f.slug = $1`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Slug, &forum.Title, &forum.User)
	return
}

func (fs *ForumService) InsertForum(forum Forum) (err error) {
	sqlQuery := `INSERT INTO public.forum (slug, title, "user")
	VALUES ($1,$2,$3)`
	_, err = fs.db.Exec(sqlQuery, forum.Slug, forum.Title, forum.User)
	return
}

func (fs *ForumService) Clean() (err error) {
	sqlQuery := `TRUNCATE forum.vote, forum.post, forum.thread, forum.forum, forum.user RESTART IDENTITY CASCADE;`
	_, err = fs.db.Exec(sqlQuery)
	return
}

func (fs *ForumService) SelectStatus() (status Status, err error) {
	sqlQuery := `
	SELECT
	(SELECT COALESCE(SUM(forum.posts), 0) FROM forum.forum WHERE posts > 0) AS post, 
	(SELECT COALESCE(SUM(forum.threads), 0) FROM forum.forum WHERE threads > 0) AS thread, 
	(SELECT COUNT(*) FROM forum.user) AS user,
	(SELECT COUNT(*) FROM forum.forum) AS forum;`
	err = fs.db.QueryRow(sqlQuery).Scan(&status.Post, &status.Thread, &status.User, &status.Forum)
	return
}
