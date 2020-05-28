// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// 🙏 Thanks & Credits to @arsmn

package embed

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/gofiber/fiber"
)

// Config holds the configuration for the middleware
type Config struct {
	// Root is a FileSystem that provides access
	// to a collection of files and directories.
	// Required. Default: nil
	Root http.FileSystem

	// ErrorHandler defines the response body when an error raised.
	// Optional. Defaul: Next for NotExistError and 403 for PermissionError and 500 for others
	ErrorHandler func(*fiber.Ctx, error)

	// Index file for serving a directory.
	// Optional. Default: "index.html"
	Index string

	// Enable directory browsing.
	// Optional. Default: false
	Browse bool
}

// New returns an embed middleware for serving files
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

		// Serve index if path is directory
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

		// Browse directory if no index found and browsing is enabled
		if stat.IsDir() {
			if cfg.Browse {
				if err := dirList(c, file); err != nil {
					cfg.ErrorHandler(c, err)
				}
				return
			}
			c.SendStatus(fiber.StatusForbidden)
			return
		}

		modTime := stat.ModTime()
		contentLength := int(stat.Size())

		// Set Content Type header
		c.Type(getFileExtension(stat.Name()))

		// Set Last Modified header
		if !modTime.IsZero() {
			c.Set(fiber.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
		}

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

func dirList(c *fiber.Ctx, f http.File) error {
	fileinfos, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	fm := make(map[string]os.FileInfo, len(fileinfos))
	filenames := make([]string, 0, len(fileinfos))
	for _, fi := range fileinfos {
		name := fi.Name()
		fm[name] = fi
		filenames = append(filenames, name)
	}

	basePathEscaped := html.EscapeString(c.Path())
	c.Write(fmt.Sprintf("<html><head><title>%s</title><style>.dir { font-weight: bold }</style></head><body>", basePathEscaped))
	c.Write(fmt.Sprintf("<h1>%s</h1>", basePathEscaped))
	c.Write(fmt.Sprintf("<ul>"))

	if len(basePathEscaped) > 1 {
		parentPathEscaped := html.EscapeString(c.Path() + "/..")
		c.Write(fmt.Sprintf(`<li><a href="%s" class="dir">..</a></li>`, parentPathEscaped))
	}

	sort.Strings(filenames)
	for _, name := range filenames {
		pathEscaped := html.EscapeString(path.Join(c.Path() + "/" + name))
		fi := fm[name]
		auxStr := "dir"
		className := "dir"
		if !fi.IsDir() {
			auxStr = fmt.Sprintf("file, %d bytes", fi.Size())
			className = "file"
		}
		c.Write(fmt.Sprintf(`<li><a href="%s" class="%s">%s</a>, %s, last modified %s</li>`,
			pathEscaped, className, html.EscapeString(name), auxStr, fi.ModTime()))
	}
	c.Write(fmt.Sprintf("</ul></body></html>"))

	c.Type("html")

	return nil
}
