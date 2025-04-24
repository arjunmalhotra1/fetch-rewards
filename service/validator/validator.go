package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type ValidationErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func (vErr ValidationErrorResponse) Error() string {
	msg := "Validation Failed:"
	for _, e := range vErr.Errors {
		msg += fmt.Sprintf("[%s: %s]", e.Field, e.Message)
	}

	return msg
}

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("retailer", validateRetailer)
	validate.RegisterValidation("purchaseDate", validateDateFormat)
	validate.RegisterValidation("purchaseTime", validateTimeFormat)
	validate.RegisterValidation("total", validateTotal)
	validate.RegisterValidation("shortDescription", validateShortDescription)
	validate.RegisterValidation("price", validatePrice)

	return validate
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

func BuildValidationErrorResponse(errs validator.ValidationErrors) ValidationErrorResponse {
	errors := []FieldError{}

	for _, err := range errs {
		errors = append(errors, FieldError{
			Field:   err.Field(),
			Rule:    err.Tag(),
			Message: fmt.Sprintf("%s failed on the '%s' rule. The receipt is invalid.", err.Field(), err.Tag()),
		})
	}

	return ValidationErrorResponse{Errors: errors}

}
