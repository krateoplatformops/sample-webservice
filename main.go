package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/krateoplatformops/sample-webservice/docs"
	"github.com/krateoplatformops/sample-webservice/internal/auth"
	"github.com/krateoplatformops/sample-webservice/internal/handlers"
	create "github.com/krateoplatformops/sample-webservice/internal/handlers/create"
	delete "github.com/krateoplatformops/sample-webservice/internal/handlers/delete"
	"github.com/krateoplatformops/sample-webservice/internal/handlers/docs"
	update "github.com/krateoplatformops/sample-webservice/internal/handlers/update"

	get "github.com/krateoplatformops/sample-webservice/internal/handlers/get"
	list "github.com/krateoplatformops/sample-webservice/internal/handlers/list"

	"github.com/krateoplatformops/snowplow/plumbing/env"
	prettylog "github.com/krateoplatformops/snowplow/plumbing/slogs/pretty"
)

var (
	serviceName = "sample-webservice"
)

// @title 		 Sample Webservice API
// @version         1.0
// @description   Sample Webservice API.
// @BasePath		/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	debugOn := flag.Bool("debug", env.Bool("PLUGIN_DEBUG", true), "dump verbose output")
	port := flag.Int("port", env.Int("PLUGIN_PORT", 8081), "port to listen on")

	flag.Parse()

	mux := http.NewServeMux()

	logLevel := slog.LevelInfo
	if *debugOn {
		logLevel = slog.LevelDebug
	}

	lh := prettylog.New(&slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	},
		prettylog.WithDestinationWriter(os.Stderr),
		prettylog.WithColor(),
		prettylog.WithOutputEmptyAttrs(),
	)
	log := slog.New(lh)

	log = log.With("service", serviceName)

	opts := handlers.HandlerOptions{
		Log: log,
	}

	healthy := int32(0)

	mux.Handle("/resource", auth.BearerAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get.Get(opts).ServeHTTP(w, r)
		case http.MethodPost:
			create.Create(opts).ServeHTTP(w, r)
		case http.MethodDelete:
			delete.Delete(opts).ServeHTTP(w, r)
		case http.MethodPatch:
			update.Update(opts).ServeHTTP(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/resources", auth.BearerAuthMiddleware(list.List(opts)))
	mux.Handle("/openapi/", docs.WrapHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 50 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}...)
	defer stop()

	go func() {
		atomic.StoreInt32(&healthy, 1)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("could not listen on %s - %v", server.Addr, err)
		}
	}()

	// Listen for the interrupt signal.
	log.Info("server is ready to handle requests", slog.Any("port", *port))
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info("server is shutting down gracefully, press Ctrl+C again to force")
	atomic.StoreInt32(&healthy, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	}

	log.Info("server gracefully stopped")
}
