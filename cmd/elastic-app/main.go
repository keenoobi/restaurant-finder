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
	if err := client.CreateIndex("places"); err != nil {
		log.Fatalf("error to create index: %s", err)
	}

	log.Println("index 'places' successfully created")
}
