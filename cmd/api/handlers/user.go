package handlers

import (
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/BarniBl/DB-HW/internal/input"
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
	UserService *forum.UserService
}

func (h *User) CreateUser(ctx echo.Context) (Err error) {
	nickName := ctx.Param("nickname")
	if nickName == "" {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	newUser := input.User{}
	if err := ctx.Bind(&newUser); err != nil {
		return
	}
	newUser.NickName = nickName

	userSlice, err := h.UserService.SelectUserByNickNameOrEmail(newUser.NickName, newUser.Email)
	if err != nil {
		return
	}

	if len(userSlice) > 0 {
		return ctx.JSON(http.StatusConflict, userSlice)
	}

	if err = h.UserService.InsertUser(newUser); err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, newUser)
}
