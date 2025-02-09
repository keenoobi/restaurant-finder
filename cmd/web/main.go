package main

import (
	"Go_Day03/internal/config"
	"Go_Day03/internal/mylogrus"
	"Go_Day03/internal/repositories/elasticsearch"
	app "Go_Day03/web-app"
	"flag"
)

func main() {
	configPath := flag.String("config", "internal/config/config.yaml", "Path to the config file")
	flag.Parse()

	logger := mylogrus.New()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load config: %s", err)
	}

	elasticClient, err := elasticsearch.NewClient(cfg.Elasticsearch.Address, logger)
	if err != nil {
		logger.Fatalf("Failed to create Elasticsearch client: %s", err)
	}

	if err := app.StartServer(cfg, logger, elasticClient); err != nil {
		logger.Fatalf("Failed to start server: %s", err)
	}
}
