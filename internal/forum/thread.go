package forum

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"strconv"
	"time"
)

type ThreadService struct {
	db *pgx.ConnPool
}

func NewThreadService(db *pgx.ConnPool) *ThreadService {
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
	where Lower(v.thread_id) = Lower($1)`
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

func (ts *ThreadService) SelectThreadBySlug(threadSlug string) (thread Thread, err error) {
	sqlQuery := `SELECT t.id, t.author, t.created, t.id, t.forum, t.message, t.slug, t.title
	FROM public.thread as t 
	where Lower(t.slug) = Lower($1) `
	var slug sql.NullString
	err = ts.db.QueryRow(sqlQuery, threadSlug).Scan(&thread.Id, &thread.Author, &thread.Created, &thread.Id, &thread.Forum, &thread.Message, &slug, &thread.Title)
	if err != nil {
		return
	}
	if slug.Valid {
		thread.Slug = slug.String
	}
	sqlQuery = `SELECT sum(v.voice)
	FROM public.vote as v
	where v.thread_id = $1`
	var votes sql.NullInt64
	err = ts.db.QueryRow(sqlQuery, thread.Id).Scan(&votes)
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
	where v.thread_id = $1`
	var votes sql.NullInt64
	err = ts.db.QueryRow(sqlQuery, thread.Id).Scan(&votes)
	if err != nil {
		return
	}
	if votes.Valid {
		thread.Votes = int(votes.Int64)
	}
	return
}

func (ts *ThreadService) InsertThread(thread Thread) (id int, err error) {
	sqlQuery := `INSERT INTO public.thread (author, created, message, title, forum, slug)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`
	err = ts.db.QueryRow(sqlQuery, thread.Author, thread.Created, thread.Message, thread.Title, thread.Forum, thread.Slug).Scan(&id)
	return
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

	defer rows.Close()

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
		where v.thread_id = $1`
		err = ts.db.QueryRow(sqlQuery, threadScan.Title).Scan(&threadScan.Votes)
		threads = append(threads, threadScan)
	}
	return
}

func (ts *ThreadService) SelectThreadByForum(forum string, limit int, since string, desc bool) (threads []Thread, err error) {
	var rows *pgx.Rows
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

	defer rows.Close()

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
		where v.thread_id = $1`
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

func (ts *ThreadService) InsertVote(vote Vote) (err error) {
	sqlQuery := `INSERT INTO public.vote (nick_name, voice, thread_id)
	VALUES ($1,$2,$3)`
	_, err = ts.db.Exec(sqlQuery, vote.NickName, vote.Voice, vote.ThreadId)
	return
}

func (ts *ThreadService) UpdateVote(vote Vote) (countUpdatedRows int64, err error) {
	sqlQuery := `
	UPDATE public.vote SET voice = $1
	where lower(vote.nick_name) = lower($2) AND vote.thread_id = $3`
	result, err := ts.db.Exec(sqlQuery, vote.Voice, vote.NickName, vote.ThreadId)
	if err != nil {
		return
	}
	countUpdatedRows = result.RowsAffected()
	return
}

func (ts *ThreadService) UpdateThread(thread Thread) (err error) {
	sqlQuery := `
	UPDATE public.thread SET message = $1,
	                         title = $2
	where thread.id = $3`
	_, err = ts.db.Exec(sqlQuery, thread.Message, thread.Title, thread.Id)
	return
}

func (ts *ThreadService) FindVote(vote Vote) (voted bool, err error) {
	sqlQuery := `SELECT v.nick_name
	FROM public.vote as v 
	where lower(v.nick_name) = lower($1) AND v.thread_id = $2`
	var vot sql.NullString
	err = ts.db.QueryRow(sqlQuery, vote.NickName, vote.ThreadId).Scan(&vot)
	return vot.Valid, err
}

func (ts *ThreadService) SelectPosts(threadID int, limit, since, sort, desc string) (Posts []Post, Err error) {
	posts := []Post{}

	var rows *pgx.Rows
	var err error
	if sort == "flat" {
		if desc == "false" {
			sqlQuery := "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
				"FROM public.post as p WHERE p.thread = $1 AND p.id > $3 ORDER BY p.id LIMIT $2"
			rows, err = ts.db.Query(sqlQuery, threadID, limit, since)
		} else {
			sqlQuery := "SELECT p.author, p.created, p.forum, p.id, p.is_edited, p.message, p.parent, p.thread " +
				"FROM public.post as p WHERE p.thread = $1 AND p.id < $3 ORDER BY p.id DESC LIMIT $2"
			rows, err = ts.db.Query(sqlQuery, threadID, limit, since)
		}

	} else if sort == "tree" {
		if desc == "false" {
			if since != "0" && since != "999999999" {
				sqlQuery := "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
					"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
					"FROM public.post as T1 WHERE T1.parent = 0 and T1.thread = $1 " +
					"union " +
					"select T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
					"FROM public.post T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
					") " +
					"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY root, PATH LIMIT $2;"
				rows, err = ts.db.Query(sqlQuery, threadID, 100000)
			} else {
				sqlQuery := "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
					"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
					"FROM public.post as T1 WHERE T1.parent = 0 and T1.thread = $1 " +
					"union " +
					"select T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
					"FROM public.post T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
					") " +
					"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY root, PATH LIMIT $2;"
				rows, err = ts.db.Query(sqlQuery, threadID, limit)
			}
		} else {
			if since != "0" && since != "999999999" {
				sqlQuery := "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
					"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
					"FROM public.post as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
					"union " +
					"select  T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST (temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1, root " +
					"FROM public.post as T2 INNER JOIN temp1 ON (temp1.id = T2.parent) " +
					") " +
					"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY PATH;"
				rows, err = ts.db.Query(sqlQuery, threadID)
			} else {
				sqlQuery := "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
					"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
					"FROM public.post as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
					"union " +
					"select  T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST (temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1, root " +
					"FROM public.post as T2 INNER JOIN temp1 ON (temp1.id = T2.parent) " +
					") " +
					"select author, created, forum, id, is_edited, message, parent, thread from temp1 WHERE id < $3 ORDER BY PATH DESC LIMIT $2;"
				rows, err = ts.db.Query(sqlQuery, threadID, limit, 1000000)
			}
		}
	} else if sort == "parent_tree" {
		if desc == "false" {
			sqlQuery := "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
				"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
				"FROM public.post as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
				"union " +
				"select T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
				"FROM public.post T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
				") " +
				"select author, created, forum, id, is_edited, message, parent, thread from temp1 ORDER BY root, PATH;"
			rows, err = ts.db.Query(sqlQuery, threadID)
		} else {
			sqlQuery := "WITH RECURSIVE temp1 (author, created, forum, id, is_edited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
				"SELECT T1.author, T1.created, T1.forum, T1.id, T1.is_edited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
				"FROM public.post as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
				"union " +
				"select  T2.author, T2.created, T2.forum, T2.id, T2.is_edited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
				"FROM public.post as T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
				") " +
				"select author, created, forum, id, is_edited, message, parent, thread  from temp1 ORDER BY root desc, PATH;"
			rows, err = ts.db.Query(sqlQuery, threadID)
		}
	}

	if sort != "parent_tree" {
		defer rows.Close()
		if err != nil {
			fmt.Println(err)
			return posts, err
		}

		for rows.Next() {
			scanPost := Post{}
			//var timetz time.Time
			err := rows.Scan(&scanPost.Author, &scanPost.Created, &scanPost.Forum,
				&scanPost.Id, &scanPost.IsEdited, &scanPost.Message, &scanPost.Parent,
				&scanPost.Thread)
			if err != nil {
				return posts, err
			}
			//scanPost.Created = timetz.Format(time.RFC3339Nano)
			posts = append(posts, scanPost)
		}
	} else {
		if err != nil {
			rows.Close()
			return posts, err
		}

		count := 0
		limitDigit, _ := strconv.Atoi(limit)

		for rows.Next() {
			scanPost := Post{}
			//var timetz time.Time
			err := rows.Scan(&scanPost.Author, &scanPost.Created, &scanPost.Forum,
				&scanPost.Id, &scanPost.IsEdited, &scanPost.Message, &scanPost.Parent,
				&scanPost.Thread)
			if err != nil {
				return posts, err
			}

			if scanPost.Parent == 0 {
				count = count + 1
			}
			if count > limitDigit && (since == "0" || since == "999999999") {
				break
			} else {
				//scanPost.Created = timetz.Format(time.RFC3339Nano)
				posts = append(posts, scanPost)
			}

		}
		rows.Close()
	}

	if since != "0" && since != "999999999" && sort == "tree" {
		limitDigit, _ := strconv.Atoi(limit)
		sinceDigit, _ := strconv.Atoi(since)
		var sincePosts []Post
		counter := 0
		//for ; posts[counter].ID <= sinceDigit && counter < len(posts); {
		//	counter++
		//}
		if desc == "false" {
			startIndex := 1000000000
			//postMinStartIndex
			minValue := 100000000000
			for i := 0; i < len(posts); i++ {
				if (posts[i].Id == sinceDigit) {
					startIndex = i + 1
					break
				}
				if (posts[i].Id > sinceDigit) && (posts[i].Id < minValue) {
					startIndex = i
					minValue = posts[i].Id
				}
			}
			sincePostsCount := 0
			counter = startIndex
			for ; sincePostsCount < limitDigit && counter < len(posts); {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter++
			}
		} else {
			startIndex := -1000000000
			//postMinStartIndex
			maxValue := 0
			for i := len(posts) - 1; i >= 0; i-- {
				if (posts[i].Id == sinceDigit) {
					startIndex = i - 1
					break
				}
				if (posts[i].Id < sinceDigit) && (posts[i].Id > maxValue) {
					startIndex = i
					maxValue = posts[i].Id
				}
			}

			//xsort.Slice(posts[0:startIndex], func(i, j int) bool { return posts[i].ID < posts[j].ID})
			sincePostsCount := 0
			counter = startIndex
			for ; sincePostsCount < limitDigit && counter >= 0; {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter--
			}
		}
		return sincePosts, nil
	}

	if since != "0" && since != "999999999" && sort == "parent_tree" {
		limitDigit, _ := strconv.Atoi(limit)
		sinceDigit, _ := strconv.Atoi(since)
		var sincePosts []Post
		counter := 0
		if desc == "false" {
			startIndex := 1000000000
			minValue := 100000000000
			for i := 0; i < len(posts); i++ {
				if (posts[i].Id == sinceDigit) {
					startIndex = i + 1
					break
				}
				if (posts[i].Id > sinceDigit) && (posts[i].Id < minValue) {
					startIndex = i
					minValue = posts[i].Id
				}
			}
			sincePostsCount := 0
			counter = startIndex
			for ; sincePostsCount < limitDigit && counter < len(posts); {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				sincePostsCount++
				counter++
			}
		} else {
			startIndex := -1000000000
			//postMinStartIndex
			maxValue := 100000000000
			for i := len(posts) - 1; i >= 0; i-- {
				if (posts[i].Id == sinceDigit) {
					startIndex = i + 1
					break
				}
				if (posts[i].Id < sinceDigit) && (posts[i].Id < maxValue) {
					startIndex = i
					maxValue = posts[i].Id
				}
			}

			//xsort.Slice(posts[0:startIndex], func(i, j int) bool { return posts[i].ID < posts[j].ID})
			sincePostsCount := 0
			counter = startIndex
			for ; sincePostsCount < limitDigit && counter < len(posts); {
				scanPost := Post{}
				scanPost = posts[counter]
				sincePosts = append(sincePosts, scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter++
			}
		}
		return sincePosts, nil
	}

	return posts, nil
}
