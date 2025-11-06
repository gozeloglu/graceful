package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// CleanUpFunc is a type for batch cleanups. It should be given to
// [Graceful.RegisterCleanupFunctions] method.
type CleanUpFunc func()

// Graceful type for registering cleanup functions and graceful shutdown. You
// can define as in the example below.
//
// Example:
//
//	g := &Graceful{}
type Graceful struct {
	functions []CleanUpFunc
}

// Shutdown listens given os signals. If any of them is triggered, the given
// [http.Server] is shutdown. Finally, all cleanup functions are triggered
// after the [http.Server] is shutdown. If no signal is passed to method, it
// means that all signals are listened. [Graceful.Shutdown] blocks the code
// until notified by any given signal.
//
// Example:
//
//	srv := &http.Server{}
//	dbCleanUpFunc := func() {
//		db.Close()
//	}
//	fileCleanUpFunc := func() {
//		f.Close()
//	}
//	g := &Graceful{}
//	g.RegisterCleanupFunctions(dbCleanUpFunc, fileCleanUpFunc)
//	g.Shutdown(srv, 5*time.Second, os.Interrupt, syscall.SIGTERM)
func (g *Graceful) Shutdown(srv *http.Server, timeout time.Duration, signals ...os.Signal) error {
	defer g.cleanup()

	signalCh := make(chan os.Signal, 1)
	defer close(signalCh)

	signal.Notify(signalCh, signals...)

	<-signalCh
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

// RegisterCleanupFunctions registers cleanup functions. You can pass the
// functions that should be cleaned up after graceful shutdown. They can be
// closing database, file, or even context cancelling.
func (g *Graceful) RegisterCleanupFunctions(functions ...CleanUpFunc) {
	if g.functions == nil {
		g.functions = make([]CleanUpFunc, 0)
	}

	for _, f := range functions {
		g.functions = append(g.functions, f)
	}
}

func (g *Graceful) cleanup() {
	for _, f := range g.functions {
		f()
	}
}
