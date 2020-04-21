// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// 🙏 Thanks & Credits to @arsmn

package embed

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber"
)

type Config struct {
	Prefix       string
	Root         http.FileSystem
	ErrorHandler func(*fiber.Ctx, error)
}

func New(config ...Config) func(*fiber.Ctx) {
	var cfg Config

	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.Root == nil {
		log.Fatal("Fiber: FileServer middleware requires root")
	}

	if cfg.Prefix == "" {
		cfg.Prefix = "/"
	}

	if !strings.HasPrefix(cfg.Prefix, "/") {
		cfg.Prefix = "/" + cfg.Prefix
	}

	if !strings.HasSuffix(cfg.Prefix, "/") {
		cfg.Prefix = cfg.Prefix + "/"
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) {
			c.Status(fiber.StatusNotFound)
			c.SendString("File not found")
		}
	}

	return func(c *fiber.Ctx) {
		p := c.Path()

		if !strings.HasPrefix(p, cfg.Prefix) {
			c.Next()
			return
		}
		p = strings.TrimPrefix(p, cfg.Prefix)

		file, err := cfg.Root.Open(filepath.Clean(p))
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}
		contentLength := int(stat.Size())

		contentType, err := detectContentType(file, stat)
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}

		c.Fasthttp.SetContentType(contentType)
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

func detectContentType(f http.File, fileInfo os.FileInfo) (string, error) {
	ext := getFileExtension(fileInfo.Name())
	contentType := mime.TypeByExtension(ext)
	if len(contentType) == 0 {
		data, err := readFileHeader(f)
		if err != nil {
			return "", err
		}
		contentType = http.DetectContentType(data)
	}

	return contentType, nil
}

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}

	return path[n:]
}

func readFileHeader(f http.File) ([]byte, error) {
	r := io.Reader(f)
	lr := &io.LimitedReader{
		R: r,
		N: 512,
	}
	data, err := ioutil.ReadAll(lr)
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	return data, err
}
