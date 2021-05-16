package http

import (
	"fmt"
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	js "github.com/lrweck/go-sampleapi/serializer/json"
	"github.com/lrweck/go-sampleapi/service"
	"github.com/pkg/errors"
)

// RedirectHandler - Interface for endpoint handlers
type ApiHandler interface {
	GetAccount(*fiber.Ctx) error
	PostAccount(*fiber.Ctx) error
	PostTransaction(*fiber.Ctx) error
}

type handler struct {
	apiService *service.ApiService
}

func NewHandler(apiServ *service.ApiService) ApiHandler {
	return &handler{
		apiService: apiServ,
	}
}

func (h *handler) GetAccount(ct *fiber.Ctx) error {
	id, err := uuid.Parse(ct.Params("accountId"))
	contentType := ct.Get("Content-Type", "application/json")

	if err != nil {
		return err
	}

	acc, err := h.apiService.FindAccount(id)

	if err != nil {
		if errors.Cause(err) == service.ErrAccountNotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}

	respBody, err := h.accSerializer(contentType).Encode(acc)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return setupResponse(ct, contentType, fiber.StatusOK, respBody)

}

func (h *handler) PostAccount(ct *fiber.Ctx) error {
	contentType := ct.Get("Content-Type", "application/json")
	account, err := h.accSerializer(contentType).Decode(ct.Body())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	err = h.apiService.StoreAccount(account)
	if err != nil {
		if errors.Cause(err) == service.ErrAccountInvalid {
			return fiber.ErrBadRequest
		}
		return fiber.ErrInternalServerError
	}

	respBody, err := h.accSerializer(contentType).Encode(account)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return setupResponse(ct, contentType, fiber.StatusCreated, respBody)
}

func (h *handler) PostTransaction(ct *fiber.Ctx) error {
	contentType := ct.Get("Content-Type")
	tx, err := h.txSerializer(contentType).Decode(ct.Body())
	if err != nil {
		return fiber.ErrInternalServerError
	}

	err = h.apiService.StoreTransaction(tx)
	fmt.Printf("4 - %s", err)
	if err != nil {
		if errors.Cause(err) == service.ErrTransactionInvalid {
			return fiber.ErrBadRequest
		}
		return fiber.ErrInternalServerError
	}

	respBody, err := h.txSerializer(contentType).Encode(tx)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return setupResponse(ct, contentType, fiber.StatusCreated, respBody)
}

// Sends a statuscode and response to the user
func setupResponse(f *fiber.Ctx, contentType string, statusCode int, body []byte) error {
	f.Append("Content-Type", contentType)
	f.Status(statusCode)

	if err := f.Send(body); err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return nil
}

// Serializes queries according to the content-type header
func (h *handler) accSerializer(contentType string) service.AccountSerializer {
	switch contentType {
	case "application/json":
		return &js.Account{}
	case "application/xml":
		return nil
	default:
		return nil
	}
}

// Serializes queries according to the content-type header
func (h *handler) txSerializer(contentType string) service.TransactionSerializer {
	switch contentType {
	case "application/json":
		return &js.Transactions{}
	case "application/xml":
		return nil
	default:
		return nil
	}
}
