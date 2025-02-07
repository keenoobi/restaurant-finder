package csv

import "Go_Day03/internal/entities"

type CSVReader interface {
	ReadCSV(filePath string) ([][]string, error)
	CSVToJSON(records [][]string) ([]entities.Place, error)
}
