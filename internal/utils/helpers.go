package utils

import (
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

// CSVToJSON преобразует записи CSV в JSON-документы.
func CSVToJSON(records [][]string) ([]map[string]interface{}, error) {
	var documents []map[string]interface{}

	// Заголовки
	headers := []string{"id", "name", "address", "phone", "lon", "lat"}

	for _, record := range records {
		doc := make(map[string]interface{})
		for i, value := range record {
			doc[headers[i]] = value
		}

		// Преобразуем id в uint64
		id, err := strconv.ParseUint(doc["id"].(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert id: %s", err)
		}
		doc["id"] = id

		// Преобразуем lat и lon в float64
		lat, err := strconv.ParseFloat(doc["lat"].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert lat: %s", err)
		}

		lon, err := strconv.ParseFloat(doc["lon"].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert lon: %s", err)
		}

		// Добавляем location как geo_point
		doc["location"] = map[string]float64{
			"lat": lat,
			"lon": lon,
		}

		// Удаляем lat и lon из документа
		delete(doc, "lat")
		delete(doc, "lon")

		documents = append(documents, doc)
	}

	return documents, nil
}
