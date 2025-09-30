package middleware

import (
	"server/internal/pkg/domain"

	"github.com/labstack/echo"
)

func AuthEchoMiddleware(service domain.SessionService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			_, err := service.CheckSession(context.Request().Header)
			if err != nil {
				return context.NoContent(401)
			}

			return next(context)
		}
	}
}
