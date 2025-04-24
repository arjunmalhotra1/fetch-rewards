package service

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/arjunmalhotra1/fetch-rewards/repository"
	internalValidator "github.com/arjunmalhotra1/fetch-rewards/service/validator"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type ReceiptService struct {
	rr        repository.ReceiptRepository
	validator *validator.Validate
}

func NewReceiptService(repository repository.ReceiptRepository) ReceiptService {
	v := internalValidator.NewValidator()
	rs := ReceiptService{
		rr:        repository,
		validator: v,
	}
	return rs
}

func (rs ReceiptService) GetPoints(uuid string) (int, bool) {
	return rs.rr.GetPoints(uuid)
}

func (rs ReceiptService) ProcessReceipts(receipt repository.Receipt) (string, error) {
	err := rs.validator.Struct(receipt)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errResponse := internalValidator.BuildValidationErrorResponse(validationErrors)
			log.Println("validation failed for the receipt: ", errResponse)
			return "", errResponse
		}
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Println("error generating the uuid", err)
		return "", err
	}

	receipt.UUID = uuid.String()
	receipt.Points = calculatePoints(receipt)
	rs.rr.StoreReceipt(receipt.UUID, receipt)

	log.Println("receipt Points: ", receipt.Points, receipt.UUID)
	return receipt.UUID, nil
}

func calculatePoints(receipt repository.Receipt) int {
	totalPoints := 0
	totalPoints += calculateRetailerPoints(receipt.Retailer)

	if isTotalRoundDollar(receipt.Total) {
		totalPoints += 50
	}

	if isTotalMultiple25(receipt.Total) {
		totalPoints += 25
	}

	m := itemsMultiple(len(receipt.Items))
	totalPoints += 5 * m

	totalPoints += itemsPoints(receipt.Items)

	if isDayOdd(receipt.PurchaseDate) {
		totalPoints += 6
	}

	if isTimeInRange(receipt.PurchaseTime) {
		totalPoints += 10
	}

	return totalPoints
}

func calculateRetailerPoints(retailer string) int {
	points := 0
	for _, v := range retailer {
		if unicode.IsLetter(v) || unicode.IsDigit(v) {
			points++
		}
	}
	return points
}

func itemsPoints(items []repository.Item) int {
	points := 0
	for _, item := range items {
		item.ShortDescription = strings.TrimSpace(item.ShortDescription)
		if len(item.ShortDescription)%3 == 0 {
			itemP, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(itemP * 0.2))
		}
	}
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
