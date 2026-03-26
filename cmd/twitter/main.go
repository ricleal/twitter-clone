package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	oapiMiddleware "github.com/oapi-codegen/nethttp-middleware"

	apiv1 "github.com/ricleal/twitter-clone/internal/api/v1"
	openapiv1 "github.com/ricleal/twitter-clone/internal/api/v1/openapi"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

const errLogKey = "error" // slog attribute key for error values

const (
	serverReadTimeout = 10 * time.Second // HTTP server read timeout
	shutdownTimeout   = 10 * time.Second // graceful shutdown deadline
)

func apiV1Router(root *http.ServeMux, logger *slog.Logger, su service.UserService, st service.TweetService) error {
	twitterAPI := apiv1.New(logger, su, st)

	swagger, err := openapiv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("error getting swagger: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	apiJSON, err := json.Marshal(swagger)
	if err != nil {
		return fmt.Errorf("error marshaling swagger: %w", err)
	}

	openapiv1.HandlerWithOptions(twitterAPI, openapiv1.StdHTTPServerOptions{
		BaseURL:    "/api/v1",
		BaseRouter: root,
		Middlewares: []openapiv1.MiddlewareFunc{
			oapiMiddleware.OapiRequestValidator(swagger),
			middleware.AllowContentType("application/json"),
			middleware.SetHeader("Content-Type", "application/json"),
		},
	})

	root.HandleFunc("GET /api/v1/api.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, writeErr := w.Write(apiJSON)
		if writeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return nil
}

func main() {
	ctx := context.Background()
	logger, err := InitLogFromEnv()
	if err != nil {
		panic(fmt.Sprintf("Error initializing logging: %v", err))
	}

	port := flag.Int("port", 8888, "Port for the HTTP server")
	flag.Parse()

	if runErr := run(ctx, logger, *port); runErr != nil {
		logger.Error("fatal error", errLogKey, runErr)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger, port int) error {
	// Set up our data store
	dbServer, err := postgres.NewStorage(ctx, logger)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer dbServer.Close()

	s := store.NewPersistentStore(dbServer.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)

	// Set up the root mux
	mux := http.NewServeMux()

	// Set up API v1
	if routerErr := apiV1Router(mux, logger, su, st); routerErr != nil {
		return fmt.Errorf("setting up api router: %w", routerErr)
	}

	// Wrap with global middleware (outermost = first to run)
	var h http.Handler = mux
	h = middleware.StripSlashes(h)
	h = middleware.Recoverer(h)
	h = middleware.Logger(h)

	return serve(ctx, logger, h, port)
}

func serve(ctx context.Context, logger *slog.Logger, handler http.Handler, port int) error {
	srv := &http.Server{
		Handler:     handler,
		Addr:        fmt.Sprintf(":%d", port),
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		ReadTimeout: serverReadTimeout,
	}

	errChan := make(chan error)
	go func() {
		logger.InfoContext(ctx, "serving http on port", "port", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %w", err)
		}
	}()
	ctx, stop := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}

	logger.InfoContext(ctx, "shutting down...")
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}
	logger.InfoContext(ctx, "Server shutdown gracefully")
	return nil
}
