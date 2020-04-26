// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// 🙏 Thanks & Credits to @arsmn

package embed

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber"
)

type Config struct {
	Root         http.FileSystem
	ErrorHandler func(*fiber.Ctx, error)
	Index        string
}

func New(config ...Config) func(*fiber.Ctx) {
	var cfg Config

	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.Root == nil {
		log.Fatal("Fiber: Embed middleware requires root")
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) {
			c.Status(fiber.StatusNotFound)
			c.SendString("File not found")
		}
	}

	if cfg.Index == "" {
		cfg.Index = "index.html"
	}

	var prefix string
	return func(c *fiber.Ctx) {

		// Set prefix
		if len(prefix) == 0 {
			prefix = c.Route().Path
		}

		// Strip prefix
		path := strings.TrimPrefix(c.Path(), prefix)
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		file, err := cfg.Root.Open(filepath.Clean(path))
		if err != nil {
			if err.Error() == "file does not exist" {
				c.Next()
				return
			}
			cfg.ErrorHandler(c, err)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}

		if stat.IsDir() {
			index, err := cfg.Root.Open(path + "/" + cfg.Index)
			if err != nil {
				cfg.ErrorHandler(c, err)
				return
			}
			indexStat, err := index.Stat()
			if err != nil {
				cfg.ErrorHandler(c, err)
				return
			}

			file = index
			stat = indexStat
		}

		contentLength := int(stat.Size())

		// Set Content Type header
		c.Type(getFileExtension(stat.Name()))

		if c.Method() == fiber.MethodGet {
			c.Fasthttp.SetBodyStream(file, contentLength)
			return
		} else if c.Method() == fiber.MethodHead {
			c.Fasthttp.ResetBody()
			c.Fasthttp.Response.SkipBody = true
			c.Fasthttp.Response.Header.SetContentLength(contentLength)
			if err := file.Close(); err != nil {
				cfg.ErrorHandler(c, err)
				return
			}
			return
		}

		c.Next()
	}
}

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}
	return path[n:]
}
