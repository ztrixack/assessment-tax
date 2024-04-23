package middlewares

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

func BasicAuth(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, password, ok := c.Request().BasicAuth()
			if !ok || !isValidCredentials(username, password) {
				log.Fields(logger.Fields{"username": username, "password": securePassword(password)}).E("Invalid username or password")
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

func isValidCredentials(username, password string) bool {
	if username != os.Getenv("ADMIN_USERNAME") || password != os.Getenv("ADMIN_PASSWORD") {
		return false
	}

	return true
}

func securePassword(pwd string) string {
	if len(pwd) < 6 {
		return "******"
	}

	return pwd[0:4] + "****"
}
