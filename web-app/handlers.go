package app

import (
	"Go_Day03/internal/entities"
	"Go_Day03/internal/interfaces/store"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

const pageSize = 10

func HandlePlacesRequest(storeClient store.Store, htmlPage string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
			return
		}

		offset := (page - 1) * pageSize
		places, total, err := storeClient.GetPlaces(pageSize, offset)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get places: %s", err), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(htmlPage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse template: %s", err), http.StatusInternalServerError)
			return
		}

		data := struct {
			Places   []entities.Place
			Total    int
			Page     int
			PrevPage int
			NextPage int
			LastPage int
		}{
			Places:   places,
			Total:    total,
			Page:     page,
			PrevPage: page - 1,
			NextPage: page + 1,
			LastPage: (total + pageSize - 1) / pageSize,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render template: %s", err), http.StatusInternalServerError)
		}
	}
}
