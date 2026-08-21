package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"stickerchallenge/internal/config"
	"stickerchallenge/internal/domain"
	"stickerchallenge/internal/httpapi"
	"stickerchallenge/internal/service"
	"stickerchallenge/internal/store"
)

func main() {
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := config.Validate(cfg); err != nil {
		log.Fatal(err)
	}
	db, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	svc := service.New(db, service.FixedClock{Value: "2116-05-01T00:00:00Z"})
	if cfg.Demo {
		runDemo(svc)
		return
	}
	fmt.Printf("sticker challenge listening on %s\n", cfg.Listen)
	if err := http.ListenAndServe(cfg.Listen, httpapi.New(svc).Handler()); err != nil {
		log.Fatal(err)
	}
}

func runDemo(svc *service.Service) {
	batch, err := svc.RegisterBatch("2116-05", "Boundary divisibility", "operator", []domain.Candidate{{ID: "a", Number: 22}, {ID: "b", Number: 25}})
	if err != nil {
		log.Fatal(err)
	}
	if _, err = svc.StartReview(batch.ID, "reviewer"); err != nil {
		log.Fatal(err)
	}
	if _, err = svc.ConfirmBatch(batch.ID, "reviewer"); err != nil {
		log.Fatal(err)
	}
	if _, err = svc.PublishBatch(batch.ID, "operator"); err != nil {
		log.Fatal(err)
	}
	snapshot, err := svc.ExportConfirmed(batch.ID, "operator")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(snapshot.Payload)
}
