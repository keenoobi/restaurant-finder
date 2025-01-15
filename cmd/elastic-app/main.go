package main

import (
	"Go_Day03/internal/elastic"
	"Go_Day03/internal/utils"
	"flag"
	"log"
)

func main() {
	// Флаги командной строки
	createIndex := flag.Bool("create-index", false, "Create the 'places' index")
	addMapping := flag.Bool("add-mapping", false, "Add mapping to the 'places' index")
	loadData := flag.Bool("load-data", false, "Load data from CSV into the 'places' index")
	csvFile := flag.String("csv", "data.csv", "Path to the CSV file with data")

	flag.Parse()

	// Создаем клиент Elasticsearch
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatalf("error to create client: %s", err)
	}

	// Создаем индекс "places"
	if *createIndex {
		if err := client.CreateIndex("places"); err != nil {
			log.Fatalf("Failed to create index: %s", err)
		}
		log.Println("Index 'places' created successfully!")
	}

	if *addMapping {
		// Схема маппинга
		mapping := `{
		  "properties": {
			"id": {
			  "type": "unsigned_long"
			},
			"name": {
			  "type": "text"
			},
			"address": {
			  "type": "text"
			},
			"phone": {
			  "type": "text"
			},
			"location": {
			  "type": "geo_point"
			}
		  }
		}`

		if err := client.AddMapping("places", mapping); err != nil {
			log.Fatalf("Failed to add mapping: %s", err)
		}
		log.Println("Mapping added successfully!")
	}

	// Загрузка данных
	if *loadData {
		// Читаем CSV-файл
		records, err := utils.ReadCSV(*csvFile)
		if err != nil {
			log.Fatalf("Failed to read CSV: %s", err)
		}

		// Преобразуем CSV в JSON
		places, err := utils.CSVToJSON(records)
		if err != nil {
			log.Fatalf("Failed to convert CSV to JSON: %s", err)
		}

		// Загружаем данные в Elasticsearch
		if err := client.BulkIndex("places", places); err != nil {
			log.Fatalf("Failed to load data: %s", err)
		}
		log.Println("Data loaded successfully!")
	}
}
