package store

import "Go_Day03/internal/entities"

type Store interface {
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
	GetClosestPlaces(lat, lon float64, limit int) ([]entities.Place, error)
}
