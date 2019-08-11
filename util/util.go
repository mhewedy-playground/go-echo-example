package util

import (
	"github.com/labstack/echo"
	"io"
)

func SkipperFn(skippedURLs []string) func(echo.Context) bool {
	return func(context echo.Context) bool {
		for _, url := range skippedURLs {
			if url == context.Request().URL.String() {
				return true
			}
		}
		return false
	}
}

type teeWriter struct {
	writers []io.Writer
}

func NewTeeWriter(writers []io.Writer) *teeWriter {
	return &teeWriter{writers: writers}
}

func (t teeWriter) Write(p []byte) (n int, err error) {
	for _, writer := range t.writers {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}
	return
}
