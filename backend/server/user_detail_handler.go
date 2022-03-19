package server

import (
	"errors"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type followReq struct {
	DstUserId string `json:"dstUserId"`
}

func followRes(follow model.Follow) gin.H {
	return gin.H{
		"srcUserId": follow.SrcUserId,
		"dstUserId": follow.DstUserId,
	}
}

type userSettingsReq struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

func userDetailRes(user model.User, grumbles []model.GrumbleRes, follows []model.User, followers []model.User, bookmarks []model.GrumbleRes, isFollow bool, isFollower bool) gin.H {
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
	bookmarksJson := make([]gin.H, 0)
	for _, g := range bookmarks {
		bookmarksJson = append(bookmarksJson, grumbleRes(g))
	}

	return gin.H{
		"user":       userRes(user),
		"grumbles":   grumblesJson,
		"follows":    followsJson,
		"followers":  followersJson,
		"bookmarks":  bookmarksJson,
		"isFollow":   isFollow,
		"isFollower": isFollower,
	}
}

func (s *Server) getUserDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}

		userId := c.Param("id")
		user, err := s.userStore.RetrieveById(userId)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		grumbles, err := s.grumbleStore.RetrieveByUserId(signinUser.Id, user.Id)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}
		model.SortGrumblesForNewest(grumbles)

		follows, err := s.followStore.RetrieveFollows(user.Id)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}
		followUsers := make([]model.User, 0)
		for _, f := range follows {
			u, err := s.userStore.RetrieveById(f.DstUserId)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, &signinUser, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			followUsers = append(followUsers, u)
		}

		followers, err := s.followStore.RetrieveFollowers(user.Id)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}
		followerUsers := make([]model.User, 0)
		for _, f := range followers {
			u, err := s.userStore.RetrieveById(f.SrcUserId)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, &signinUser, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			followerUsers = append(followerUsers, u)
		}

		bookmarks, err := s.grumbleStore.RetrieveBookmarkedGrumblesByUserId(signinUser.Id, user.Id)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}
		// 最新日時順
		sort.Slice(bookmarks, func(i, j int) bool {
			return bookmarks[i].CreatedAt.After(bookmarks[j].CreatedAt)
		})

		isFollow, isFollower, err := s.followStore.RetrieveFollowRelation(signinUser.Id, userId)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &signinUser, "success to get user detail")
		c.JSON(http.StatusOK, userDetailRes(user, grumbles, followUsers, followerUsers, bookmarks, isFollow, isFollower))
	}
}

func (s *Server) postFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}

		var req followReq
		if err := c.BindJSON(&req); err != nil {
			s.Warn(c, http.StatusBadRequest, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := s.followStore.Create(signinUser.Id, req.DstUserId); err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &signinUser, "success to follow")
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
		return
	}
}

func (s *Server) postUnFollow() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}

		var req followReq
		if err := c.BindJSON(&req); err != nil {
			s.Warn(c, http.StatusBadRequest, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := s.followStore.Delete(signinUser.Id, req.DstUserId); err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &signinUser, "success to unfollow")
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
		return
	}
}

func (s *Server) postUserSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}
		var req userSettingsReq
		if err := c.BindJSON(&req); err != nil {
			s.Warn(c, http.StatusBadRequest, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if err := model.ValidateUserName(req.Name); err != nil {
			s.Warn(c, http.StatusBadRequest, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateUserProfile(req.Profile); err != nil {
			s.Warn(c, http.StatusBadRequest, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		signinUser.Name = req.Name
		signinUser.Profile = req.Profile
		if err := s.userStore.Update(&signinUser); err != nil {
			s.Error(c, http.StatusInternalServerError, &signinUser, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &signinUser, "success to change user setting")
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
		return
	}
}
