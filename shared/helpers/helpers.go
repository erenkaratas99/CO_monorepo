package helpers

import "github.com/labstack/echo/v4"

func WithTags(tags ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("swagger_tag", tags)
			return next(c)
		}
	}
}
