package forum

import (
	"database/sql"

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
	SELECT distinct nick_name, email, full_name, about
	FROM "user" as u
			 LEFT JOIN post as p ON lower(p.author) = lower(nick_name)
			 LEFT JOIN thread as t ON lower(t.author) = lower(nick_name)
	WHERE lower(nick_name) < lower($3) AND (lower(p.forum) = lower($1) OR lower(t.forum) = lower($1))
	ORDER BY nick_name DESC 
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
		err := rows.Scan(&user.NickName, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (us *UserService) SelectUsersByForum(forum string, limit int, since string) (users []User, err error) {
	sqlQuery := `
	SELECT distinct nick_name, email, full_name, about
	FROM "user" as u
			 LEFT JOIN post as p ON lower(p.author) = lower(nick_name)
			 LEFT JOIN thread as t ON lower(t.author) = lower(nick_name)
	WHERE lower(nick_name) > lower($3) AND (lower(p.forum) = lower($1) OR lower(t.forum) = lower($1))
	ORDER BY nick_name ASC
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
		err := rows.Scan(&user.NickName, &user.Email, &user.FullName, &user.About)
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
