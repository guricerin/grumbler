package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type postGrumbleReq struct {
	UserId  string `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=300"`
}

func (s *Server) postGrumble() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req postGrumbleReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := s.authorizationCheck(c)
		if err != nil {
			c.JSON(http.StatusForbidden, errorRes(err))
			return
		}
		err = s.grumbleStore.Create(req.Content, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
		return
	}
}
