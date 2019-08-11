package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mhewedy/echo-example/controllers"
	"github.com/mhewedy/echo-example/util"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func main() {
	e := echo.New()

	// Logger
	e.Logger.SetOutput(util.NewTeeWriter([]io.Writer{
		os.Stdout, &lumberjack.Logger{
			Filename:   "/tmp/myapp.log",
			MaxSize:    500,
			MaxAge:     28,
			MaxBackups: 3,
			Compress:   true,
		},
	}))

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			SigningKey: []byte("secret"),
			Skipper:    util.SkipperFn([]string{"/api/v1/login", "/"}),
		}))

	// Routes
	e.POST("/api/v1/login", controllers.Login)
	e.GET("/", controllers.Home)
	e.GET("/api/v1/me", controllers.Me)

	e.Logger.Fatal(e.Start(":8000"))
}
