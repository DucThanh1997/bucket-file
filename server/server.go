package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	_ "github.com/streadway/amqp"
	"github.com/gin-gonic/gin"
	jwt "github.com/namcuongq/gin-jwt"
	"github.com/natefinch/lumberjack"
)

type server struct {
	*gin.Engine
}

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
}

type CorsConfig struct {
	AllowOrigin  string
	MaxAge       string
	AllowMethods string
	AllowHeaders string
}

type AuthenConfig struct {
	SecretKey     string
	ExpiredHour   uint64
	Authenticator func(c *gin.Context) (map[string]interface{}, error)
	Verification  func(map[string]interface{}) (bool, error)
}

type query struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	SortBy string
}

func GetUserId(c *gin.Context) string {
	return fmt.Sprintf("%v", c.Value("id"))
}

func PrepQuery(c *gin.Context) query {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	page = page - 1
	if page < 0 {
		page = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return query{page, limit, search, sort, sortBy}
}

func New() *server {
	router := &server{gin.New()}
	router.Use(gin.Recovery())

	return router
}

func (r *server) Cors(config CorsConfig) {
	r.Use(corsMiddleware(config))
}

func (r *server) Authen(authen AuthenConfig) (*jwt.JWT, error) {
	return jwt.New(&jwt.JWT{
		SecretKey:     authen.SecretKey,
		ExpiredHour:   authen.ExpiredHour,
		Authenticator: authen.Authenticator,
		Verification:  authen.Verification,
	})
}

func Data(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Code: http.StatusOK,
		Data: data,
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, response{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	})
}

func Success(c *gin.Context, mess ...string) {
	var message = "Success"
	if mess != nil {
		message = mess[0]
	}

	c.JSON(http.StatusOK, response{
		Code:    http.StatusOK,
		Message: message,
	})
}

func BadRequest(c *gin.Context, mess ...string) {
	var message = "Bad Request"
	if mess != nil {
		message = mess[0]
	}

	c.JSON(http.StatusBadRequest, response{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func NotFound(c *gin.Context, mess ...string) {
	var message = "Not Found"
	if mess != nil {
		message = mess[0]
	}

	c.JSON(http.StatusNotFound, response{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func InternalServerError(c *gin.Context, mess ...string) {
	var message = "Internal Server Error"
	if mess != nil {
		message = mess[0]
	}

	c.JSON(http.StatusInternalServerError, response{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

func corsMiddleware(config CorsConfig) gin.HandlerFunc {
	if len(config.AllowOrigin) < 1 {
		config.AllowOrigin = "*"
	}

	if len(config.MaxAge) < 1 {
		config.MaxAge = "86400"
	}

	if len(config.AllowHeaders) < 1 {
		config.AllowHeaders = "POST, GET, PUT, DELETE"
	}

	if len(config.AllowMethods) < 1 {
		config.AllowMethods = "Authorization, Content-Type"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
		c.Writer.Header().Set("Access-Control-Max-Age", config.MaxAge)
		c.Writer.Header().Set("Access-Control-Allow-Methods", config.AllowHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Headers", config.AllowMethods)
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

type accessLog struct {
	*log.Logger
}

type LoggerConfig struct {
	Filename string
	MaxSize  int
	MaxAge   int
	Compress bool
}

func NewLogger(config LoggerConfig) *accessLog {
	var logger = accessLog{log.New(nil, "", log.LstdFlags)}

	if config.Filename == "" {
		config.Filename = "access.log"
	}

	if config.MaxSize < 1 {
		config.MaxSize = 5
	}

	if config.MaxAge < 1 {
		config.MaxAge = 5
	}

	logger.SetOutput(&lumberjack.Logger{
		Filename: config.Filename,
		MaxSize:  config.MaxSize,
		MaxAge:   config.MaxAge,
		Compress: config.Compress,
	})

	return &logger
}

func (logger *accessLog) SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		userAgent := c.GetHeader("User-Agent")

		logger.Printf("| %v | %v | %v | %v %v | %v\n",
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
			userAgent)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}

