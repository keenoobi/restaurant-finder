package app

import (
	"Go_Day03/internal/config"
	"Go_Day03/internal/interfaces/logger"
	"Go_Day03/internal/interfaces/store"

	"fmt"
	"net/http"
)

func StartServer(cfg *config.Config, logger logger.Logger, storeClient store.Store) error {
	http.HandleFunc("/", HandlePlacesRequest(storeClient, cfg.WebApp.HTMLPage))

	http.HandleFunc("/api/places", HandlePlacesAPI(storeClient))

	// Регистрируем защищенный обработчик для поиска ближайших ресторанов
	http.HandleFunc("/api/recommend", JWTMiddleware(cfg, HandleRecommendRequest(storeClient)))

	// Регистрируем обработчик для генерации токена
	http.HandleFunc("/api/get_token", HandleGetToken(cfg))

	logger.Info(fmt.Sprintf("Starting server on %s", cfg.WebApp.Port))
	if err := http.ListenAndServe(cfg.WebApp.Port, nil); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}

	return nil
}
