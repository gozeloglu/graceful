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
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}
	fileCleanUpFunc := func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}

	g := &graceful.Graceful{}
	g.RegisterCleanupFunctions(dbCleanUpFunc, fileCleanUpFunc)
	err := g.Shutdown(srv, 5*time.Second, os.Interrupt, syscall.SIGTERM)
	if err != nil {
		log.Fatal(err)
	}
}
