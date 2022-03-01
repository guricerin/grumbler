package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type postGrumbleReq struct {
	Content string `json:"content"`
}

type getGrumblesReq struct {
	UserId string `json:"user_id" binding:"required"`
}

func grumbleRes(g model.Grumble) gin.H {
	return gin.H{
		"id":         g.Pk,
		"content":    g.Content,
		"user_id":    g.UserId,
		"created_at": g.CreatedAt,
	}
}

func (s *Server) postGrumble() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req postGrumbleReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateGrumble(req.Content); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		err = s.grumbleStore.Create(req.Content, user)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
		return
	}
}

func (s *Server) getGrumbles() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getGrumblesReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		grumbles, err := s.grumbleStore.RetrieveByUserId(req.UserId)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		grumblesJson := make([]gin.H, 0)
		for _, g := range grumbles {
			grumblesJson = append(grumblesJson, grumbleRes(g))
		}
		c.JSON(http.StatusOK, gin.H{
			"grumbles": grumblesJson,
		})
	}
}
