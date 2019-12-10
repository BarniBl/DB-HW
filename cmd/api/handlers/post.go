package handlers

import (
	"github.com/BarniBl/DB-HW/internal/forum"
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

	postSlice, err := h.PostService.SelectPostById(id)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if len(postSlice) == 0 {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	post := postSlice[0]

	fullPost := output.FullPost{Post: post}

	userSlice, err := h.UserService.SelectUserByNickName(post.Author)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if len(userSlice) == 0 {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	fullPost.Author = userSlice[0]

	forumSlice, err := h.ForumService.SelectForumBySlug(post.Forum)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if len(forumSlice) == 0 {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	fullPost.Forum = forumSlice[0]

	threadSlice, err := h.ThreadService.SelectThreadById(post.Thread)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if len(forumSlice) == 0 {
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	fullPost.Thread = threadSlice[0]

	return ctx.JSON(http.StatusOK, fullPost)
}
