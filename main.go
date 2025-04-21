package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Item struct {
	ShortDescription string `json:"shortDescription" validate:"required,shortDescription"`
	Price            string `json:"price" validate:"required,price"`
}

type Receipt struct {
	Retailer     string `json:"retailer" validate:"required,retailer"`
	PurchaseDate string `json:"purchaseDate" validate:"required,purchaseDate"`
	PurchaseTime string `json:"purchaseTime" validate:"required,purchaseTime"`
	Items        []Item `json:"items" validate:"required"`
	Total        string `json:"total" validate:"required,total"`
}

// TODO: May be use contexts

var validate *validator.Validate

func main() {
	validate = validator.New()
	validate := validator.New()
	validate.RegisterValidation("retailer", validateRetailer)
	validate.RegisterValidation("purchaseDate", validateDateFormat)
	validate.RegisterValidation("purchaseDate", validateTimeFormat)
	validate.RegisterValidation("total", validateTotal)
	validate.RegisterValidation("shortDescription", validateShortDescription)
	validate.RegisterValidation("price", validatePrice)

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

	err := validate.Struct(receipt)
	if err != nil {

	}

}

func validateTimeFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04", fl.Field().String())
	return err == nil
}

func validateDateFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

func validateTotal(fl validator.FieldLevel) bool {
	total := fl.Field().String()
	matched, _ := regexp.MatchString("^\\d+\\.\\d{2}$", total)
	return matched
}

func validateShortDescription(fl validator.FieldLevel) bool {
	total := fl.Field().String()
	matched, _ := regexp.MatchString("^[\\w\\s\\-]+$", total)
	return matched
}

func validatePrice(fl validator.FieldLevel) bool {
	total := fl.Field().String()
	matched, _ := regexp.MatchString("^\\d+\\.\\d{2}$", total)
	return matched
}

// Custom validation for Retailer
func validateRetailer(fl validator.FieldLevel) bool {
	retailer := fl.Field().String()
	matched, _ := regexp.MatchString("^[\\w\\s\\-&]+$", retailer)
	return matched
}
