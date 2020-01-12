package handlers

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/BarniBl/DB-HW/internal/input"
	"github.com/BarniBl/DB-HW/internal/output"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type Post struct {
	ForumService  *forum.ForumService
	UserService   *forum.UserService
	ThreadService *forum.ThreadService
	PostService   *forum.PostService
}

func (h *Post) GetFullPost(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if id < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	post, err := h.PostService.SelectPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find post"})
		}
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	fullPost := output.FullPost{Post: post}

	user, err := h.UserService.SelectUserByNickName(post.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find user"})
		}
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	fullPost.Author = user

	forum, err := h.ForumService.SelectFullForumBySlug(post.Forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find forum"})
		}
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	fullPost.Forum = forum

	thread, err := h.ThreadService.SelectThreadById(post.Thread)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	fullPost.Thread = thread

	return ctx.JSON(http.StatusOK, fullPost)
}

func (h *Post) EditMessage(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if id < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	editMessage := input.Message{}
	if err := ctx.Bind(&editMessage); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	num, err := h.PostService.UpdatePostMessage(editMessage.Message, id)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	if num != 1 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	post, err := h.PostService.SelectPostById(id)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	return ctx.JSON(http.StatusNotFound, post)
}
