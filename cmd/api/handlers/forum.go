package handlers

import (
	"fmt"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/BarniBl/DB-HW/internal/input"
	"github.com/BarniBl/DB-HW/internal/output"
	"github.com/labstack/echo"
	"net/http"
)

type Forum struct {
	ForumService  *forum.ForumService
	UserService   *forum.UserService
	ThreadService *forum.ThreadService
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

func (h *Forum) CreateThread(ctx echo.Context) (Err error) {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	newThread := input.Thread{}
	if err := ctx.Bind(&newThread); err != nil {
		fmt.Println(err)
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	newThread.Forum = slug

	forumSlice, err := h.ForumService.SelectForumBySlug(newThread.Forum)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	if len(forumSlice) != 1 {
		return ctx.JSON(http.StatusNotFound, forumSlice)
	}

	if err := h.UserService.CheckUser(newThread.Author); err != nil {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	threadCopySlice, err := h.ThreadService.SelectThreadByTitle(newThread.Title)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	if len(threadCopySlice) > 0 {
		return ctx.JSON(http.StatusConflict, threadCopySlice)
	}

	err = h.ThreadService.InsertThread(newThread)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	threadSlice, err := h.ThreadService.SelectThreadByTitle(newThread.Title)
	if len(threadSlice) != 1 {
		return ctx.JSON(http.StatusConflict, output.ErrorMessage{Message: "Error"})
	}

	return ctx.JSON(http.StatusCreated, threadSlice[0])
}

func (h *Forum) GetForumDetails(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	forumSlice, err := h.ForumService.SelectForumBySlug(slug)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}
	if len(forumSlice) != 1 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	return ctx.JSON(http.StatusOK, forumSlice[0])
}
