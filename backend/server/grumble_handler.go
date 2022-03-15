package server

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type postGrumbleReq struct {
	Content string `json:"content"`
}

func replyInfoForGrumbleResRes(r model.ReplyInfoFoGrumbleRes) gin.H {
	return gin.H{
		"dstGrumblePk": r.DstGrumblePk,
		"dstUserId":    r.DstUserId,
		"repliedCount": r.RepliedCount,
	}
}

func grumbleRes(g model.GrumbleRes) gin.H {
	return gin.H{
		"pk":                       g.Pk,
		"content":                  g.Content,
		"userId":                   g.UserId,
		"userName":                 g.UserName,
		"createdAt":                g.CreatedAt.Format("2006/01/02 15:04:05"),
		"reply":                    replyInfoForGrumbleResRes(g.Reply),
		"bookmarkedCount":          g.BookmarkedCount,
		"isBookmarkedBySigninUser": g.IsBookmarkedBySigninUser,
	}
}

type bookmarkReq struct {
	GrumblePk string `json:"grumblePk"`
}

func grumbleDetailRes(mainGrumble model.GrumbleRes, ancestors []model.GrumbleRes, replies []model.GrumbleRes) gin.H {
	ancestorsJson := make([]gin.H, 0)
	for _, r := range ancestors {
		ancestorsJson = append(ancestorsJson, grumbleRes(r))
	}
	repliesJson := make([]gin.H, 0)
	for _, r := range replies {
		repliesJson = append(repliesJson, grumbleRes(r))
	}
	return gin.H{
		"target":    grumbleRes(mainGrumble),
		"ancestors": ancestorsJson,
		"replies":   repliesJson,
	}
}

type postReplyReq struct {
	Content      string `json:"content"`
	DstGrumblePk string `json:"dstGrumblePk"`
}

func (s *Server) getGrumbleDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		grumblePk := c.Param("grumble_pk")
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		mainGrumble, err := s.grumbleStore.RetrieveByPk(grumblePk, user.Id)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ancestors, err := s.grumbleStore.RetrieveReplyAncestors(grumblePk, user.Id)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		// 最古日時順
		sort.Slice(ancestors, func(i, j int) bool {
			return ancestors[i].CreatedAt.Before(ancestors[j].CreatedAt)
		})
		replies, err := s.grumbleStore.RetrieveByReplyDstPk(grumblePk, user.Id)
		if err != nil {
			// todo
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		// 最古日時順
		sort.Slice(replies, func(i, j int) bool {
			return replies[i].CreatedAt.Before(replies[j].CreatedAt)
		})

		c.JSON(http.StatusOK, grumbleDetailRes(mainGrumble, ancestors, replies))
		return
	}
}

func (s *Server) getTimeline() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		grumbles, err := s.grumbleStore.RetrieveByUserId(user.Id, user.Id)
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
			gs, err := s.grumbleStore.RetrieveByUserId(user.Id, f.DstUserId)
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

func (s *Server) postGrumble() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		var req postGrumbleReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateGrumble(req.Content); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		_, err = s.grumbleStore.Create(req.Content, user)
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

type deleteGrumbleReq struct {
	GrumblePk string `json:"grumblePk"`
}

func (s *Server) postDeleteGrumble() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		var req deleteGrumbleReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		err = s.grumbleStore.Delete(req.GrumblePk, user.Id)
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
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		var req bookmarkReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if _, err := s.grumbleStore.CreateBookmark(req.GrumblePk, signinUser.Id); err != nil {
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

func (s *Server) postDeleteBookmark() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		var req bookmarkReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if err := s.grumbleStore.DeleteBookmark(req.GrumblePk, signinUser.Id); err != nil {
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

func (s *Server) postReply() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			log.Printf("postReply() 2: %s\n", err.Error())
			c.JSON(http.StatusUnauthorized, errorRes(err))
			return
		}

		var req postReplyReq
		if err := c.BindJSON(&req); err != nil {
			// todo
			log.Printf("postReply() 0: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateGrumble(req.Content); err != nil {
			// todo
			log.Printf("postReply() 1: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		grumble, err := s.grumbleStore.Create(req.Content, user)
		if err != nil {
			// todo
			log.Printf("postReply() 3: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		_, err = s.grumbleStore.CreateReply(grumble.Pk, req.DstGrumblePk)
		if err != nil {
			// todo
			log.Printf("postReply() 4: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
		return
	}
}
