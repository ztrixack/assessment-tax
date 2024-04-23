package middlewares

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

func Logger(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}

			fields := logger.Fields{
				"remote_ip":  c.RealIP(),
				"host":       req.Host,
				"method":     req.Method,
				"uri":        req.RequestURI,
				"user_agent": req.UserAgent(),
				"status":     res.Status,
				"referer":    req.Referer(),
				"latency":    stop.Sub(start).String(),
				"bytes_in":   cl,
				"bytes_out":  strconv.FormatInt(res.Size, 10),
			}

			if err != nil {
				log.Err(err).Fields(fields).E("Request failed")
			} else {
				log.Fields(fields).D("Request completed")
			}

			return err
		}
	}
}
