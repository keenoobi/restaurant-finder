package utils

import (
	"Go_Day03/internal/entities"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type CSVReaderImpl struct {
	delimiter string
}

func NewCSVReader(delimiter string) *CSVReaderImpl {
	return &CSVReaderImpl{delimiter: delimiter}
}

func (r *CSVReaderImpl) ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = []rune(r.delimiter)[0]

	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read the first row: %s", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %s", err)
	}

	return records, nil
}

func (r *CSVReaderImpl) CSVToJSON(records [][]string) ([]entities.Place, error) {
	var places []entities.Place

	for _, record := range records {
		place := entities.Place{}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert id: %s", err)
		}
		place.ID = id + 1

		place.Name = record[1]
		place.Address = record[2]
		place.Phone = record[3]

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
