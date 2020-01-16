package main

import (
	"database/sql"
	"fmt"
	"github.com/BarniBl/DB-HW/cmd/api/handlers"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"io/ioutil"
)

var (
	connectionString = "postgres://forum:7396@localhost:5432/forum?sslmode=disable"
	host             = "0.0.0.0:5000"
)

func main() {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = LoadSchemaSQL(db)
	if err != nil {
		fmt.Println(err)
	}
	userService := forum.NewUserService(db)
	threadService := forum.NewThreadService(db)
	forumService := forum.NewForumService(db)
	postService := forum.NewPostService(db)

	user := handlers.User{UserService: userService}
	forum := handlers.Forum{ForumService: forumService, UserService: userService, ThreadService: threadService}
	post := handlers.Post{PostService: postService, ForumService: forumService, UserService: userService, ThreadService: threadService}

	e := echo.New()
	e.POST("/user/:nickname/create", user.CreateUser)
	e.GET("/user/:nickname/profile", user.GetProfile)
	e.POST("/user/:nickname/profile", user.EditProfile)

	e.POST("/forum/create", forum.CreateForum)
	e.POST("/forum/:slug/create", forum.CreateThread)
	e.GET("/forum/:slug/details", forum.GetForumDetails)
	e.GET("/forum/:slug/threads", forum.GetForumThreads)
	e.GET("/forum/:slug/users", forum.GetForumUsers)

	e.GET("/post/:id/details", post.GetFullPost)
	e.POST("/post/:id/details", post.EditMessage)

	e.GET("/thread/:slug_or_id/details", post.GetThread)
	e.POST("/thread/:slug_or_id/details", post.EditThread)
	e.GET("/thread/:slug_or_id/posts", post.GetPosts)
	e.POST("/thread/:slug_or_id/create", post.CreatePosts)
	e.POST("/thread/:slug_or_id/vote", post.CreateVote)

	e.POST("/service/clear", forum.Clean)
	e.GET("/service/status", forum.Status)

	e.Use(middleware.Logger())
	e.Logger.Warnf("start listening on %s", host)
	if err := e.Start(host); err != nil {
		e.Logger.Errorf("server error: %s", err)
	}

	e.Logger.Warnf("shutdown")

}

const dbSchema = "dum_hw_pdb.sql"

func LoadSchemaSQL(db *sql.DB) error {

	content, err := ioutil.ReadFile(dbSchema)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(content)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}
