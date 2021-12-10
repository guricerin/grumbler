package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/util"
)

const SESSION_TOKEN string = "sesstoken"

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
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysessions", store))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"home": "hohoho",
		})
	})
	router.GET("/login", func(c *gin.Context) {
		// todo
		c.JSON(http.StatusOK, gin.H{
			"hoge": false,
		})
	})
	router.POST("/login", s.postLogin())
	router.POST("/signup", s.postSignup())

	menu := router.Group("/menu")
	menu.Use(s.authenticationMiddleware())
	{
		menu.GET("/user/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"is_logged_in": true,
			})
		})
	}

	s.router = router
}

// 認証
func (s *Server) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get(SESSION_TOKEN)
		if token == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			// セッション固定化攻撃対策 : 認証毎に新たなトークンを発行
			newToken, err := createUuid()
			if err != nil {
				// todo: err msgをユーザ用に変更
				c.JSON(http.StatusInternalServerError, errorRes(err))
				c.Abort()
				return
			}

			oldToken, ok := token.(string)
			if !ok {
				// todo: err msgをユーザ用に変更
				err = errors.New("token is not string")
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
}

func (s *Server) fetchUserFromSession(c *gin.Context) error {
	token := sessions.Default(c).Get(SESSION_TOKEN)
	if token == nil {
		// todo
	}
	// todo: トークンをキーにログインユーザ情報を取得
	return nil
}

func errorRes(err error) gin.H {
	return gin.H{"error": err.Error()}
}
