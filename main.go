package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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
	UUID         string
	Points       int
}

type ValidationErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

var validate *validator.Validate

var pointsMap map[string]Receipt

func main() {
	validate = validator.New()
	//validate := validator.New()
	validate.RegisterValidation("retailer", validateRetailer)
	validate.RegisterValidation("purchaseDate", validateDateFormat)
	validate.RegisterValidation("purchaseTime", validateTimeFormat)
	validate.RegisterValidation("total", validateTotal)
	validate.RegisterValidation("shortDescription", validateShortDescription)
	validate.RegisterValidation("price", validatePrice)

	router := chi.NewMux()

	pointsMap = map[string]Receipt{}

	router.Get("/", healthz)
	router.Post("/receipts/process", processReceipts)
	router.Get("/receipts/{id}/points", getPoints)
	http.ListenAndServe(":8086", router)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"status": "ok"})
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	points := pointsMap[id].Points

	render.JSON(w, r, map[string]int{"points": points})
}

func buildValidationErrorResponse(errs validator.ValidationErrors) ValidationErrorResponse {
	errors := []FieldError{}

	for _, err := range errs {
		errors = append(errors, FieldError{
			Field:   err.Field(),
			Rule:    err.Tag(),
			Message: fmt.Sprintf("%s failed on the '%s' rule. The receipt is invalid. Please refer the API documentation", err.Field(), err.Tag()),
		})
	}

	return ValidationErrorResponse{Errors: errors}

}

func processReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	receiptBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(receiptBody, &receipt)

	fmt.Println(string(receiptBody))

	err := validate.Struct(receipt)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errResponse := buildValidationErrorResponse(validationErrors)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, errResponse)
			return
		}
	}

	uuid, _ := uuid.NewUUID()

	receipt.UUID = uuid.String()
	receipt.Points = calculatePoints(receipt)
	pointsMap[receipt.UUID] = receipt

	// Return the id
	fmt.Println("receipt.Points: ", receipt.Points, receipt.UUID)
}

func calculateRetailerPoints(retailer string) int {
	points := 0
	for _, v := range retailer {
		if unicode.IsLetter(v) || unicode.IsDigit(v) {
			points++
		}
	}
	fmt.Println("calculateRetailerPoints", points)
	return points
}

func calculatePoints(receipt Receipt) int {
	totalPoints := 0
	totalPoints += calculateRetailerPoints(receipt.Retailer)

	if isTotalRoundDollar(receipt.Total) {
		// TODO make these constants on the top
		fmt.Println("isTotalRoundDollar adding 50")
		totalPoints += 50
	}

	if isTotalMultiple25(receipt.Total) {
		fmt.Println("isTotalMultiple25 adding 25")
		totalPoints += 25
	}

	m := itemsMultiple(len(receipt.Items))
	fmt.Println("itemsMultiple adding", 5*m)
	totalPoints += 5 * m

	totalPoints += itemsPoints(receipt.Items)

	if isDayOdd(receipt.PurchaseDate) {
		fmt.Println("isDayOdd adding", 6)
		totalPoints += 6
	}

	if isTimeInRange(receipt.PurchaseTime) {
		fmt.Println("isTimeInRange adding", 10)
		totalPoints += 10
	}

	return totalPoints
}

func itemsPoints(items []Item) int {
	points := 0
	for _, item := range items {
		item.ShortDescription = strings.TrimSpace(item.ShortDescription)
		if len(item.ShortDescription)%3 == 0 {
			itemP, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(itemP * 0.2))
		}
	}

	fmt.Println("itemPoints", points)
	return points
}

func isTimeInRange(purchaseTime string) bool {
	t, _ := time.Parse("15:04", purchaseTime)

	h := t.Hour() * 100
	m := t.Minute()
	total := h + m

	return total > 1400 && total < 1600
}

func isDayOdd(purchaseDate string) bool {
	t, _ := time.Parse("2006-01-02", purchaseDate)
	return t.Day()%2 == 1

}

func itemsMultiple(length int) int {
	return length / 2
}

func isTotalMultiple25(total string) bool {
	amount, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return false
	}

	return math.Mod(amount, 0.25) == 0
}

func isTotalRoundDollar(total string) bool {
	amount, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return false
	}
	return amount == math.Floor(amount)
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
