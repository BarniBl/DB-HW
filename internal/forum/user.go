package forum

import (
	"bytes"
	"database/sql"
	xsort "sort"
	"strconv"
	"strings"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) SelectUserByNickNameOrEmail(nickName, email string) (users []User, err error) {
	sqlQuery := `SELECT u.nick_name, u.email, u.full_name, u.about
	FROM public.user as u 
	where lower(u.nick_name) = lower($1) or lower(u.email) = lower($2)`
	rows, err := us.db.Query(sqlQuery, nickName, email)
	if err != nil {
		return users, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		userScan := User{}
		err := rows.Scan(&userScan.NickName, &userScan.Email, &userScan.FullName, &userScan.About)
		if err != nil {
			return users, err
		}
		users = append(users, userScan)
	}
	return users, nil
}

func (us *UserService) SelectUserByNickName(nickName string) (user User, err error) {
	sqlQuery := `SELECT u.nick_name, u.email, u.full_name, u.about
	FROM public.user as u 
	where lower(u.nick_name) = lower($1)`
	err = us.db.QueryRow(sqlQuery, nickName).Scan(&user.NickName, &user.Email, &user.FullName, &user.About)
	return
}

func (us *UserService) SelectUsersByForumDesc(forum string, limit int, since string) (users []User, err error) {
	sqlQuery := `
	SELECT distinct LOWER(nick_name) COLLATE "C", nick_name, email, full_name, about
	FROM public.user as u
			 LEFT JOIN public.post as p ON lower(p.author) = lower(nick_name)
			 LEFT JOIN public.thread as t ON lower(t.author) = lower(nick_name)
	WHERE lower(nick_name) < lower($3) COLLATE "C" AND (lower(p.forum) = lower($1) OR lower(t.forum) = lower($1))
	ORDER BY LOWER(nick_name) COLLATE "C" DESC 
	LIMIT $2`
	rows, err := us.db.Query(sqlQuery, forum, limit, since)
	if err != nil {
		return
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		user := User{}
		temp := ""
		err := rows.Scan(&temp, &user.NickName, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (us *UserService) SelectUsersByForum(forum string, limit int, since string) (users []User, err error) {
	sqlQuery := `
	SELECT distinct LOWER(nick_name) COLLATE "C", nick_name, email, full_name, about
	FROM public.user as u
			 LEFT JOIN public.post as p ON lower(p.author) = lower(nick_name)
			 LEFT JOIN public.thread as t ON lower(t.author) = lower(nick_name)
	WHERE lower(nick_name) > lower($3) COLLATE "C" AND (lower(p.forum) = lower($1) OR lower(t.forum) = lower($1))
	ORDER BY LOWER(nick_name) COLLATE "C" ASC
	LIMIT $2`
	rows, err := us.db.Query(sqlQuery, forum, limit, since)
	if err != nil {
		return
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		user := User{}
		temp := ""
		err := rows.Scan(&temp, &user.NickName, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (us *UserService) SelectUsersByForumAntiSince(forum string, limit int) (users []User, err error) {
	sqlQuery := `
	SELECT distinct lower(nick_name) COLLATE "C", nick_name, email, full_name, about
	FROM public.user as u
			 LEFT JOIN public.post as p ON lower(p.author) = lower(nick_name)
			 LEFT JOIN public.thread as t ON lower(t.author) = lower(nick_name)
	WHERE (lower(p.forum) = lower($1) OR lower(t.forum) = lower($1))
	ORDER BY lower(nick_name) COLLATE "C" ASC
	LIMIT $2`
	rows, err := us.db.Query(sqlQuery, forum, limit)
	if err != nil {
		return
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		user := User{}
		temp := ""
		err := rows.Scan(&temp, &user.NickName, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (us *UserService) InsertUser(user User) error {
	sqlQuery := `INSERT INTO public.user (nick_name,email,full_name,about)
	VALUES ($1,$2,$3,$4)`
	_, err := us.db.Exec(sqlQuery, user.NickName, user.Email, user.FullName, user.About)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateUser(user User) error {
	sqlQuery := `UPDATE public.user
	SET email = $1, 
		full_name = $2, 	
		about = $3
		WHERE nick_name = $4`
	_, err := us.db.Exec(sqlQuery, user.Email, user.FullName, user.About, user.NickName)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) FindUserByNickName(nickName string) (findNickName string, err error) {
	sqlQuery := `SELECT u.nick_name
	FROM public.user as u 
	where lower(u.nick_name) = lower($1)`
	err = us.db.QueryRow(sqlQuery, nickName).Scan(&findNickName)
	return
}

func (us *UserService) SelectAllUsersByForum(slug, limit, since, desc string) (users []User, err error) {
	var rows *sql.Rows
	//if since == "" {
	if desc == "false" {
		sqlQuery := "SELECT u.about, u.email, u.full_name, u.nick_name " +
			`FROM public."user" as u ` +
			"WHERE u.nick_name IN ( " +
			"SELECT t.author AS nick_name " +
			"FROM public.thread as t " +
			"WHERE lower(t.forum) = lower($1) " +
			"UNION " +
			"SELECT p.author AS nick_name " +
			"FROM public.post as p " +
			"WHERE lower(p.forum) = lower($1) ) " +
			"ORDER BY lower(u.nick_name) " +
			"LIMIT $2;"
		rows, err = us.db.Query(sqlQuery, slug, limit)
	} else {
		sqlQuery := "SELECT u.about, u.email, u.full_name, u.nick_name " +
			`FROM public."user" as u ` +
			"WHERE u.nick_name IN ( " +
			"SELECT t.author AS nick_name " +
			"FROM public.thread as t " +
			"WHERE lower(t.forum) = lower($1) " +
			"UNION " +
			"SELECT p.author AS nick_name " +
			"FROM public.post as p " +
			"WHERE lower(p.forum) = lower($1) ) " +
			"ORDER BY u.nick_name DESC " +
			"LIMIT $2;"
		rows, err = us.db.Query(sqlQuery, slug, limit)
	}
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		scanUser := User{}
		err = rows.Scan(&scanUser.About, &scanUser.Email, &scanUser.FullName,
			&scanUser.NickName)
		if err != nil {
			return
		}
		users = append(users, scanUser)
	}

	resUsers := []User{}

	limitDigit, _ := strconv.Atoi(limit)

	if desc == "false" {

		xsort.Slice(users, func(i, j int) bool { return bytes.Compare([]byte(strings.ToLower(users[i].NickName)), []byte(strings.ToLower(users[j].NickName))) < 0 })

		if since == "" {
			for i := 0; i < limitDigit && i < len(users); i++ {
				resUsers = append(resUsers, users[i])
			}
		} else {
			j := 0
			for i := 0; j < limitDigit && i < len(users); {
				if bytes.Compare([]byte(strings.ToLower(users[i].NickName)), []byte(strings.ToLower(since))) > 0 {
					resUsers = append(resUsers, users[i])
					j++
				}
				i++
			}
		}
	} else {

		xsort.Slice(users, func(i, j int) bool { return bytes.Compare([]byte(strings.ToLower(users[i].NickName)), []byte(strings.ToLower(users[j].NickName))) > 0 })

		if since == "" {
			for i := 0; i < limitDigit && i < len(users); i++ {
				resUsers = append(resUsers, users[i])
			}
		} else {
			j := 0
			for i := 0; j < limitDigit && i < len(users); {
				if bytes.Compare([]byte(strings.ToLower(users[i].NickName)), []byte(strings.ToLower(since))) < 0 {
					resUsers = append(resUsers, users[i])
					j++
				}
				i++
			}
		}
	}

	return
}
