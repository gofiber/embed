// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// 🙏 Thanks & Credits to @arsmn

package embed

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
)

func Test_Embed(t *testing.T) {
	app := *fiber.New()

	app.Use("/test", New(Config{
		Root: http.Dir("./testdata"),
	}))

	app.Use("/dir", New(Config{
		Root:   http.Dir("./testdata"),
		Browse: true,
	}))

	app.Get("/", func(c *fiber.Ctx) {
		c.SendString("Hello, World!")
	})

	tests := []struct {
		name         string
		url          string
		statusCode   int
		contentType  string
		modifiedTime string
	}{
		{
			name:         "Should be returns status 200 with suitable content-type",
			url:          "/test/index.html",
			statusCode:   200,
			contentType:  "text/html",
			modifiedTime: "Thu, 07 May 2020 15:40:26 GMT",
		},
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:         "Should be returns status 200 with suitable content-type",
			url:          "/test/test.json",
			statusCode:   200,
			contentType:  "application/json",
			modifiedTime: "Tue, 21 Apr 2020 11:58:15 GMT",
		},
		{
			name:         "Should be returns status 200 with suitable content-type",
			url:          "/test/main.css",
			statusCode:   200,
			contentType:  "text/css",
			modifiedTime: "Tue, 21 Apr 2020 11:34:42 GMT",
		},
		{
			name:       "Should be returns status 404",
			url:        "/test/nofile.js",
			statusCode: 404,
		},
		{
			name:       "Should be returns status 404",
			url:        "/test/nofile",
			statusCode: 404,
		},
		{
			name:        "Should be returns status 200",
			url:         "/",
			statusCode:  200,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:       "Should be returns status 403",
			url:        "/test/inner",
			statusCode: 403,
		},
		{
			name:        "Should list the directory contents",
			url:         "/dir/inner",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:         "Should be returns status 200",
			url:          "/dir/inner/fiber.png",
			statusCode:   200,
			contentType:  "image/png",
			modifiedTime: "Thu, 07 May 2020 14:44:24 GMT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf(`%s: %s`, t.Name(), err)
			}

			if resp.StatusCode != tt.statusCode {
				t.Fatalf(`%s: StatusCode: got %v - expected %v`, t.Name(), resp.StatusCode, tt.statusCode)
			}

			if tt.contentType != "" {
				ct := resp.Header.Get("Content-Type")
				if ct != tt.contentType {
					t.Fatalf(`%s: Content-Type: got %s - expected %s`, t.Name(), ct, tt.contentType)
				}
			}

			if tt.modifiedTime != "" {
				lm := resp.Header.Get("Last-Modified")
				if lm != tt.modifiedTime {
					t.Fatalf(`%s: Last-Modified: got %s - expected %s`, t.Name(), lm, tt.modifiedTime)
				}
			}
		})
	}
}