package main

import (
	"database/sql"
	"fmt"
	"github.com/BarniBl/DB-HW/cmd/api/handlers"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/labstack/echo"
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

	user := handlers.User{UserService: forum.NewUserService(db)}

	e := echo.New()
	e.POST("/user/:nickname/create", user.CreateUser)

	e.Logger.Warnf("start listening on %s", host)
	if err := e.Start(host); err != nil {
		e.Logger.Errorf("server error: %s", err)
	}

	e.Logger.Warnf("shutdown")

}
