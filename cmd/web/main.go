package main

import (
	"Go_Day03/internal/config"
	"Go_Day03/internal/entities"
	"Go_Day03/internal/interfaces/store"
	"Go_Day03/internal/logger"
	"Go_Day03/internal/repositories/elasticsearch"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

const (
	pageSize = 10
)

func main() {
	configPath := flag.String("config", "internal/config/config.yaml", "Path to the config file")
	flag.Parse()

	logger := logger.NewLogrusLogger()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load config: %s", err)
	}

	elasticClient, err := elasticsearch.NewClient(cfg.Elasticsearch.Address, logger)
	if err != nil {
		logger.Fatalf("Failed to create Elasticsearch client: %s", err)
	}

	// Используем интерфейс Store
	var storeClient store.Store = elasticClient

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
			return
		}

		offset := (page - 1) * pageSize
		places, total, err := storeClient.GetPlaces(pageSize, offset)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get places: %s", err), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(cfg.WebApp.HTMLPage)
		if err != nil {
			logger.Fatalf("Can't parse index.html file: %s", err)
		}

		data := struct {
			Places   []entities.Place
			Total    int
			Page     int
			PrevPage int
			NextPage int
			LastPage int
		}{
			Places:   places,
			Total:    total,
			Page:     page,
			PrevPage: page - 1,
			NextPage: page + 1,
			LastPage: (total + pageSize - 1) / pageSize,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render template: %s", err), http.StatusInternalServerError)
		}
	})

	logger.Info(fmt.Sprintf("Starting server on %s", cfg.WebApp.Port))
	if err := http.ListenAndServe(cfg.WebApp.Port, nil); err != nil {
		logger.Fatalf("Failed to start server: %s", err)
	}
}
