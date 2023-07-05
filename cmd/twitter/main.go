package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	apiv1 "github.com/ricleal/twitter-clone/internal/api/v1"
	openapiv1 "github.com/ricleal/twitter-clone/internal/api/v1/openapi"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

// apiV1Router sets up the API v1 router and mounts it on the root router.
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
	r.Use(oapiMiddleware.OapiRequestValidator(swagger))
	r.Use(middleware.AllowContentType("application/json")) //nolint:goconst // ignore
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	root.Mount("/api/v1", http.StripPrefix("/api/v1", openapiv1.HandlerFromMux(twitterAPI, r)))

	apiJSON, err := json.Marshal(swagger)
	if err != nil {
		return fmt.Errorf("error marshaling swagger: %w", err)
	}
	root.Get("/api/v1/api.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(apiJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	return nil
}

// printRoutes prints out the routes registered on a chi router.
func printRoutes(ctx context.Context, r chi.Router) {
	walkFunc := func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Ctx(ctx).Debug().Str("method", method).Str("route", route).Msg("registered route")
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error walking routes")
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

	// Set up our data store
	dbServer, err := postgres.NewStorage(ctx)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error connecting to database")
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
	if err := apiV1Router(root, su, st); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error setting up openapi router")
	}

	// Print out the routes if we're in debug mode
	if log.Ctx(ctx).GetLevel() <= zerolog.DebugLevel {
		printRoutes(ctx, root)
	}

	// Start the server
	if err := serve(ctx, root, *port); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error serving http")
	}
}

// serve starts the HTTP server and handles shutdown gracefully.
func serve(ctx context.Context, handler http.Handler, port int) error {
	srv := &http.Server{
		Handler:     handler,
		Addr:        fmt.Sprintf(":%d", port),
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		ReadTimeout: 10 * time.Second,
	}

	errChan := make(chan error)
	go func() {
		log.Ctx(ctx).Info().Int("port", port).Msg("serving http on port")
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

	log.Ctx(ctx).Info().Msg("shutting down...")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown gracefully: %w", err)
	}
	log.Ctx(ctx).Info().Msg("Server shutdown gracefully")
	return nil
}
