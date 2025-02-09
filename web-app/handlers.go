package app

import (
	"Go_Day03/internal/entities"
	"Go_Day03/internal/interfaces/store"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

const pageSize = 10

type PlacesResponse struct {
	Name     string           `json:"name"`
	Total    int              `json:"total"`
	Places   []entities.Place `json:"places"`
	Page     int              `json:"page"`
	PrevPage int              `json:"prev_page"`
	NextPage int              `json:"next_page"`
	LastPage int              `json:"last_page"`
}

// getPlacesData получает данные для списка мест
func getPlacesData(storeClient store.Store, page int) (*PlacesResponse, error) {
	offset := (page - 1) * pageSize
	places, total, err := storeClient.GetPlaces(pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get places: %s", err)
	}

	lastPage := (total + pageSize - 1) / pageSize

	return &PlacesResponse{
		Name:     "Places",
		Total:    total,
		Places:   places,
		Page:     page,
		PrevPage: page - 1,
		NextPage: page + 1,
		LastPage: lastPage,
	}, nil
}

// HandlePlacesRequest обрабатывает запросы для HTML-страницы
func HandlePlacesRequest(storeClient store.Store, htmlPage string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
			return
		}

		data, err := getPlacesData(storeClient, page)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get places: %s", err), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(htmlPage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render template: %s", err), http.StatusInternalServerError)
		}
	}
}

// HandlePlacesAPI обрабатывает запросы для API
func HandlePlacesAPI(storeClient store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Invalid 'page' value: '%s'", pageStr),
			})
			return
		}

		data, err := getPlacesData(storeClient, page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to get places: %s", err),
			})
			return
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Failed to encode response: %s", err),
			})
			return
		}
	}
}
