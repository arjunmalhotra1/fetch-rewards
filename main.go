package main

import (
	"net/http"

	"github.com/arjunmalhotra1/fetch-rewards/handler"
	"github.com/arjunmalhotra1/fetch-rewards/repository"
	"github.com/arjunmalhotra1/fetch-rewards/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo := repository.NewReceiptRepository()
	service := service.NewReceiptService(repo)
	handler := handler.NewHandler(service)

	router := chi.NewMux()
	router.Get("/", handler.Healthz)
	router.Post("/receipts/process", handler.ProcessReceipts)
	router.Get("/receipts/{id}/points", handler.GetPoints)

	http.ListenAndServe(":8086", router)
}
