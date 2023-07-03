package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/ricleal/twitter-clone/internal/api"
	"github.com/ricleal/twitter-clone/internal/api/openapiv1"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/store"
)

func openAPIRouter(root *chi.Mux, twitterServer *api.TwitterAPI) *chi.Mux {

	swagger, err := openapiv1.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	r := chi.NewRouter()
	r.Use(oapiMiddleware.OapiRequestValidator(swagger))

	root.Mount("/api/v1", http.StripPrefix("/api/v1", openapiv1.HandlerFromMux(twitterServer, r)))

	apiJSON, _ := json.Marshal(swagger)
	root.Get("/api/v1/api.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(apiJSON)
	})

	return r

}

func main() {
	ctx := context.Background()
	ctx, err := InitLogFromEnv(ctx)
	if err != nil {
		panic(fmt.Sprintf("Error initializing logging: %v", err))
	}

	port := flag.Int("port", 8889, "Port for the HTTP server")
	flag.Parse()

	// Set up our data store
	dbServer, err := postgres.NewHandler(ctx)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error connecting to database")
	}
	defer dbServer.Close()
	s := store.NewSQLStore(dbServer.DB())
	st := service.NewTweetService(s)
	su := service.NewUserService(s)
	twitterAPI := api.New(su, st)

	// Set up router
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)
	root.Use(middleware.StripSlashes)

	openAPIRouter(root, twitterAPI)

	// TODO: Delete this
	// DEBUG: Print out all routes
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(root, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	// Start the server
	if err := serve(ctx, root, *port); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error serving http")
	}

}

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
