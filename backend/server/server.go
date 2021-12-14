package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
	"github.com/guricerin/grumbler/backend/util"
)

const (
	SESSION_TOKEN string = "sesstoken"
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
		MaxAge:   60 * 60 * 24 * 7, // 寿命は一週間
		HttpOnly: true,             // JSなどからのクッキーへのアクセスを禁止
	})
	router.Use(sessions.Sessions("grumbler_session", store))
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost"},
	// 	AllowMethods:     []string{"GET", "POST", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	// 	AllowCredentials: true,
	// 	// PreFlight要求がキャッシュされる時間
	// 	MaxAge: 24 * time.Hour,
	// }))
	corsconf := cors.DefaultConfig()
	corsconf.AllowOrigins = []string{"http://localhost"}
	router.Use(cors.New(corsconf))

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"home": "hohoho",
		})
	})
	router.POST("/api/login", s.postLogin())
	router.POST("/api/signup", s.postSignup())

	auth := router.Group("/api/auth")
	auth.Use(s.authenticationMiddleware())
	{
		auth.GET("/user/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"is_logged_in": true,
			})
		})
		auth.GET("/user/:id/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"is_logged_in": false,
			})
		})
	}

	s.router = router
}

// 認証
func (s *Server) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get(SESSION_TOKEN)
		if v == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		oldToken, ok := v.(string)
		if !ok {
			// todo: err msgをユーザ用に変更
			err := errors.New("token is not string")
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

		s.sessionStore.Update(oldToken, newToken)
		session.Set(SESSION_TOKEN, newToken)
		session.Save()
		c.Next()
	}
}

func (s *Server) fetchUserFromSession(c *gin.Context) (user model.User, err error) {
	user = model.User{}
	session := sessions.Default(c)
	v := session.Get(SESSION_TOKEN)
	if v == nil {
		err = errors.New("cookie value not set.")
		return
	}
	token, ok := v.(string)
	if !ok {
		err = errors.New("cookie value is not string.")
		return
	}
	sess, err := s.sessionStore.RetrieveByToken(token)
	if err != nil {
		return
	}
	user, err = s.userStore.RetrieveByPk(sess.UserPk)
	return
}

func errorRes(err error) gin.H {
	es := []string{err.Error()}
	return gin.H{"error": es}
}
