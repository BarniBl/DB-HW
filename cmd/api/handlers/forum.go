package handlers

import (
	"database/sql"
	"fmt"
	"github.com/BarniBl/DB-HW/internal/forum"
	"github.com/BarniBl/DB-HW/internal/input"
	"github.com/BarniBl/DB-HW/internal/output"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
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

	forum, err := h.ForumService.SelectFullForumBySlug(newForum.Slug)
	if err == nil {
		return ctx.JSON(http.StatusConflict, forum)
	}
	if err != sql.ErrNoRows {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	_, err = h.UserService.FindUserByNickName(newForum.User)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find user"})
		}
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

	_, err := h.ForumService.SelectForumBySlug(newThread.Forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find forum"})
		}
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	_, err = h.UserService.FindUserByNickName(newThread.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find user"})
		}
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	thread, err := h.ThreadService.SelectThreadByTitle(newThread.Title)
	if err == nil {
		return ctx.JSON(http.StatusConflict, thread)
	}
	if err != sql.ErrNoRows {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	err = h.ThreadService.InsertThread(newThread)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	thread, err = h.ThreadService.SelectThreadByTitle(newThread.Title)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	return ctx.JSON(http.StatusCreated, thread)
}

func (h *Forum) GetForumDetails(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	forum, err := h.ForumService.SelectFullForumBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Can't find forum"})
		}
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	return ctx.JSON(http.StatusOK, forum)
}

func (h *Forum) GetForumThreads(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	limitStr := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if limit < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	sinceStr := ctx.QueryParam("since")
	since, err := strconv.Atoi(sinceStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if since < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	descStr := ctx.QueryParam("desc")
	desc, err := strconv.ParseBool(descStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	if desc == true {
		threads, err := h.ThreadService.SelectThreadByForumDesc(slug, limit, since)
		if err != nil {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
		}
		if len(threads) == 0 {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
		}
		return ctx.JSON(http.StatusOK, threads)
	}
	threads, err := h.ThreadService.SelectThreadByForum(slug, limit, since)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if len(threads) == 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	return ctx.JSON(http.StatusOK, threads)
}

func (h *Forum) GetForumUsers(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.JSON(http.StatusBadRequest, output.ErrorMessage{Message: "Error"})
	}

	limitStr := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if limit < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	sinceStr := ctx.QueryParam("since")
	since, err := strconv.Atoi(sinceStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if since < 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	descStr := ctx.QueryParam("desc")
	desc, err := strconv.ParseBool(descStr)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}

	if desc == true {
		users, err := h.UserService.SelectUsersByForumDesc(slug, limit, since)
		if err != nil {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
		}
		if len(users) == 0 {
			ctx.Logger().Warn(err)
			return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
		}
		return ctx.JSON(http.StatusOK, users)
	}
	users, err := h.UserService.SelectUsersByForum(slug, limit, since)
	if err != nil {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	if len(users) == 0 {
		ctx.Logger().Warn(err)
		return ctx.JSON(http.StatusNotFound, output.ErrorMessage{Message: "Error"})
	}
	return ctx.JSON(http.StatusOK, users)
}

func (h *Forum) Clean(ctx echo.Context) (Err error) {
	defer func() {
		if bodyErr := ctx.Request().Body.Close(); bodyErr != nil {
			Err = errors.Wrap(Err, bodyErr.Error())
		}
	}()

	ctx.Response().Header().Set("Content-Type", "application/json")

	err := h.ForumService.Clean()
	if err != nil {
		return err
	}


	if err := ctx.JSON(200, nil); err != nil {
		return err
	}

	return nil
}