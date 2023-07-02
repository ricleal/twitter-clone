package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/ricleal/twitter-clone/internal/api"
	"github.com/ricleal/twitter-clone/internal/api/openapi"
	"github.com/ricleal/twitter-clone/internal/service"
	"github.com/ricleal/twitter-clone/internal/service/repository/postgres"
	"github.com/ricleal/twitter-clone/internal/service/store"
	"github.com/rs/zerolog/log"
)

func main() {

	ctx := context.Background()
	ctx, err := InitLogFromEnv(ctx)
	if err != nil {
		panic(fmt.Sprintf("Error initializing logging: %v", err))
	}

	var port = flag.Int("port", 8889, "Port for the HTTP server")
	flag.Parse()

	swagger, err := openapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Set up our data store
	dbServer, err := postgres.NewHandler(ctx)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error connecting to database")
	}
	defer dbServer.Close()
	store := store.NewSQLStore(dbServer.DB())
	st := service.NewTweetService(store)
	su := service.NewUserService(store)
	twitterServer := api.New(su, st)

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// Middleware logging
	// Logger
	logger := httplog.NewLogger("httplog-example", httplog.Options{
		JSON: true,
	})
	r.Use(httplog.RequestLogger(logger))

	// register the http handlers
	openapi.HandlerFromMux(twitterServer, r)

	s := &http.Server{
		Handler:     r,
		Addr:        fmt.Sprintf(":%d", *port),
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// And we serve HTTP until the world ends.
	log.Ctx(ctx).Info().Int("port", *port).Msg("serving http on port")
	log.Ctx(ctx).Fatal().Err(s.ListenAndServe()).Msg("http server quit")
	// log.Fatal(s.ListenAndServe())
}
