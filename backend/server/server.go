package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
	"github.com/guricerin/grumbler/backend/util"
)

const (
	SESSION_TOKEN string = "grumbler_sesstoken"
)

type Server struct {
	cfg          util.Config
	router       *gin.Engine
	userStore    userStore
	sessionStore sessionStore
	grumbleStore grumbleStore
}

func NewServer(cfg util.Config, db *sql.DB) Server {
	s := Server{
		cfg:          cfg,
		userStore:    NewUserStore(db),
		sessionStore: NewSessionStore(db),
		grumbleStore: NewGrumbleStore(db),
	}
	s.setupRouter()
	return s
}

func (s *Server) Run() error {
	url := fmt.Sprintf("%s:%s", s.cfg.ServerHost, s.cfg.ServerPort)
	return s.router.Run(url)
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "OPTION"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		// PreFlight要求がキャッシュされる時間
		MaxAge: 24 * time.Hour,
	}))

	router.GET("/api/signin-check", s.signinCheck())
	router.POST("/api/signin", s.postSignIn())
	router.POST("/api/signup", s.postSignUp())

	auth := router.Group("/api/auth")
	auth.Use(s.authenticationMiddleware())
	{
		auth.GET("/user/:id", s.getUser())
		auth.GET("/search", s.getSearch())
		auth.POST("/user/:id/signout", s.postSignOut())
		auth.POST("/user/:id/unsubscribe", s.postUnsubscribe())
		auth.POST("/grumble", s.postGrumble())
		auth.GET("/grumbles", s.getGrumbles())
	}

	s.router = router
}

// 認証
func (s *Server) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := s.fetchSessToken(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			c.Abort()
			return
		}
		_, err = s.sessionStore.RetrieveByToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			c.Abort()
			return
		}

		c.Next()
	}
}

// セッション固定化攻撃対策 : 認証毎に新たなトークンを発行
func (s *Server) resetSessToken(c *gin.Context) error {
	oldToken, err := s.fetchSessToken(c)
	if err != nil {
		return err
	}

	newToken, err := createUuid()
	if err != nil {
		return err
	}

	err = s.sessionStore.Update(oldToken, newToken)
	if err != nil {
		return err
	}

	s.setCookie(c, newToken)
	return nil
}

func (s *Server) fetchSessToken(c *gin.Context) (string, error) {
	token, err := c.Cookie(SESSION_TOKEN)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Server) setCookie(c *gin.Context, token string) {
	week := 60 * 60 * 24 * 7
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SESSION_TOKEN, token, week, "/api", "localhost", false, true)
}

func (s *Server) deleteCookie(c *gin.Context) (err error) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SESSION_TOKEN, "dummy", -1, "/api", "localhost", false, true)
	return
}

func (s *Server) fetchUserFromSession(c *gin.Context) (user model.User, err error) {
	user = model.User{}
	token, err := s.fetchSessToken(c)
	if err != nil {
		return
	}
	log.Printf("fetchUserFromSession() token: %s\n", token)
	sess, err := s.sessionStore.RetrieveByToken(token)
	if err != nil {
		log.Println("failed to RetrieveByToken")
		return
	}
	user, err = s.userStore.RetrieveByPk(sess.UserPk)
	return
}

// ページリソースのユーザと、それにアクセスしようとしているユーザは同一か
func (s *Server) authorizationCheck(c *gin.Context) (model.User, error) {
	userId := c.Param("id")
	log.Printf("id: %s\n", userId)
	dummy := model.User{}
	rsrcUser, err := s.userStore.RetrieveById(userId)
	if err != nil {
		return dummy, err
	}
	curUser, err := s.fetchUserFromSession(c)
	if err != nil {
		return dummy, err
	}

	if rsrcUser.Pk == curUser.Pk {
		return rsrcUser, nil
	} else {
		return dummy, errors.New("wrong user")
	}
}

func userRes(user model.User) gin.H {
	return gin.H{
		"id":      user.Id,
		"name":    user.Name,
		"profile": user.Profile,
	}
}

func errorRes(err error) gin.H {
	es := []string{err.Error()}
	return gin.H{"error": es}
}
