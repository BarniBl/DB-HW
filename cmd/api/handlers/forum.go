package handlers

import (
	"database/sql"
	"fmt"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/labstack/echo"
	"math"
	"net/http"
	"strconv"
)

type Forum struct {
	ForumService  *forum.ForumService
	UserService   *forum.UserService
	ThreadService *forum.ThreadService
}

func (h *Forum) CreateForum(ctx echo.Context) (Err error) {

	newForum := forum.Forum{}
	if err := ctx.Bind(&newForum); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	fullForum, err := h.ForumService.SelectFullForumBySlug(newForum.Slug)
	if err == nil {
		return ctx.JSON(http.StatusConflict, fullForum)
	}
	if err != sql.ErrNoRows {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	nickName, err := h.UserService.FindUserByNickName(newForum.User)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	newForum.User = nickName

	if err = h.ForumService.InsertForum(newForum); err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	return ctx.JSON(http.StatusCreated, newForum)
}

func (h *Forum) CreateThread(ctx echo.Context) (Err error) {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	newThread := forum.Thread{}
	if err := ctx.Bind(&newThread); err != nil {
		fmt.Println(err)
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	newThread.Forum = slug

	threadForum, err := h.ForumService.SelectForumBySlug(newThread.Forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find forum"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	newThread.Forum = threadForum.Slug

	_, err = h.UserService.FindUserByNickName(newThread.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find user"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	if newThread.Slug != "" {
		thread, err := h.ThreadService.SelectThreadBySlug(newThread.Slug)
		if err == nil {
			return ctx.JSON(http.StatusConflict, thread)
		}
		if err != sql.ErrNoRows {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
		}
	}

	id, err := h.ThreadService.InsertThread(newThread)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	thread, err := h.ThreadService.SelectThreadById(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	return ctx.JSON(http.StatusCreated, thread)
}

func (h *Forum) GetForumDetails(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	fullForum, err := h.ForumService.SelectFullForumBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find forum"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	return ctx.JSON(http.StatusOK, fullForum)
}

func (h *Forum) GetForumThreads(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}

	limitStr := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	if limit < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	since := ctx.QueryParam("since")

	descStr := ctx.QueryParam("desc")
	desc, err := strconv.ParseBool(descStr)
	if err != nil {
		desc = false
	}
	_, err = h.ForumService.SelectForumBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	threads, err := h.ThreadService.SelectThreadByForum(slug, limit, since, desc)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	if len(threads) == 0 {
		threads := []forum.Thread{}
		return ctx.JSON(http.StatusOK, threads)
	}
	return ctx.JSON(http.StatusOK, threads)
}

func (h *Forum) GetForumUsers(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, forum.ErrorMessage{Message: "Error"})
	}
	_, err := h.ForumService.SelectForumBySlug(slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Can't find forum"})
	}
	limitStr := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = math.MaxInt32
	}
	if limit < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}

	since := ctx.QueryParam("since")

	descStr := ctx.QueryParam("desc")
	desc, err := strconv.ParseBool(descStr)
	if err != nil {
		desc = false
	}

	if desc == true {
		if since == "" {
			since = "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
		}
		users, err := h.UserService.SelectUsersByForumDesc(slug, limit, since)
		if err != nil {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		if len(users) == 0 {
			useres := []forum.User{}
			return ctx.JSON(http.StatusOK, useres)
		}
		return ctx.JSON(http.StatusOK, users)
	}
	if since == "" {
		users, err := h.UserService.SelectUsersByForumAntiSince(slug, limit)
		if err != nil {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
		}
		if len(users) == 0 {
			useres := []forum.User{}
			return ctx.JSON(http.StatusOK, useres)
		}
		return ctx.JSON(http.StatusOK, users)
	}
	users, err := h.UserService.SelectUsersByForum(slug, limit, since)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, forum.ErrorMessage{Message: "Error"})
	}
	if len(users) == 0 {
		useres := []forum.User{}
		return ctx.JSON(http.StatusOK, useres)
	}
	return ctx.JSON(http.StatusOK, users)
}

func (h *Forum) Clean(ctx echo.Context) error {
	err := h.ForumService.Clean()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, nil)
}

func (h *Forum) Status(ctx echo.Context) error {

	status, err := h.ForumService.SelectStatus()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, status)
}
