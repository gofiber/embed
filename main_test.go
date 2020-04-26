// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// 🙏 Thanks & Credits to @arsmn

package embed

import (
	"fmt"
	"net/http"
	"testing"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gobuffalo/packr/v2"
	"github.com/gofiber/fiber"
	"github.com/markbates/pkger"
)

func Test_Embed(t *testing.T) {
	app := *fiber.New()

	// pkger
	app.Use("/pkger", New(Config{
		Root: pkger.Dir("/testdata"),
	}))

	// packr
	app.Use("/packr", New(Config{
		Root: packr.New("box", "./testdata"),
	}))

	// go.rice
	app.Use("/rice", New(Config{
		Root: rice.MustFindBox("./testdata").HTTPBox(),
	}))

	app.Get("/", func(c *fiber.Ctx) {
		c.SendString("Hello, World!")
	})

	embedders := []string{"pkger", "packr", "rice"}

	tests := []struct {
		name        string
		fileName    string
		statusCode  int
		contentType string
	}{
		// {
		// 	name:        "Should be returns status 200 with suitable content-type",
		// 	fileName:    "index.html",
		// 	statusCode:  200,
		// 	contentType: "text/html",
		// },
		// {
		// 	name:        "Should be returns status 200 with suitable content-type",
		// 	fileName:    "test.json",
		// 	statusCode:  200,
		// 	contentType: "application/json",
		// },
		// {
		// 	name:        "Should be returns status 200 with suitable content-type",
		// 	fileName:    "main.css",
		// 	statusCode:  200,
		// 	contentType: "text/css",
		// },
		// {
		// 	name:       "Should be returns status 404",
		// 	fileName:   "nofile.js",
		// 	statusCode: 404,
		// },
		{
			name:       "Should be returns status 404",
			fileName:   "nofile/",
			statusCode: 404,
		},
	}

	for _, embedder := range embedders {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				url := fmt.Sprintf("/%s/%s", embedder, tt.fileName)
				req, _ := http.NewRequest("GET", url, nil)
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
			})
		}
	}
}
