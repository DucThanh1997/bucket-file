package controller

import (
	container_camera "bucket_file/cmd/api/container"
	"bucket_file/server"

	"github.com/gin-gonic/gin"
)

var (
	container *container_camera.Container
)

func NewRouter(c *container_camera.Container) error {
	container = c
	router := server.New()
	router.Use(gin.Logger())

	router.Cors(server.CorsConfig{})

	authen, err := router.Authen(server.AuthenConfig{
		SecretKey:     c.Config.TokenSecretKey,
		ExpiredHour:   c.Config.TokenExpiredHour,
		Authenticator: login,
		Verification:  userVerify,
	})
	if err != nil {
		return err
	}
	v1 := router.Group("/api/v1")
	v1.POST("/users", CreateUser)
	v1.POST("/login", authen.LoginHandler)
	// v1.Use(authen.TokenAuthMiddleware())

	user := v1.Group("/users")
	{
		user.PUT("", UpdateUser)
		user.GET("", GetListUsers)
		user.PATCH("", ChangePassword)
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/refresh", RefreshToken)
	}

	bucket := v1.Group("/bucket")
	{
		bucket.PUT("", CreateBucket)
		bucket.GET("", GetBucket)
		// bucket.PUT("/file", Create_file)
		// bucket.GET("/file", GetFile)
	}

	return router.Run(c.Config.Listen)
}
