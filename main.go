package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

func main() {
	validator.New()
	router := chi.NewMux()

	router.Get("/", healthz)
	router.Post("/receipts/process", processReceipts)
	http.ListenAndServe(":8086", router)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"status": "ok"})
}

func processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	receiptBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(receiptBody, &receipt)

	fmt.Println(string(receiptBody))

}

func validateReceipt() {

}
