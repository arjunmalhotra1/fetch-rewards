package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/arjunmalhotra1/fetch-rewards/repository"
	"github.com/arjunmalhotra1/fetch-rewards/service"
	"github.com/arjunmalhotra1/fetch-rewards/service/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	rs service.ReceiptService
}

type successResponsePostReceipt struct {
	Id string `json:"id"`
}

type successResponseGetPoints struct {
	Points int `json:"points"`
}

func NewHandler(service service.ReceiptService) Handler {
	return Handler{
		rs: service,
	}
}

func (h Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"status": "ok"})
}

func (h Handler) GetPoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	points, ok := h.rs.GetPoints(id)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"description": "No receipt found for that ID."})
		return
	}
	render.JSON(w, r, successResponseGetPoints{Points: points})
}

func (h Handler) ProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var receipt repository.Receipt
	receiptBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading the request body", err)
	}
	json.Unmarshal(receiptBody, &receipt)

	log.Println("received receipt request for: ", string(receiptBody))

	id, err := h.rs.ProcessReceipts(receipt)
	if err != nil {
		if vErr, ok := err.(validator.ValidationErrorResponse); ok {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, vErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, "Internal Server error")
		return
	}
	render.JSON(w, r, successResponsePostReceipt{Id: id})
}
