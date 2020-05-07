// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// 🙏 Thanks & Credits to @arsmn

package embed

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/gofiber/fiber"
)

type Config struct {
	Root         http.FileSystem
	ErrorHandler func(*fiber.Ctx, error)
	Index        string
	DirList      bool
}

// New returns
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
			if os.IsNotExist(err) {
				c.Next()
				return
			}
			if os.IsPermission(err) {
				c.SendStatus(fiber.StatusForbidden)
				return
			}
			c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	if cfg.Index == "" {
		cfg.Index = "/index.html"
	}

	if !strings.HasPrefix(cfg.Index, "/") {
		cfg.Index = "/" + cfg.Index
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

		file, err := cfg.Root.Open(path)
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}

		if stat.IsDir() {
			indexPath := strings.TrimSuffix(path, "/") + cfg.Index
			index, err := cfg.Root.Open(indexPath)
			if err == nil {
				indexStat, err := index.Stat()
				if err == nil {
					file = index
					stat = indexStat
				}
			}
		}

		if stat.IsDir() {
			if cfg.DirList {
				if err := dirList(c, file); err != nil {
					cfg.ErrorHandler(c, err)
				}
				return
			}
			c.SendStatus(fiber.StatusForbidden)
			return
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

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

func dirList(c *fiber.Ctx, f http.File) error {
	dirs, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

	c.Type("html")
	fmt.Fprintf(c.Fasthttp, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: path.Join(c.Path(), "/", name)}
		fmt.Fprintf(c.Fasthttp, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
	}
	fmt.Fprintf(c.Fasthttp, "</pre>\n")

	return nil
}
