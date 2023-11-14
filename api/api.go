package api

import (
	"net/http"

	"app/api/handler"
	"app/config"
	"app/storage"
)

func NewApi(cfg *config.Config, storage storage.StorageI) {

	handler := handler.NewHandler(cfg, storage)

	http.HandleFunc("/category", handler.Category)
	http.HandleFunc("/product", handler.Product)
}
