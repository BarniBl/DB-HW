package forum

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/input"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (us *UserService) SelectUserByNickNameOrEmail(nickName, email string) (userSl []input.User, er error) {
	sqlQuery := `SELECT u.nick_name, u.email, u.full_name, u.about
	FROM public.user as u 
	where u.nick_name = $1 or u.email = $2`
	userSlice := make([]input.User, 0)
	rows, err := us.db.Query(sqlQuery, nickName, email)
	if err != nil {
		return userSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		user := input.User{}
		err := rows.Scan(&user.NickName, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return userSlice, err
		}
		userSlice = append(userSlice, user)
	}
	return userSlice, nil
}

func (us *UserService) SelectUserByNickName(nickName string) (userSl []input.User, er error) {
	sqlQuery := `SELECT u.nick_name, u.email, u.full_name, u.about
	FROM public.user as u 
	where u.nick_name = $1`
	userSlice := make([]input.User, 0)
	rows, err := us.db.Query(sqlQuery, nickName)
	if err != nil {
		return userSlice, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			er = err
		}
	}()

	for rows.Next() {
		user := input.User{}
		err := rows.Scan(&user.NickName, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return userSlice, err
		}
		userSlice = append(userSlice, user)
	}
	return userSlice, nil
}

func (us *UserService) InsertUser(user input.User) error {
	sqlQuery := `INSERT INTO public.user (nick_name,email,full_name,about)
	VALUES ($1,$2,$3,$4)
	returning nick_name`
	var nickName string
	err := us.db.QueryRow(sqlQuery, user.NickName, user.Email, user.FullName, user.About).Scan(&nickName)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateUser(user input.User) error {
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
