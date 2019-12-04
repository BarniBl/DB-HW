package serve

import "github.com/labstack/echo"

func JSONData(ctx echo.Context, statusCode int, data interface{}) {
	ctx.Response().Header().Set("Content-Type", "application/json")
}

func JSONError(ctx echo.Context, statusCode int, textError string) {

}
