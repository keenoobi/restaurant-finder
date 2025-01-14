package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func ReadCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("fail to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("fail to read CSV: %s", err)
	}

	return records, nil
}

func CSVToJSON(records [][]string) ([]map[string]interface{}, error) {
	var documents []map[string]interface{}

	// Первая строка - это заголовки
	headers := records[0]

	for _, record := range records[1:] {
		doc := make(map[string]interface{})
		for i, value := range record {
			doc[headers[i]] = value
		}

		// Преобразуем lat и lon в float64
		lat, err := strconv.ParseFloat(doc["lat"].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("fail to transform lat: %s", err)
		}

		lon, err := strconv.ParseFloat(doc["lon"].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("fail to transform lon: %s", err)
		}

		// Добавляем location как geo_point
		doc["location"] = map[string]float64{
			"lat": lat,
			"lon": lon,
		}

		delete(doc, "lat")
		delete(doc, "lon")

		documents = append(documents, doc)
	}

	return documents, nil
}
