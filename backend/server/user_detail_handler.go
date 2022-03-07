package server

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type followReq struct {
	SrcUserId string `josn:"srcUserId"`
	DstUserId string `json:"dstUserId"`
}

func followRes(follow model.Follow) gin.H {
	return gin.H{
		"srcUserId": follow.SrcUserId,
		"dstUserId": follow.DstUserId,
	}
}

func userDetailRes(user model.User, grumbles []model.GrumbleRes, follows []model.User, followers []model.User, isFollow bool, isFollower bool) gin.H {
	grumblesJson := make([]gin.H, 0)
	for _, g := range grumbles {
		grumblesJson = append(grumblesJson, grumbleRes(g))
	}
	followsJson := make([]gin.H, 0)
	for _, f := range follows {
		followsJson = append(followsJson, userRes(f))
	}
	followersJson := make([]gin.H, 0)
	for _, f := range followers {
		followersJson = append(followersJson, userRes(f))
	}

	return gin.H{
		"user":       userRes(user),
		"grumbles":   grumblesJson,
		"follows":    followsJson,
		"followers":  followersJson,
		"isFollow":   isFollow,
		"isFollower": isFollower,
	}
}

func (s *Server) getUserDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			log.Printf("getUserDetail() 0: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		userId := c.Param("id")
		user, err := s.userStore.RetrieveById(userId)
		if err != nil {
			// todo
			log.Printf("getUserDetail() 1: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		grumbles, err := s.grumbleStore.RetrieveByUserId(user.Id)
		if err != nil {
			// todo
			log.Printf("getUserDetail() 2: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		// 最新日時順
		sort.Slice(grumbles, func(i, j int) bool {
			return grumbles[i].CreatedAt.After(grumbles[j].CreatedAt)
		})

		follows, err := s.followStore.RetrieveFollows(user.Id)
		if err != nil {
			// todo
			log.Printf("getUserDetail() 3: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		followUsers := make([]model.User, 0)
		for _, f := range follows {
			u, err := s.userStore.RetrieveById(f.DstUserId)
			if err != nil {
				// todo
				log.Printf("getUserDetail() 4: %s\n", err.Error())
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			followUsers = append(followUsers, u)
		}

		followers, err := s.followStore.RetrieveFollowers(user.Id)
		if err != nil {
			// todo
			log.Printf("getUserDetail() 5: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		followerUsers := make([]model.User, 0)
		for _, f := range followers {
			u, err := s.userStore.RetrieveById(f.SrcUserId)
			if err != nil {
				// todo
				log.Printf("getUserDetail() 6: %s\n", err.Error())
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			followerUsers = append(followerUsers, u)
		}

		isFollow, isFollower, err := s.followStore.RetrieveFollowRelation(signinUser.Id, userId)
		if err != nil {
			// todo
			log.Printf("getUserDetail() 7: %s\n", err.Error())
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		c.JSON(http.StatusOK, userDetailRes(user, grumbles, followUsers, followerUsers, isFollow, isFollower))
	}
}

func (s *Server) postFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req followReq
		if err := c.BindJSON(&req); err != nil {
			log.Printf("postFollow(): %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := s.followStore.Create(req.SrcUserId, req.DstUserId); err != nil {
			log.Printf("postFollow(): %s\n", err.Error())
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

func (s *Server) postUnFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req followReq
		if err := c.BindJSON(&req); err != nil {
			log.Printf("postUnFollow(): %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := s.followStore.Delete(req.SrcUserId, req.DstUserId); err != nil {
			log.Printf("postUnFollow(): %s\n", err.Error())
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
