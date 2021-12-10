package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type signupUserReq struct {
	Id       string `json:"id" binding:"required, alphanum, min=1, max=255"`
	Name     string `json:"name" binding:"required, min=1"`
	Password string `json:"password" binding:"required, min=8, max=255"`
}

type loginUserReq struct {
	Id       string `json:"id" binding:"required, alphanum, min=1, max=255"`
	Password string `json:"password" binding:"required, min=8, max=255"`
}

// POST /signup
func (s *Server) postSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signupUserReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		_, err := s.userStore.RetrieveById(req.Id)
		if err != nil && err == sql.ErrNoRows {
			// ok
			hashedPassword, err := encryptPassword(req.Password)
			if err != nil {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			signupUser := model.User{
				Id:       req.Id,
				Name:     req.Name,
				Password: hashedPassword,
			}
			err = s.userStore.Create(signupUser)
			if err != nil {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}

			token, err := createUuid()
			if err != nil {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			err = s.sessionStore.Create(token, signupUser)
			if err != nil {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
			session := sessions.Default(c)
			session.Set(SESSION_TOKEN, token)
			session.Save()

			url := fmt.Sprintf("/user/%s", signupUser.Id)
			c.Redirect(http.StatusFound, url)
		} else if err != nil && err != sql.ErrNoRows {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
		} else {
			// duplicate id
			msg := fmt.Sprintf("the user id '%s' is already used.", req.Id)
			e := errors.New(msg)
			c.JSON(http.StatusBadRequest, errorRes(e))
		}
	}
}

// POST /login
func (s *Server) postLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginUserReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := s.userStore.RetrieveById(req.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				err := errors.New("id or password is wrong.")
				c.JSON(http.StatusBadRequest, errorRes(err))
				return
			} else {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				return
			}
		}

		if !verifyPasswordHash(user.Password, req.Password) {
			err := errors.New("id or password is wrong.")
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		token, err := createUuid()
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		err = s.sessionStore.Create(token, user)
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		session := sessions.Default(c)
		session.Set(SESSION_TOKEN, token)
		session.Save()

		url := fmt.Sprintf("/user/%s", req.Id)
		c.Redirect(http.StatusFound, url)
	}
}
