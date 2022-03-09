package server

import (
	"errors"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func (s *Server) getSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		signinUser, err := s.fetchUserFromSession(c)
		if err != nil {
			// todo
			log.Printf("getSearch() 0: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		query := c.Query("q")
		if query == "" {
			// todo
			log.Printf("getSearch() 1: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(errors.New("検索対象の文字列が指定されていません。")))
			return
		}
		kind := c.Query("k")
		if kind == "" {
			// todo
			log.Printf("getSearch() 2: %s\n", err.Error())
			c.JSON(http.StatusBadRequest, errorRes(errors.New("検索種別が指定されていません。")))
			return
		}

		switch kind {
		case "user_id":
			users, err := s.userStore.SearchById(query)
			if err != nil {
				c.JSON(http.StatusBadRequest, errorRes(err))
				return
			}
			usersJson := make([]gin.H, 0)
			for _, u := range users {
				usersJson = append(usersJson, userRes(u))
			}
			c.JSON(http.StatusOK, gin.H{
				"users": usersJson,
			})
			return
		case "user_name":
			users, err := s.userStore.SearchByName(query)
			if err != nil {
				c.JSON(http.StatusBadRequest, errorRes(err))
				return
			}
			usersJson := make([]gin.H, 0)
			for _, u := range users {
				usersJson = append(usersJson, userRes(u))
			}
			c.JSON(http.StatusOK, gin.H{
				"users": usersJson,
			})
			return
		case "grumble":
			grumbles, err := s.grumbleStore.Search(signinUser.Id, query)
			if err != nil {
				c.JSON(http.StatusBadRequest, errorRes(err))
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
			return
		default:
			c.JSON(http.StatusBadRequest, errorRes(errors.New("不正な検索種別が指定されています。")))
			return
		}
	}
}
