package main

import (
	"Go_Day03/internal/config"
	"Go_Day03/internal/mylogrus"
	"Go_Day03/internal/repositories/elasticsearch"
	"Go_Day03/internal/usecases"
	"Go_Day03/internal/utils"
	"encoding/json"
	"flag"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to the config file")
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

	csvReader := utils.NewCSVReader(cfg.CSV.Delimiter)

	if err := elasticClient.CreateIndex(cfg.Elasticsearch.Index); err != nil {
		logger.Fatalf("Failed to create index: %s", err)
	}

	// Преобразуем маппинг из YAML в JSON
	mappingJSON, err := json.Marshal(cfg.Elasticsearch.Mapping)
	if err != nil {
		logger.Fatalf("Failed to serialize mapping to JSON: %s", err)
	}

	// Передаем JSON в Elasticsearch
	if err := elasticClient.AddMapping(cfg.Elasticsearch.Index, string(mappingJSON)); err != nil {
		logger.Fatalf("Failed to add mapping: %s", err)
	}

	loadDataUseCase := usecases.NewLoadDataUseCase(elasticClient, csvReader, logger, cfg.Elasticsearch.Index)
	if err := loadDataUseCase.Execute(cfg.CSV.FilePath); err != nil {
		logger.Fatalf("Failed to load data: %s", err)
	}
}
