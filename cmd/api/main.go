package main

import (
	"database/sql"
	"fmt"
	"github.com/BarniBl/DB-HW/cmd/api/handlers"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

var (
	connectionString = "postgres://postgres:7396@localhost:5432/forum?sslmode=disable"
	host             = "0.0.0.0:8080"
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

	userService := forum.NewUserService(db)
	threadService := forum.NewThreadService(db)
	forumService := forum.NewForumService(db)
	postService := forum.NewPostService(db)

	user := handlers.User{UserService: userService}
	forum := handlers.Forum{ForumService: forumService, UserService: userService, ThreadService: threadService}

	e := echo.New()
	e.POST("/user/:nickname/create", user.CreateUser)
	e.GET("/user/:nickname/profile", user.GetProfile)
	e.POST("/user/:nickname/profile", user.EditProfile)

	e.POST("/forum/create", forum.CreateForum)
	e.POST("/forum/:slug/create", forum.CreateThread)
	e.GET("/forum/:slug/details", forum.GetForumDetails)
	e.GET("/forum/:slug/threads", forum.GetForumThreads)
	e.GET("/forum/:slug/users", forum.GetForumUsers)

	e.Use(middleware.Logger())
	e.Logger.Warnf("start listening on %s", host)
	if err := e.Start(host); err != nil {
		e.Logger.Errorf("server error: %s", err)
	}

	e.Logger.Warnf("shutdown")

}
