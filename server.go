package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mhewedy/echo-example/controllers"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func SkipperFn(skipURLs []string) func(echo.Context) bool {
	return func(context echo.Context) bool {
		for _, url := range skipURLs {
			if url == context.Request().URL.String() {
				return true
			}
		}
		return false
	}
}

type TeeWriter struct {
	w []io.Writer
}

func (t TeeWriter) Write(p []byte) (n int, err error) {
	for _, writer := range t.w {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}
	return
}

func main() {
	e := echo.New()

	// Logger
	e.Logger.SetOutput(TeeWriter{[]io.Writer{
		os.Stdout, &lumberjack.Logger{
			Filename:   "/tmp/myapp.log",
			MaxSize:    500,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
	}})

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(
		middleware.JWTConfig{
			SigningKey: []byte("secret"),
			Skipper:    SkipperFn([]string{"/api/v1/login", "/"}),
		}))

	// Routes
	e.POST("/api/v1/login", controllers.Login)
	e.GET("/", controllers.Home)
	e.GET("/api/v1/me", controllers.Me)

	e.Logger.Fatal(e.Start(":8000"))
}
