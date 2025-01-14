package main

import (
	"Go_Day03/internal/elastic"
	"log"
)

func main() {
	// Создаем клиент Elasticsearch
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatalf("error to create client: %s", err)
	}

	// Создаем индекс "places"
	if err = client.CreateIndex("places"); err != nil {
		log.Fatalf("error to create index: %s", err)
	}

	log.Println("index 'places' successfuly created")

	// Схема маппинга
	mapping := `{
		"properties": {
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

	if err = client.AddMapping("places", mapping); err != nil {
		log.Fatalf("error to add mapping")
	}

	log.Println("mapping successfuly added")
}
