# graceful

`graceful` is a package to apply graceful shutdown for Go projects. It provides cleanup option to cleaning something while
shutting down the server.

## Installation

```shell
go get github.com/gozeloglu/graceful
```

## Example

```go
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gozeloglu/graceful"
)

func main() {
	db := sql.DB{}
	f := os.File{}
	srv := &http.Server{}

	// Create cleanup functions
	dbCleanUpFunc := func() {
		db.Close()
	}
	fileCleanUpFunc := func() {
		f.Close()
	}

	g := &graceful.Graceful{}
	g.RegisterCleanupFunctions(dbCleanUpFunc, fileCleanUpFunc)
	err := g.Shutdown(srv, 5*time.Second, os.Interrupt, syscall.SIGTERM)
	if err != nil {
		log.Fatal(err)
	}
}

```

## LICENSE
[Apache License](LICENSE)