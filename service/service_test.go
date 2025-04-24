package service

import (
	"errors"
	"testing"

	"github.com/arjunmalhotra1/fetch-rewards/repository"
	"github.com/stretchr/testify/assert"
)

var test1 = repository.Receipt{
	Retailer:     "M&M Corner Market",
	PurchaseDate: "2022-03-20",
	PurchaseTime: "14:33",
	Items: []repository.Item{
		{
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		}, {
			ShortDescription: "Gatorade",
			Price:            "2.25",
		},
	},
	Total: "9.00",
}

var test2 = repository.Receipt{
	Retailer:     "Target",
	PurchaseDate: "2022-01-01",
	PurchaseTime: "13:01",
	Items: []repository.Item{
		{
			ShortDescription: "Mountain Dew 12PK",
			Price:            "6.49",
		}, {
			ShortDescription: "Emils Cheese Pizza",
			Price:            "12.25",
		}, {
			ShortDescription: "Knorr Creamy Chicken",
			Price:            "1.26",
		}, {
			ShortDescription: "Doritos Nacho Cheese",
			Price:            "3.35",
		}, {
			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
			Price:            "12.00",
		},
	},
	Total: "35.35",
}

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		input          repository.Receipt
		expectedOutput int
	}{
		{
			input:          test1,
			expectedOutput: 109,
		},
		{
			input:          test2,
			expectedOutput: 28,
		},
	}

	for _, test := range tests {
		points := calculatePoints(test.input)
		assert.Equal(t, test.expectedOutput, points, "points are not equal")
	}
}

var test3 = repository.Receipt{
	Retailer: "$1000",
}

func TestProcessReciepts(t *testing.T) {
	testRepo := repository.NewReceiptRepository()
	testService := NewReceiptService(testRepo)

	tests := []struct {
		input  repository.Receipt
		output string
		err    error
	}{{
		input:  test3,
		output: "",
		err:    errors.New(""),
	}}

	for _, test := range tests {
		uuid, err := testService.ProcessReceipts(test.input)
		assert.NotNil(t, err)
		assert.Equal(t, uuid, "")
		assert.ErrorContains(t, err, "Retailer: failed on the 'retailer' rule")
		assert.ErrorContains(t, err, "PurchaseTime: failed on the 'required' rule")

	}

}
