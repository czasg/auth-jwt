# go-fill
[![LICENSE](https://img.shields.io/github/license/mashape/apistatus.svg?style=flat-square&label=License)](https://github.com/czasg/go-fill/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/czasg/go-fill/branch/main/graph/badge.svg?token=OkiSH6DMqf)](https://codecov.io/gh/czasg/go-fill)
[![GitHub Stars](https://img.shields.io/github/stars/czasg/go-fill.svg?style=flat-square&label=Stars&logo=github)](https://github.com/czasg/go-fill/stargazers)

Auto filling zero-value、env-value、default-value into a pointer-struct.

## zero-value
use `fill.Fill`, it can fill most zero-value.
```go
package main

import (
	"fmt"
	"github.com/czasg/go-fill"
)

type Response struct {
	Code    int
	Message string
	Data    *Data
}

type Data struct {
	Trace string
	Body  []byte
}

func main() {
	response := Response{}   // response.Data is nil.
	_ = fill.Fill(&response) // response.Data is zero-value.
	fmt.Println(response.Data == nil)
	fmt.Println(response.Data.Body == nil)
}
```

## env-value
`fill.FillEnv` is equal of `fill.Fill(v, fill.OptEnv)`.
```go
package main

import (
	"fmt"
	"github.com/czasg/go-fill"
	"os"
)

type Config struct {
	Env     string
	Version string
	Postgres
}

type Postgres struct {
	Addr     string
	User     string
	Password string
	Database string
}

func main() {
	_ = os.Setenv("ENV", "test")
	_ = os.Setenv("VERSION", "v0.0.1")
	_ = os.Setenv("POSTGRES_ADDR", "localhost:5432")
	_ = os.Setenv("POSTGRES_USER", "postgres")
	_ = os.Setenv("POSTGRES_PASSWORD", "postgres")
	_ = os.Setenv("POSTGRES_DATABASE", "postgres")
	cfg := Config{}
	_ = fill.FillEnv(&cfg)
	fmt.Println(cfg.Env, cfg.Version, cfg.Postgres)
}
```

### env tag list
|tag|comment|
|---|---|
|env:"fieldName"|default is struct field name, you can also point a new fieldName.|
|env:",default=value"|set default env value.|
|env:",require"|it return an err when env is not found.|
|env:",empty"|set current fieldName to an empty string like "".|
|env:",sep=_"|when struct into struct, sep is the connector, default is "_".|

```go
package main

import (
	"fmt"
	"github.com/czasg/go-fill"
	"os"
)

type Config struct {
	RPC      `env:"GRPC"`
	Redis    `env:"RDS"`
	Postgres `env:"PG"`
}

type RPC struct {
	gRPC `env:",empty"`
}

type gRPC struct {
	Addr string
}

type Redis struct {
	Addr     string `env:",sep=!"`
	Password string `env:",sep=@"`
	DB       int    `env:",sep=#"`
}

type Postgres struct {
	Addr     string `env:",default=localhost:5432"`
	User     string `env:",default=postgres"`
	Password string `env:",default=postgres"`
	Database string `env:",default=postgres"`
}

func main() {
	_ = os.Setenv("GRPC_ADDR", "localhost:9000")
	_ = os.Setenv("RDS!ADDR", "localhost:6379")
	_ = os.Setenv("RDS@PASSWORD", "123456")
	_ = os.Setenv("RDS#DB", "10")
	cfg := Config{}
	_ = fill.FillEnv(&cfg)
	fmt.Println(cfg)
}
```


## default-value
`fill.FillDefault` is equal of `fill.Fill(v, fill.OptDefault)`.
```go
package main

import (
	"fmt"
	"github.com/czasg/go-fill"
)

type Request struct {
	Trace      string `default:"no-trace"`
	PageSize   int    `default:"20"`
	PageNumber int    `default:"1"`
}

func main() {
	request := Request{}
	request.Trace = "no-zero-value will not be filled with default-value"
	_ = fill.FillDefault(&request)
	fmt.Println(request)
}
```
