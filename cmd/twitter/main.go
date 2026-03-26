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

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

func apiV1Router(root *chi.Mux, su service.UserService, st service.TweetService) error {
	twitterAPI := apiv1.New(su, st)

	swagger, err := openapiv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("error getting swagger: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidator(swagger)) //nolint:staticcheck // pending oapi-codegen upgrade
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	root.Mount("/api/v1", http.StripPrefix("/api/v1", openapiv1.HandlerFromMux(twitterAPI, r)))

	apiJSON, err := json.Marshal(swagger)
	if err != nil {
		return fmt.Errorf("error marshaling swagger: %w", err)
	}
	root.Get("/api/v1/api.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, writeErr := w.Write(apiJSON)
		if writeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return nil
}

func printRoutes(ctx context.Context, r chi.Router) {
	walkFunc := func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		slog.DebugContext( //nolint:sloglint // global logger configured via slog.SetDefault
			ctx,
			"registered route",
			"method",
			method,
			"route",
			route,
		)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		slog.ErrorContext( //nolint:sloglint // global logger configured via slog.SetDefault
			ctx,
			"error walking routes",
			errLogKey,
			err,
		)
	}
}

func main() {
	ctx := context.Background()
	ctx, err := InitLogFromEnv(ctx)
	if err != nil {
		panic(fmt.Sprintf("Error initializing logging: %v", err))
	}

	port := flag.Int("port", 8888, "Port for the HTTP server")
	flag.Parse()

	if runErr := run(ctx, *port); runErr != nil {
		slog.ErrorContext( //nolint:sloglint // global logger configured via slog.SetDefault
			ctx,
			"fatal error",
			errLogKey,
			runErr,
		)
		os.Exit(1)
	}
}

func run(ctx context.Context, port int) error {
	// Set up our data store
	dbServer, err := postgres.NewStorage(ctx)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer dbServer.Close()

	s := store.NewPersistentStore(dbServer.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)

	// Set up the root router
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)
	root.Use(middleware.StripSlashes)

	// Set up API v1
	if routerErr := apiV1Router(root, su, st); routerErr != nil {
		return fmt.Errorf("setting up api router: %w", routerErr)
	}

	// Print out the routes if we're in debug mode
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		printRoutes(ctx, root)
	}

	return serve(ctx, root, port)
}

func serve(ctx context.Context, handler http.Handler, port int) error {
	srv := &http.Server{
		Handler:     handler,
		Addr:        fmt.Sprintf(":%d", port),
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		ReadTimeout: serverReadTimeout,
	}

	errChan := make(chan error)
	go func() {
		slog.InfoContext(ctx, "serving http on port", "port", port)
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

	slog.InfoContext(ctx, "shutting down...") //nolint:sloglint // global logger configured via slog.SetDefault
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}
	slog.InfoContext( //nolint:sloglint // global logger configured via slog.SetDefault
		ctx,
		"Server shutdown gracefully",
	)
	return nil
}
