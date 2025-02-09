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

	logger.Info(fmt.Sprintf("Starting server on %s", cfg.WebApp.Port))
	if err := http.ListenAndServe(cfg.WebApp.Port, nil); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}

	return nil
}
