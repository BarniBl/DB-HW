package handlers

import (
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/BarniBl/DB-HW/internal/input"
	"github.com/BarniBl/DB-HW/internal/output"
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
	UserService *forum.UserService
}

func (h *User) CreateUser(ctx echo.Context) (Err error) {
	nickName := ctx.Param("nickname")
	if nickName == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}
	newUser := input.User{}
	if err := ctx.Bind(&newUser); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}
	newUser.NickName = nickName

	userSlice, err := h.UserService.SelectUserByNickNameOrEmail(newUser.NickName, newUser.Email)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}

	if len(userSlice) > 0 {
		return ctx.JSON(http.StatusConflict, userSlice)
	}

	if err = h.UserService.InsertUser(newUser); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}

	return ctx.JSON(http.StatusCreated, newUser)
}

func (h *User) GetProfile(ctx echo.Context) (Err error) {
	nickName := ctx.Param("nickname")
	if nickName == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}
	userSlice, err := h.UserService.SelectUserByNickName(nickName)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}

	if len(userSlice) != 1 {
		return ctx.JSON(http.StatusNotFound, userSlice)
	}

	return ctx.JSON(http.StatusOK, userSlice[0])
}

func (h *User) EditProfile(ctx echo.Context) (Err error) {
	nickName := ctx.Param("nickname")
	if nickName == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}
	editUser := input.User{}
	if err := ctx.Bind(&editUser); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}
	editUser.NickName = nickName

	userSlice, err := h.UserService.SelectUserByNickNameOrEmail(editUser.NickName, editUser.Email)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}
	if len(userSlice) > 1 {
		return ctx.JSON(http.StatusConflict, userSlice)
	}
	if userSlice[0].NickName != editUser.NickName {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message:"Error"})
	}

	if err = h.UserService.UpdateUser(editUser); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message:"Error"})
	}

	return ctx.JSON(http.StatusOK, editUser)
}
