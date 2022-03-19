package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type searchReq struct {
	Keyword string `json:"keyword"`
	Kind    string `json:"kind"`
}

func (s *Server) getSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}

		var req searchReq
		if err := c.BindJSON(&req); err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		switch req.Kind {
		case "user_id":
			users, err := s.userStore.Search(req.Keyword, UserIdSearch)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, &signinUser, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			usersJson := make([]gin.H, 0)
			for _, u := range users {
				usersJson = append(usersJson, userRes(u))
			}

			s.Info(c, http.StatusOK, &signinUser, "success to search user_id")
			c.JSON(http.StatusOK, gin.H{
				"users": usersJson,
			})
			return
		case "user_name":
			users, err := s.userStore.Search(req.Keyword, UserNameSearch)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, &signinUser, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			usersJson := make([]gin.H, 0)
			for _, u := range users {
				usersJson = append(usersJson, userRes(u))
			}

			s.Info(c, http.StatusOK, &signinUser, "success to search user_name")
			c.JSON(http.StatusOK, gin.H{
				"users": usersJson,
			})
			return
		case "grumble":
			grumbles, err := s.grumbleStore.Search(signinUser.Id, req.Keyword)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, &signinUser, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			model.SortGrumblesForNewest(grumbles)

			grumblesJson := make([]gin.H, 0)
			for _, g := range grumbles {
				grumblesJson = append(grumblesJson, grumbleRes(g))
			}

			s.Info(c, http.StatusOK, &signinUser, "success to search grumble")
			c.JSON(http.StatusOK, gin.H{
				"grumbles": grumblesJson,
			})
			return
		default:
			s.Warn(c, http.StatusBadRequest, &signinUser, err)
			c.JSON(http.StatusBadRequest, errorRes(errors.New("不正な検索種別が指定されています。")))
			return
		}
	}
}
