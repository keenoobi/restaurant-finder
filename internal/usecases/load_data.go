package usecases

import (
	"Go_Day03/internal/interfaces/csv"
	"Go_Day03/internal/interfaces/elastic"
	"fmt"
	"log"
)

type LoadDataUseCase struct {
	elasticClient elastic.ElasticClient
	csvReader     csv.CSVReader
	indexName     string
}

func NewLoadDataUseCase(elasticClient elastic.ElasticClient, csvReader csv.CSVReader, indexName string) *LoadDataUseCase {
	return &LoadDataUseCase{
		elasticClient: elasticClient,
		csvReader:     csvReader,
		indexName:     indexName,
	}
}

func (uc *LoadDataUseCase) Execute(csvFile string) error {
	// Читаем CSV-файл
	records, err := uc.csvReader.ReadCSV(csvFile)
	if err != nil {
		return fmt.Errorf("failed to read CSV: %s", err)
	}

	// Преобразуем CSV в JSON (структуры Place)
	places, err := uc.csvReader.CSVToJSON(records)
	if err != nil {
		return fmt.Errorf("failed to convert CSV to JSON: %s", err)
	}

	// Загружаем данные в Elasticsearch
	if err := uc.elasticClient.BulkIndex(uc.indexName, places); err != nil {
		return fmt.Errorf("failed to load data: %s", err)
	}

	log.Println("Data loaded successfully!")
	return nil
}
