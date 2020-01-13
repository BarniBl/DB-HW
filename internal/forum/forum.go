package forum

import (
	"database/sql"
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
	where lower(f.slug) = lower($1)`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Slug, &forum.Title, &forum.User)
	if err != nil {
		return
	}
	sqlQuery = `SELECT count(*)
	FROM public.thread as t
	where lower(t.forum) = lower($1)`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Threads)
	if err != nil {
		return
	}
	sqlQuery = `
	SELECT count(*)
	FROM public.post as p
	where lower(p.forum) = lower($1)`
	err = fs.db.QueryRow(sqlQuery, slug).Scan(&forum.Posts)
	return
}

func (fs *ForumService) SelectForumBySlug(slug string) (forum Forum, err error) {
	sqlQuery := `
	SELECT f.slug, f.title, f.user
	FROM public.forum as f
	where lower(f.slug) = lower($1)`
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
	sqlQuery := `TRUNCATE public.vote, public.post, public.thread, public.forum, public.user RESTART IDENTITY CASCADE;`
	_, err = fs.db.Exec(sqlQuery)
	return
}

func (fs *ForumService) SelectStatus() (status Status, err error) {
	sqlQuery := `
	SELECT
	(SELECT COALESCE(SUM(public.posts), 0) FROM public.forum WHERE posts > 0) AS post, 
	(SELECT COALESCE(SUM(public.threads), 0) FROM public.forum WHERE threads > 0) AS thread, 
	(SELECT COUNT(*) FROM public.user) AS user,
	(SELECT COUNT(*) FROM public.forum) AS forum;`
	err = fs.db.QueryRow(sqlQuery).Scan(&status.Post, &status.Thread, &status.User, &status.Forum)
	return
}
