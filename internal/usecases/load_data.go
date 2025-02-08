package usecases

import (
	"Go_Day03/internal/interfaces/csv"
	"Go_Day03/internal/interfaces/elastic"
	"Go_Day03/internal/interfaces/logger"
	"fmt"
)

type LoadDataUseCase struct {
	elasticClient elastic.ElasticClient
	csvReader     csv.CSVReader
	logger        logger.Logger
	indexName     string
}

func NewLoadDataUseCase(elasticClient elastic.ElasticClient, csvReader csv.CSVReader, logger logger.Logger, indexName string) *LoadDataUseCase {
	return &LoadDataUseCase{
		elasticClient: elasticClient,
		csvReader:     csvReader,
		logger:        logger,
		indexName:     indexName,
	}
}

func (uc *LoadDataUseCase) Execute(csvFile string) error {
	uc.logger.Info("Reading CSV file")
	records, err := uc.csvReader.ReadCSV(csvFile)
	if err != nil {
		uc.logger.WithFields(map[string]interface{}{
			"file": csvFile,
		}).Errorf("Failed to read CSV: %s", err)
		return fmt.Errorf("failed to read CSV: %s", err)
	}

	uc.logger.Info("Converting CSV to JSON")
	places, err := uc.csvReader.CSVToJSON(records)
	if err != nil {
		uc.logger.WithFields(map[string]interface{}{
			"file": csvFile,
		}).Errorf("Failed to convert CSV to JSON: %s", err)
		return fmt.Errorf("failed to convert CSV to JSON: %s", err)
	}

	uc.logger.WithFields(map[string]interface{}{
		"index": uc.indexName,
	}).Info("Loading data into Elasticsearch")

	if err := uc.elasticClient.BulkIndex(uc.indexName, places); err != nil {
		uc.logger.WithFields(map[string]interface{}{
			"index": uc.indexName,
		}).Errorf("Failed to load data: %s", err)
		return fmt.Errorf("failed to load data: %s", err)
	}

	uc.logger.Info("Data loaded successfully")
	return nil
}
