package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type getUserReq struct {
	Id string `json:"id"`
}

type signupUserReq struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type signinUserReq struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func (s *Server) signinCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}

		err = s.resetSessToken(c)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, nil, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &user, "ok")
		c.JSON(http.StatusOK, userRes(user))
	}
}

func (s *Server) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		user, err := s.userStore.RetrieveById(userId)
		if err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(errors.New("bad request")))
			return
		}

		s.Info(c, http.StatusOK, &user, "ok")
		c.JSON(http.StatusOK, userRes(user))
	}
}

func (s *Server) postSignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signupUserReq
		if err := c.BindJSON(&req); err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateUserName(req.Name); err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateUserId(req.Id); err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := model.ValidateUserPassword(req.Password); err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		_, err := s.userStore.RetrieveById(req.Id)
		if err != nil && err == sql.ErrNoRows {
			// idがだぶってないのでok
			hashedPassword, err := encryptPassword(req.Password)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			signupUser := model.User{
				Id:       req.Id,
				Name:     req.Name,
				Password: hashedPassword,
				Profile:  "", // 後から設定させる
			}
			err = s.userStore.Create(&signupUser)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}

			token, err := createUuid()
			if err != nil {
				s.Error(c, http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			err = s.sessionStore.Create(token, signupUser)
			if err != nil {
				s.Error(c, http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
			s.setCookie(c, token)

			s.Info(c, http.StatusOK, &signupUser, "success to signup")
			c.JSON(http.StatusOK, gin.H{
				"id":      signupUser.Id,
				"name":    signupUser.Name,
				"profile": "",
			})
		} else if err != nil && err != sql.ErrNoRows {
			s.Error(c, http.StatusInternalServerError, nil, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
		} else {
			// duplicate id
			msg := fmt.Sprintf("ユーザID '%s' は既に使用されています。", req.Id)
			e := errors.New(msg)
			s.Warn(c, http.StatusBadRequest, nil, e)
			c.JSON(http.StatusBadRequest, errorRes(e))
		}
	}
}

func (s *Server) postSignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signinUserReq
		if err := c.BindJSON(&req); err != nil {
			s.Warn(c, http.StatusBadRequest, nil, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := s.userStore.RetrieveById(req.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				err := errors.New("ユーザIDまたはパスワードが異なります。")
				s.Warn(c, http.StatusBadRequest, nil, err)
				c.JSON(http.StatusBadRequest, errorRes(err))
				return
			} else {
				s.Error(c, http.StatusInternalServerError, nil, err)
				c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
				return
			}
		}

		if !verifyPasswordHash(user.Password, req.Password) {
			err := errors.New("ユーザIDまたはパスワードが異なります。")
			s.Warn(c, http.StatusBadRequest, &user, err)
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		token, err := createUuid()
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &user, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		if err = s.sessionStore.Create(token, user); err != nil {
			s.Error(c, http.StatusInternalServerError, &user, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.setCookie(c, token)
		s.Info(c, http.StatusOK, &user, "success to signin")
		c.JSON(http.StatusOK, userRes(user))
	}
}

func (s *Server) postSignOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}
		err = s.deleteCookie(c)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &user, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		if err := s.sessionStore.DeleteByUserPk(user.Pk); err != nil {
			s.Error(c, http.StatusInternalServerError, &user, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &user, "success to signout")
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	}
}

func (s *Server) postUnsubscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := s.fetchUserFromSession(c)
		if err != nil {
			s.Warn(c, http.StatusUnauthorized, nil, err)
			c.JSON(http.StatusUnauthorized, errorRes(errors.New("unauthorized")))
			return
		}
		err = s.deleteCookie(c)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &user, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		err = s.userStore.DeleteByPk(user.Pk)
		if err != nil {
			s.Error(c, http.StatusInternalServerError, &user, err)
			c.JSON(http.StatusInternalServerError, errorRes(errors.New("server error")))
			return
		}

		s.Info(c, http.StatusOK, &user, "success to unsubscribe")
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	}
}
