package handlers

import (
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/BarniBl/DB-HW/internal/input"
	"github.com/BarniBl/DB-HW/internal/output"
	"github.com/labstack/echo"
	"net/http"
)

type Forum struct {
	ForumService *forum.ForumService
	UserService  *forum.UserService
}

func (h *Forum) CreateForum(ctx echo.Context) (Err error) {

	newForum := input.Forum{}
	if err := ctx.Bind(&newForum); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	forumSlice, err := h.ForumService.SelectForumBySlug(newForum.Slug)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	if len(forumSlice) > 0 {
		return ctx.JSON(http.StatusConflict, forumSlice)
	}

	if err := h.UserService.CheckUser(newForum.User); err != nil {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	if err = h.ForumService.InsertForum(newForum); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	return ctx.JSON(http.StatusCreated, newForum)
}
