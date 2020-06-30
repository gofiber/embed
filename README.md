# Embed

![Release](https://img.shields.io/github/release/gofiber/embed.svg)
[![Discord](https://img.shields.io/badge/discord-join%20channel-7289DA)](https://gofiber.io/discord)
![Test](https://github.com/gofiber/embed/workflows/Test/badge.svg)
![Security](https://github.com/gofiber/embed/workflows/Security/badge.svg)
![Linter](https://github.com/gofiber/embed/workflows/Linter/badge.svg)

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
| Index | `string` | Index file name | `index.html` |
| Browse | `bool` | Enable directory browsing | `false` |
| Root | `http.FileSystem` | http.FileSystem to use | `nil` |
| ErrorHandler | `func(*fiber.Ctx, error)` | Error handler | `InternalServerError` |

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

  app.Use("/assets", embed.New(embed.Config{
    Root:   pkger.Dir("/assets"),
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

  app.Use("/assets", embed.New(embed.Config{
    Root:   packr.New("Assets Box", "/assets"),
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

  "github.com/GeertJohan/go.rice"
)

func main() {
  app := fiber.New()

  app.Use("/assets", embed.New(embed.Config{
    Prefix: "/assets/",
    Root:   rice.MustFindBox("assets").HTTPBox(),
  }))

  app.Listen(8080)
}
```

### fileb0x

```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/gofiber/embed"
  "<Your go module>/myEmbeddedFiles"
)

func main() {
  app := fiber.New()

  app.Use("/assets", embed.New(embed.Config{
    Root:   myEmbeddedFiles.HTTP,
  }))

  app.Listen(8080)
}
```

### statik

```go
package main

import (
	"log
  "github.com/gofiber/fiber"
  "github.com/gofiber/embed"
	
	"<Your go module>/statik"
	"github.com/rakyll/statik/fs"
)

func main() {

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Use("/", embed.New(embed.Config{
		Root: statikFS,
	}))

	app.Listen(8080)
}
```