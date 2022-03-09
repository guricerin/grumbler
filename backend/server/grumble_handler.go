package server

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type postGrumbleReq struct {
	Content string `json:"content"`
}

func grumbleRes(g model.GrumbleRes) gin.H {
	return gin.H{
		"pk":        g.Pk,
		"content":   g.Content,
		"userId":    g.UserId,
		"userName":  g.UserName,
		"createdAt": g.CreatedAt.Format("2006/01/02 15:04:05"),
	}
}

type bookmarkReq struct {
	GrumbleOk string `json:"grumblePk"`
	ByUserId  string `json:"byUserId"`
}

func (s *Server) getTimeline() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		grumbles, err := s.grumbleStore.RetrieveByUserId(user.Id)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		follows, err := s.followStore.RetrieveFollows(user.Id)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		for _, f := range follows {
			gs, err := s.grumbleStore.RetrieveByUserId(f.DstUserId)
			if err != nil {
				// todo
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			grumbles = append(grumbles, gs...)
		}
		// 最新日時順
		sort.Slice(grumbles, func(i, j int) bool {
			return grumbles[i].CreatedAt.After(grumbles[j].CreatedAt)
		})

		grumblesJson := make([]gin.H, 0)
		for _, g := range grumbles {
			grumblesJson = append(grumblesJson, grumbleRes(g))
		}
		c.JSON(http.StatusOK, gin.H{
			"grumbles": grumblesJson,
		})
	}
}

func (s *Server) getUserGrumbles() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		grumbles, err := s.grumbleStore.RetrieveByUserId(userId)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		// 最新日時順
		sort.Slice(grumbles, func(i, j int) bool {
			return grumbles[i].CreatedAt.After(grumbles[j].CreatedAt)
		})

		grumblesJson := make([]gin.H, 0)
		for _, g := range grumbles {
			grumblesJson = append(grumblesJson, grumbleRes(g))
		}
		c.JSON(http.StatusOK, gin.H{
			"grumbles": grumblesJson,
		})
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

func (s *Server) postBookmark() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req bookmarkReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if _, err := s.grumbleStore.CreateBookmark(req.GrumbleOk, req.ByUserId); err != nil {
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
