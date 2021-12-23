package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
}

func NewServer(cfg util.Config, db *sql.DB) Server {
	s := Server{
		cfg:          cfg,
		userStore:    NewUserStore(db),
		sessionStore: NewSessionStore(db),
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
	// todo: 'secret'は設定ファイルで指定可能にすべき？
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/api",
		MaxAge:   60 * 60 * 24 * 7, // 寿命は一週間
		HttpOnly: true,             // JSなどからのクッキーへのアクセスを禁止
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions(SESSION_TOKEN, store))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "OPTION"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		// PreFlight要求がキャッシュされる時間
		MaxAge: 24 * time.Hour,
	}))

	router.POST("/api/signin", s.postSignIn())
	router.POST("/api/signup", s.postSignUp())
	router.GET("/api/signin-check", s.signinCheck())

	auth := router.Group("/api/auth")
	auth.Use(s.authenticationMiddleware())
	{
		auth.GET("/user/:id", s.getUser())
		auth.GET("/search", s.getSearch())
		auth.POST("/user/:id/signout", s.postSignOut())
		auth.POST("/user/:id/unsubscribe", s.postUnsubscribe())
	}

	s.router = router
}

// 認証
func (s *Server) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		oldToken, err := s.fetchSessToken(session)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorRes(err))
			c.Abort()
			return
		}

		// セッション固定化攻撃対策 : 認証毎に新たなトークンを発行
		newToken, err := createUuid()
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			c.Abort()
			return
		}

		err = s.sessionStore.Update(oldToken, newToken)
		if err != nil {
			// todo: err msgをユーザ用に変更
			c.JSON(http.StatusInternalServerError, errorRes(err))
			c.Abort()
			return
		}
		session.Clear()
		session.Set(SESSION_TOKEN, newToken)
		session.Save()
		c.Next()
	}
}

func (s *Server) fetchSessToken(session sessions.Session) (string, error) {
	v := session.Get(SESSION_TOKEN)
	if v == nil {
		err := errors.New("unauthenticated")
		return "", err
	}
	token, ok := v.(string)
	if !ok {
		err := errors.New("illegal token")
		return "", err
	}
	return token, nil
}

func (s *Server) deleteCookie(c *gin.Context) (err error) {
	c.SetCookie(SESSION_TOKEN, "dummy", -1, "/api", "localhost", false, true)
	return
}

func (s *Server) fetchUserFromSession(c *gin.Context) (user model.User, err error) {
	user = model.User{}
	session := sessions.Default(c)
	token, err := s.fetchSessToken(session)
	if err != nil {
		return
	}
	sess, err := s.sessionStore.RetrieveByToken(token)
	if err != nil {
		return
	}
	user, err = s.userStore.RetrieveByPk(sess.UserPk)
	return
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
