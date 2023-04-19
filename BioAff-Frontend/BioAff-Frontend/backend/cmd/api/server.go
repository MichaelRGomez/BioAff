// BIOAFF/backend/cmd/api/server.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {

	//http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//The shutdown() function should return its error to this channel
	shutdownError := make(chan error)

	//start a background Goroutine
	go func() {

		//create a quit/exit channel which carries os.Signal values
		quit := make(chan os.Signal, 1)

		//listen for SIGINT and SIGTERM signals
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		//block until a signal is received
		s := <-quit

		//log a message
		app.logger.PrintInfo("shuttding down server", map[string]string{
			"signal": s.String(),
		})

		//create a context with a 20-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		//call the shutdown function
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		//log a message about that Goroutine
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		app.wg.wait()
		shutdownError <- nil
	}()

	//starting our server
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	//check if the shutdown process has been initiated
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	//block for notification from shutdown() function
	err = <-shutdownError
	if err != nil {
		return err
	}

	//graceful shutdown was successful
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
