package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zltl/xoidc/internal/pkg/db"
	"github.com/zltl/xoidc/internal/pkg/exampleop"
	"github.com/zltl/xoidc/internal/pkg/storage"
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

	mdb := db.New("localhost", 5432, "postgres", "postgres", "xoidc")
	err := mdb.Open()
	if err != nil {
		log.Fatal(err)
	}

	// the OpenIDProvider interface needs a Storage interface handling various checks and state manipulations
	// this might be the layer for accessing your database
	// in this example it will be handled in-memory
	ustore := storage.NewUserStore(issuer)
	storage := storage.NewStorage(ustore, mdb)

	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)

	router := exampleop.SetupServer(issuer, storage, logger, false)

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
