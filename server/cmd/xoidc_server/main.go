package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/server/internal/pkg/api"
	"github.com/zltl/xoidc/server/internal/pkg/exampleop"
	"github.com/zltl/xoidc/server/internal/pkg/storage"
	"golang.org/x/exp/slog"
)

func main() {
	//we will run on :9998
	port := "9998"
	//which gives us the issuer: http://localhost:9998/
	issuer := fmt.Sprintf("http://localhost:%s/", port)

	log.SetFormatter(&log.TextFormatter{
		DisableQuote:  true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)

	log.SetLevel(log.TraceLevel)

	// the OpenIDProvider interface needs a Storage interface handling various checks and state manipulations
	// this might be the layer for accessing your database
	// in this example it will be handled in-memory
	storage := &storage.Storage{
		PGHost:     "localhost",
		PGPort:     5432,
		PGUsername: "postgres",
		PGPassword: "postgres",
		PGDBName:   "xoidc",
	}

	err := storage.Open()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)

	router := exampleop.SetupServer(issuer, storage, logger, false)
	h := api.Handler{
		Store: storage,
	}
	router.Route("/api/oidc", h.Serve)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("server listening on http://localhost:%s/", port)
	log.Println("press ctrl+c to stop")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
