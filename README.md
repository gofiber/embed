### Embed
FileServer middleware for Fiber

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/session
```

### Signature
```go
embed.New(config ...embed.Config) func(c *fiber.Ctx)
```

### Config
| Property | Type | Description | Default |
| :--- | :--- | :--- | :--- |
| Prefix | `string` | Path prefix | `/` |
| Root | `http.FileSystem` | http.FileSystem to use | `nil` |
| ErrorHandler | `func(*fiber.Ctx, error)` | Error handler | `404 File not found` |

### Examples
### pkger

```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/gofiber/embed"

	"github.com/markbates/pkger"
)

func main() {
	app := fiber.New()
	dir := pkger.Dir("/assets")

	app.Use(embed.New(embed.Config{
		Prefix: "/assets",
		Root:  dir,
	}))

	app.Listen(8080)
}
```

### packr

```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/gofiber/embed"

	"github.com/gobuffalo/packr/v2"
)

func main() {
	app := fiber.New()
	assetsBox := packr.New("Assets Box", "/assets")

	app.Use(embed.New(embed.Config{
		Prefix: "/assets",
		Root:   assetsBox,
	}))

	app.Listen(8080)
}
```

### go.rice

```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/gofiber/embed"

	rice "github.com/GeertJohan/go.rice"
)

func main() {
	app := fiber.New()
	assetsBox := rice.MustFindBox("assets")

	app.Use(embed.New(embed.Config{
		Prefix: "/assets",
		Root:   assetsBox.HTTPBox(),
	}))

	app.Listen(8080)
}
```