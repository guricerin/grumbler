package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
)

type getUserReq struct {
	Id string `json:"id" binding:"required,alphanum,min=1,max=255"`
}

type signupUserReq struct {
	Id       string `json:"id" binding:"required,alphanum,min=1,max=255"`
	Name     string `json:"name" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type signinUserReq struct {
	Id       string `json:"id" binding:"required,alphanum,min=1,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

func (s *Server) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req getUserReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := s.userStore.RetrieveById(req.Id)
		if err != nil {
			// todo
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   user.Id,
			"name": user.Name,
		})
	}
}

// ページリソースのユーザと、それにアクセスしようとしているユーザは同一か
func (s *Server) authorizationCheck(c *gin.Context) (bool, error) {
	userId := c.Param("id")
	rsrcUser, err := s.userStore.RetrieveById(userId)
	if err != nil {
		return false, err
	}
	curUser, err := s.fetchUserFromSession(c)
	if err != nil {
		return false, err
	}

	ok := rsrcUser.Pk == curUser.Pk
	return ok, nil
}

func (s *Server) postSignUp() gin.HandlerFunc {
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
				Profile:  "", // 後から設定させる
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

			c.JSON(http.StatusOK, gin.H{
				"id":   signupUser.Id,
				"name": signupUser.Name,
			})
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

func (s *Server) postSignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signinUserReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			log.Printf("%s\n", err.Error())
			return
		}

		user, err := s.userStore.RetrieveById(req.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				err := errors.New("id or password is wrong.")
				c.JSON(http.StatusBadRequest, errorRes(err))
				log.Printf("%s\n", err.Error())
				return
			} else {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				log.Printf("%s\n", err.Error())
				return
			}
		}

		if !verifyPasswordHash(user.Password, req.Password) {
			err := errors.New("id or password is wrong.")
			c.JSON(http.StatusBadRequest, errorRes(err))
			log.Printf("%s\n", err.Error())
			return
		}

		token, err := createUuid()
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			log.Printf("%s\n", err.Error())
			return
		}

		if err = s.sessionStore.Create(token, user); err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			log.Printf("%s\n", err.Error())
			return
		}
		session := sessions.Default(c)
		session.Set(SESSION_TOKEN, token)
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"id":   user.Id,
			"name": user.Name,
		})
	}
}

func (s *Server) postSignOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := s.authorizationCheck(c)
		if err != nil || !ok {
			c.JSON(http.StatusForbidden, errorRes(errors.New("forbidden")))
			return
		}
		session := sessions.Default(c)
		token, err := s.fetchSessToken(session)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		if err := s.sessionStore.DeleteByToken(token); err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		// セッション破棄
		session.Clear()
		// クッキー削除
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	}
}

func (s *Server) postUnsubscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := s.authorizationCheck(c)
		if err != nil || !ok {
			c.JSON(http.StatusForbidden, errorRes(errors.New("forbidden")))
			return
		}
		session := sessions.Default(c)
		token, err := s.fetchSessToken(session)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		sess, err := s.sessionStore.RetrieveByToken(token)
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		err = s.userStore.DeleteByPk(sess.UserPk)
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		// セッション破棄
		session.Clear()
		// クッキー削除
		session.Options(sessions.Options{MaxAge: -1})
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	}
}
