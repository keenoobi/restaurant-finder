package utils

import (
	"Go_Day03/internal/models"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ReadCSV читает CSV-файл и возвращает его записи.
func ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' // Используем табуляцию как разделитель

	// Пропускаем первую строку (заголовки)
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read the first row: %s", err)
	}

	// Читаем оставшиеся строки
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %s", err)
	}

	return records, nil
}

// CSVToJSON преобразует записи CSV в слайс структур Place.
func CSVToJSON(records [][]string) ([]models.Place, error) {
	var places []models.Place

	for _, record := range records {
		// Создаем структуру Place
		place := models.Place{}

		// Преобразуем id в int
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert id: %s", err)
		}
		place.ID = id

		// Заполняем остальные поля
		place.Name = record[1]
		place.Address = record[2]
		place.Phone = record[3]

		// Преобразуем lat и lon в float64
		lat, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert lat: %s", err)
		}
		place.Location.Lat = lat

		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert lon: %s", err)
		}
		place.Location.Lon = lon

		places = append(places, place)
	}

	return places, nil
}
