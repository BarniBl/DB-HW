package forum

import (
	"database/sql"
	"time"
)

type ThreadService struct {
	db *sql.DB
}

func NewThreadService(db *sql.DB) *ThreadService {
	return &ThreadService{db: db}
}

func (ts *ThreadService) SelectThreadByTitleAndForum(title string, forum string) (thread Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where Lower(t.title) = Lower($1) AND Lower(t.forum) = Lower($2)`
	var slug sql.NullString
	err = ts.db.QueryRow(sqlQuery, title, forum).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &slug, &thread.Title)
	if err != nil {
		return
	}
	if slug.Valid {
		thread.Slug = slug.String
	}
	sqlQuery = `SELECT sum(v.voice)
	FROM public.vote as v
	where Lower(v.thread) = Lower($1)`
	var votes sql.NullInt64
	err = ts.db.QueryRow(sqlQuery, title).Scan(&votes)
	if err != nil {
		return
	}
	if votes.Valid {
		thread.Votes = int(votes.Int64)
	}
	return
}

func (ts *ThreadService) SelectThreadById(id int) (thread Thread, err error) {
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

func (ts *ThreadService) InsertThread(thread Thread) error {
	sqlQuery := `INSERT INTO public.thread (author, created, message, title, forum, slug)
	VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := ts.db.Exec(sqlQuery, thread.Author, thread.Created, thread.Message, thread.Title, thread.Forum, thread.Slug)
	if err != nil {
		return err
	}
	return nil
}

func (ts *ThreadService) SelectThreadByForumDesc(forum string, limit int, since time.Time) (threads []Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where lower(t.forum) = lower($1) AND t.created <= $3
	ORDER BY t.created DESC
	LIMIT $2`
	rows, err := ts.db.Query(sqlQuery, forum, limit, since)
	if err != nil {
		return
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		threadScan := Thread{}
		slug := sql.NullString{}
		err := rows.Scan(&threadScan.Author, &threadScan.Created, &threadScan.Id, &threadScan.Forum, &threadScan.Message, &slug, &threadScan.Title)
		if err != nil {
			return threads, err
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

func (ts *ThreadService) SelectThreadByForum(forum string, limit int, since string, desc bool) (threads []Thread, err error) {
	var rows *sql.Rows
	if since == "" && !desc {
		sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
		FROM public.thread as t 
		WHERE lower(t.forum) = lower($1)
		ORDER BY t.created 
		LIMIT $2`
		rows, err = ts.db.Query(sqlQuery, forum, limit)
	} else if since != "" && !desc {
		sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
		FROM public.thread as t 
		WHERE lower(t.forum) = lower($1) AND t.created >= $3
		ORDER BY t.created 
		LIMIT $2`
		rows, err = ts.db.Query(sqlQuery, forum, limit, since)
	} else if since == "" && desc {
		sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
		FROM public.thread as t 
		WHERE lower(t.forum) = lower($1)
		ORDER BY t.created DESC 
		LIMIT $2`
		rows, err = ts.db.Query(sqlQuery, forum, limit)
	} else {
		sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
		FROM public.thread as t 
		WHERE lower(t.forum) = lower($1) AND t.created <= $3
		ORDER BY t.created DESC 
		LIMIT $2`
		rows, err = ts.db.Query(sqlQuery, forum, limit, since)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		threadScan := Thread{}
		slug := sql.NullString{}
		err := rows.Scan(&threadScan.Author, &threadScan.Created, &threadScan.Id, &threadScan.Forum, &threadScan.Message, &slug, &threadScan.Title)
		if err != nil {
			return threads, err
		}
		if slug.Valid {
			threadScan.Slug = slug.String
		}
		sqlQuery := `SELECT sum(v.voice)
		FROM public.vote as v
		where v.thread = $1`
		err = ts.db.QueryRow(sqlQuery, threadScan.Title).Scan(&threadScan.Votes)
		threads = append(threads, threadScan)
	}
	return
}

func (ts *ThreadService) FindThreadBySlug(slug string) (thread Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where lower(t.slug) = lower($1)`
	err = ts.db.QueryRow(sqlQuery, slug).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title)
	return
}

func (ts *ThreadService) FindThreadById(id int) (thread Thread, err error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where t.id = $1`
	err = ts.db.QueryRow(sqlQuery, id).Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title)
	return
}
