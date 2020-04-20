### Embed
FileServer middleware for Fiber

Special thanks and credits to [Alireza Salary](https://github.com/arsmn)

### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/embed
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

### pkger example

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

### packr example

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

### go.rice example

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
