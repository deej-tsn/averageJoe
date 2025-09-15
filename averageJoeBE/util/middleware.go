package util

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func ExtractWebsocketToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Request().Header["Sec-Websocket-Protocol"]
		if ok {
			// should be in format auth-token.<token>
			tokenString := strings.Split(token[0], "-")[1]
			c.Request().Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
		}
		return next(c)

	}
}
