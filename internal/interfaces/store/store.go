package store

import "Go_Day03/internal/entities"

type Store interface {
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
}
