package handlers

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
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
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	if id < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	post, err := h.PostService.SelectPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find post"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	fullPost := forum.FullPost{Post: post}

	user, err := h.UserService.SelectUserByNickName(post.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	fullPost.Author = user

	fullForum, err := h.ForumService.SelectFullForumBySlug(post.Forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find forum"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	fullPost.Forum = fullForum

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
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	if id < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	editMessage := forum.Message{}
	if err := ctx.Bind(&editMessage); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	num, err := h.PostService.UpdatePostMessage(editMessage.Message, id)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	if num != 1 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	post, err := h.PostService.SelectPostById(id)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	return ctx.JSON(http.StatusNotFound, post)
}

func (h *Post) CreatePosts(ctx echo.Context) error {
	slugOrIdStr := ctx.Param("slug_or_id")
	var newPosts []forum.Post
	if err := ctx.Bind(&newPosts); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}
	var thread forum.Thread
	id, err := strconv.Atoi(slugOrIdStr)
	if err != nil {
		slug := slugOrIdStr
		thread, err = h.ThreadService.FindThreadBySlug(slug)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
		}
	} else {
		thread, err = h.ThreadService.FindThreadById(id)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
		}
	}
	createdTime := time.Now()
	for i := 0; i < len(newPosts); i++ {
		newPosts[i].Thread = thread.Id
		newPosts[i].Forum = thread.Forum
		newPosts[i].Created = createdTime
		_, err := h.UserService.FindUserByNickName(newPosts[i].Author)
		if err != nil {
			return ctx.JSON(http.StatusConflict, forum.ErrorMessage{Message: "Can't find user"})
		}
		if newPosts[i].Parent != 0 {
			_, err = h.PostService.SelectPostById(newPosts[i].Parent)
			if err != nil {
				return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find post"})
			}
		}
		lastId, err := h.PostService.InsertPost(newPosts[i])
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		newPosts[i].Id = lastId
		newPosts[i].IsEdited = false
	}
	return ctx.JSON(http.StatusCreated, newPosts)
}
