package forum

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/input"
)

type ThreadService struct {
	db *sql.DB
}

func NewThreadService(db *sql.DB) *ThreadService {
	return &ThreadService{db: db}
}

func (ts *ThreadService) SelectThreadByTitle(title string) (thread input.Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title, t.votes
	FROM public.thread as t 
	where t.title = $1`
	err = ts.db.QueryRow(sqlQuery, title).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, )
	if err != nil {
		return
	}
	sqlQuery = `SELECT sum(v.voice)
	FROM public.vote as v
	where v.thread = $1`
	err = ts.db.QueryRow(sqlQuery, title).Scan(&thread.Votes)
	return
}

func (ts *ThreadService) SelectThreadById(id int) (thread input.Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where t.id = $1`
	err = ts.db.QueryRow(sqlQuery, id).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title)
	if err != nil {
		return
	}
	sqlQuery = `SELECT sum(v.voice)
	FROM public.vote as v
	where v.thread = $1`
	err = ts.db.QueryRow(sqlQuery, thread.Title).Scan(&thread.Votes)
	return
}

func (ts *ThreadService) InsertThread(thread input.Thread) error {
	sqlQuery := `INSERT INTO public.thread (author, created, message, title, forum)
	VALUES ($1,$2,$3,$4,$5)`
	_, err := ts.db.Exec(sqlQuery, thread.Author, thread.Created, thread.Message, thread.Title, thread.Forum)
	if err != nil {
		return err
	}
	return nil
}

func (ts *ThreadService) SelectThreadByForumDesc(forum string, limit, since int) (threads []input.Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where t.forum = $1
	ORDER BY t.created DESC
	LIMIT $2 OFFSET $3`
	rows, err := ts.db.Query(sqlQuery, forum, limit, since)
	if err != nil {
		return threadSlice, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		threadScan := input.Thread{}
		slug := sql.NullString{}
		err := rows.Scan(&threadScan.Author, &threadScan.Created, &threadScan.Id, &threadScan.Forum, &threadScan.Message, &slug, &threadScan.Title)
		if err != nil {
			return
		}
		if slug.Valid {
			threadScan.Slug = slug.String
		}
		sqlQuery = `SELECT sum(v.voice)
		FROM public.vote as v
		where v.thread = $1`
		err = ts.db.QueryRow(sqlQuery, threadScan.Title).Scan(&threadScan.Votes)
		threads = append(threads, threadScan)
	}
	return
}

func (ts *ThreadService) SelectThreadByForum(forum string, limit, since int) (threads []input.Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where t.forum = $1
	ORDER BY t.created
	LIMIT $2 OFFSET $3`
	rows, err := ts.db.Query(sqlQuery, forum, limit, since)
	if err != nil {
		return threadSlice, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		threadScan := input.Thread{}
		slug := sql.NullString{}
		err := rows.Scan(&threadScan.Author, &threadScan.Created, &threadScan.Id, &threadScan.Forum, &threadScan.Message, &slug, &threadScan.Title)
		if err != nil {
			return
		}
		if slug.Valid {
			threadScan.Slug = slug.String
		}
		sqlQuery = `SELECT sum(v.voice)
		FROM public.vote as v
		where v.thread = $1`
		err = ts.db.QueryRow(sqlQuery, threadScan.Title).Scan(&threadScan.Votes)
		threads = append(threads, threadScan)
	}
	return
}
