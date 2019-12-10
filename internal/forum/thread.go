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

func (ts *ThreadService) SelectThreadByTitle(title string) (threadSl []input.Thread, er error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title, t.votes
	FROM public.thread as t 
	where t.title = $1`
	threadSlice := make([]input.Thread, 0)
	rows, err := ts.db.Query(sqlQuery, title)
	if err != nil {
		return threadSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		thread := input.Thread{}
		slug := sql.NullString{}
		votes := sql.NullInt64{}
		err := rows.Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &slug, &thread.Title, &votes)
		if err != nil {
			return threadSlice, err
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		if votes.Valid {
			thread.Votes = int(votes.Int64)
		}
		threadSlice = append(threadSlice, thread)
	}
	return threadSlice, nil
}

func (ts *ThreadService) SelectThreadById(id int) (threadSl []input.Thread, er error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title, t.votes
	FROM public.thread as t 
	where t.id = $1`
	threadSlice := make([]input.Thread, 0)
	rows, err := ts.db.Query(sqlQuery, id)
	if err != nil {
		return threadSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		thread := input.Thread{}
		slug := sql.NullString{}
		votes := sql.NullInt64{}
		err := rows.Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &slug, &thread.Title, &votes)
		if err != nil {
			return threadSlice, err
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		if votes.Valid {
			thread.Votes = int(votes.Int64)
		}
		threadSlice = append(threadSlice, thread)
	}
	return threadSlice, nil
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

func (ts *ThreadService) SelectThreadByForumDesc(forum string, limit, since int) (threadSl []input.Thread, er error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title, t.votes
	FROM public.thread as t 
	where t.forum = $1
	ORDER BY t.created DESC
	LIMIT $2 OFFSET $3`
	threadSlice := make([]input.Thread, 0)
	rows, err := ts.db.Query(sqlQuery, forum)
	if err != nil {
		return threadSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		thread := input.Thread{}
		slug := sql.NullString{}
		votes := sql.NullInt64{}
		err := rows.Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &slug, &thread.Title, &votes)
		if err != nil {
			return threadSlice, err
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		if votes.Valid {
			thread.Votes = int(votes.Int64)
		}
		threadSlice = append(threadSlice, thread)
	}
	return threadSlice, nil
}

func (ts *ThreadService) SelectThreadByForum(forum string, limit, since int) (threadSl []input.Thread, er error) {
	sqlQuery := `SELECT t.author, t.created, t.id, t.forum, t.message, t.slug, t.title, t.votes
	FROM public.thread as t 
	where t.forum = $1
	ORDER BY t.created
	LIMIT $2 OFFSET $3`
	threadSlice := make([]input.Thread, 0)
	rows, err := ts.db.Query(sqlQuery, forum)
	if err != nil {
		return threadSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		thread := input.Thread{}
		slug := sql.NullString{}
		votes := sql.NullInt64{}
		err := rows.Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &slug, &thread.Title, &votes)
		if err != nil {
			return threadSlice, err
		}
		if slug.Valid {
			thread.Slug = slug.String
		}
		if votes.Valid {
			thread.Votes = int(votes.Int64)
		}
		threadSlice = append(threadSlice, thread)
	}
	return threadSlice, nil
}
