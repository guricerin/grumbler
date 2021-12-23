package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

func (s *Server) getSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, errorRes(errors.New("検索対象のワードが指定されていません。")))
			return
		}
		kind := c.Query("k")
		if kind == "" {
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
			usersJson := make([]model.UserJson, 0)
			for _, u := range users {
				usersJson = append(usersJson, u.ToJson())
			}
			c.JSON(http.StatusOK, gin.H{
				"users": usersJson,
			})
			return
		case "user_name":
			// todo
			return
		case "grumble":
			// todo
			return
		default:
			c.JSON(http.StatusBadRequest, errorRes(errors.New("不正な検索種別が指定されています。")))
			return
		}
	}
}
