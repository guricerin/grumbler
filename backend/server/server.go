package server

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guricerin/grumbler/backend/model"
	"github.com/guricerin/grumbler/backend/util"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	SESSION_TOKEN string = "grumbler_sesstoken"
)

type Server struct {
	cfg          util.Config
	logger       zerolog.Logger
	router       *gin.Engine
	userStore    userStore
	sessionStore sessionStore
	grumbleStore grumbleStore
	followStore  followStore
}

func NewServer(cfg util.Config, db *sql.DB) Server {
	s := Server{
		cfg:          cfg,
		userStore:    NewUserStore(db),
		sessionStore: NewSessionStore(db),
		grumbleStore: NewGrumbleStore(db),
		followStore:  NewFollowStore(db),
	}
	s.setupRouter()
	s.setupLogger(false)
	return s
}

func (s *Server) Run() error {
	url := fmt.Sprintf("%s:%s", s.cfg.ServerHost, s.cfg.ServerPort)
	return s.router.Run(url)
}

func (s *Server) setupRouter() {
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{s.cfg.FrontOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "OPTION"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		// PreFlight要求がキャッシュされる時間
		MaxAge: 24 * time.Hour,
	}))

	router.Use(s.limitReqBodySizeMiddleware())
	router.Use(s.RequestBodyLog)

	router.GET("/api/signin-check", s.signinCheck())
	router.POST("/api/signin", s.postSignIn())
	router.POST("/api/signup", s.postSignUp())

	auth := router.Group("/api/auth")
	auth.Use(s.authenticationMiddleware())
	{
		auth.POST("/search", s.getSearch())
		auth.POST("/signout", s.postSignOut())
		auth.POST("/unsubscribe", s.postUnsubscribe())
		auth.POST("/grumble", s.postGrumble())
		auth.POST("/delete-grumble", s.postDeleteGrumble())
		auth.POST("/follow", s.postFollow())
		auth.POST("/unfollow", s.postUnFollow())
		auth.POST("/settings", s.postUserSettings())
		auth.POST("/bookmark", s.postBookmark())
		auth.POST("/delete-bookmark", s.postDeleteBookmark())
		auth.POST("/reply", s.postReply())
		auth.POST("/regrumble", s.postRegrumble())
		auth.POST("/delete-regrumble", s.postDeleteRegrumble())
		auth.GET("/grumble/:grumble_pk", s.getGrumbleDetail())

		auth.GET("/user/:id", s.getUser())
		auth.GET("/user/:id/detail", s.getUserDetail())
		auth.GET("/user/:id/timeline", s.getTimeline())
		auth.POST("/user/:id/settings", s.postUserSettings())
	}

	s.router = router
}

func (s *Server) setupLogger(isDebug bool) {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	// ログファイルのローテーション
	rotator := &lumberjack.Logger{
		Filename:   s.cfg.LogFilePath,
		MaxSize:    s.cfg.LogFileMaxSize, // mbyte
		MaxBackups: s.cfg.LogFileMaxBackups,
		MaxAge:     s.cfg.LogFileMaxAge, // 古いログファイルの寿命（day）
		Compress:   true,                // 古いログファイルをgzipで圧縮
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(io.MultiWriter(os.Stderr, rotator)).With().Stack().Timestamp().Logger()
	s.logger = logger
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

func (s *Server) limitReqBodySizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(s.cfg.RequestContentLengthMaxByte))
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
	c.SetCookie(SESSION_TOKEN, token, week, "/api", s.cfg.ServerDomain, true, true)
}

func (s *Server) deleteCookie(c *gin.Context) (err error) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SESSION_TOKEN, "dummy", -1, "/api", s.cfg.ServerDomain, true, true)
	return
}

func (s *Server) fetchUserFromSession(c *gin.Context) (user model.User, err error) {
	user = model.User{}
	token, err := s.fetchSessToken(c)
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
