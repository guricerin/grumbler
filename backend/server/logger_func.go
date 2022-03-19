package server

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
	"github.com/rs/zerolog"
)

func commonLog(e *zerolog.Event, c *gin.Context, statusCode int, user *model.User) *zerolog.Event {
	caller := ""
	_, fileName, line, ok := runtime.Caller(2)
	if ok {
		caller = fmt.Sprintf("%s:%d", fileName, line)
	}

	// todo: bodyを読み取れていない
	buf := make([]byte, 4096)
	len, err := c.Request.Body.Read(buf)
	if err != nil {
		len = 0
	}

	userId := ""
	userName := ""
	if user != nil {
		userId = user.Id
		userName = user.Name
	}

	return e.
		Str("method", c.Request.Method).
		Str("host", c.Request.Host).
		Str("url", c.Request.RequestURI).
		Str("client_ip", c.ClientIP()).
		Str("user_agent", c.Request.UserAgent()).
		Int64("content_length", c.Request.ContentLength).
		Str("content_type", c.Request.Header.Get("Content-Type")).
		Str("request_body", string(buf[0:len])).
		Str("caller", caller).
		Dict("user", zerolog.Dict().
			Str("id", userId).
			Str("name", userName)).
		Int("status_code", statusCode)
}

func (s *Server) Info(c *gin.Context, statusCode int, user *model.User, msg string) {
	commonLog(s.logger.Info(), c, statusCode, user).
		Msg(msg)
}

func (s *Server) Warn(c *gin.Context, statusCode int, user *model.User, err error) {
	commonLog(s.logger.Warn(), c, statusCode, user).
		Stack().
		Err(err).
		Send()
}

func (s *Server) Error(c *gin.Context, statusCode int, user *model.User, err error) {
	commonLog(s.logger.Error(), c, statusCode, user).
		Stack().
		Err(err).
		Send()
}
