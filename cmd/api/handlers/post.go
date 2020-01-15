package handlers

import (
	"database/sql"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
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

	related := ctx.QueryParam("related")

	post, err := h.PostService.SelectPostById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find post"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	fullPost := forum.FullPost{Post: post}

	if strings.Contains(related, "user") {
		user, err := h.UserService.SelectUserByNickName(post.Author)
		if err != nil {
			if err == sql.ErrNoRows {
				return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
			}
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		fullPost.Author = user
	}

	if strings.Contains(related, "forum") {
		fullForum, err := h.ForumService.SelectFullForumBySlug(post.Forum)
		if err != nil {
			if err == sql.ErrNoRows {
				return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find forum"})
			}
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		fullPost.Forum = fullForum
	}

	if strings.Contains(related, "thread") {
		thread, err := h.ThreadService.SelectThreadById(post.Thread)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, "")
		}
		fullPost.Thread = thread
	}

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
	post, err := h.PostService.SelectPostById(id)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	if editMessage.Message != "" && editMessage.Message != post.Message {
		num, err := h.PostService.UpdatePostMessage(editMessage.Message, id)
		if err != nil {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		if num != 1 {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find post"})
		}
		post.Message = editMessage.Message
		post.IsEdited = true
	}

	return ctx.JSON(http.StatusOK, post)
}

func (h *Post) CreatePosts(ctx echo.Context) error {
	createdTime := time.Now().Format(time.RFC3339Nano)
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
	for i := 0; i < len(newPosts); i++ {
		newPosts[i].Thread = thread.Id
		newPosts[i].Forum = thread.Forum
		newPosts[i].Created = createdTime
		_, err := h.UserService.FindUserByNickName(newPosts[i].Author)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
		}
		if newPosts[i].Parent != 0 {
			err = h.PostService.FindPostById(newPosts[i].Parent, newPosts[i].Thread)
			if err != nil {
				return ctx.JSON(http.StatusConflict, forum.ErrorMessage{Message: "Can't find post"})
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

func (h *Post) EditThread(ctx echo.Context) error {
	slugOrIdStr := ctx.Param("slug_or_id")
	var editThread forum.Thread
	if err := ctx.Bind(&editThread); err != nil {
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
	if editThread.Message != "" {
		thread.Message = editThread.Message
	}
	if editThread.Title != "" {
		thread.Title = editThread.Title
	}
	if editThread.Message == "" && editThread.Title == "" {
		return ctx.JSON(http.StatusOK, thread)
	}
	err = h.ThreadService.UpdateThread(thread)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't update thread"})
	}
	return ctx.JSON(http.StatusOK, thread)
}
func (h *Post) CreateVote(ctx echo.Context) error {
	slugOrIdStr := ctx.Param("slug_or_id")
	var newVote forum.Vote
	if err := ctx.Bind(&newVote); err != nil {
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

	_, err = h.UserService.FindUserByNickName(newVote.NickName)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
	}
	newVote.ThreadId = thread.Id
	voted, err := h.ThreadService.FindVote(newVote)
	if voted {
		err = h.ThreadService.UpdateVote(newVote)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't vote"})
		}
	} else {
		err = h.ThreadService.InsertVote(newVote)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't vote"})
		}
	}

	thread, err = h.ThreadService.SelectThreadById(newVote.ThreadId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
	}
	return ctx.JSON(http.StatusOK, thread)
}

func (h *Post) GetThread(ctx echo.Context) error {
	slugOrIdStr := ctx.Param("slug_or_id")

	var thread forum.Thread
	id, err := strconv.Atoi(slugOrIdStr)
	if err != nil {
		slug := slugOrIdStr
		thread, err = h.ThreadService.SelectThreadBySlug(slug)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
		}
	} else {
		thread, err = h.ThreadService.SelectThreadById(id)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
		}
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (h *Post) GetPosts(ctx echo.Context) error {
	slugOrIdStr := ctx.Param("slug_or_id")

	limit := ctx.QueryParam("limit")
	since := ctx.QueryParam("since")
	sort := ctx.QueryParam("sort")
	desc := ctx.QueryParam("desc")

	if limit == "" {
		limit = "100"
	}

	if sort == "" {
		sort = "flat"
	}
	if desc == "" {
		desc = "false"
	}
	if since == "" {
		if desc == "false" {
			since = "0"
		} else {
			since = "999999999"
		}
	}

	var thread forum.Thread
	id, err := strconv.Atoi(slugOrIdStr)
	if err != nil {
		thread, err = h.ThreadService.SelectThreadBySlug(slugOrIdStr)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
		}
	} else {
		thread, err = h.ThreadService.SelectThreadById(id)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find thread"})
		}
	}

	posts, err := h.ThreadService.SelectPosts(thread.Id, limit, since, sort, desc)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't read posts"})
	}
	if len(posts) == 0 {
		postss := []Post{}
		return ctx.JSON(http.StatusOK, postss)
	}
	return ctx.JSON(http.StatusOK, posts)
}
