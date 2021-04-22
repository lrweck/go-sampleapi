package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

// Handler pra criar uma account.
// Chama DBCreateAccount pra armazenar as informações
func HandlerCreateAccount(c *fiber.Ctx) error {

	p := PostAccount{}

	if err := c.BodyParser(&p); err != nil {
		return fiberError(c, fiber.StatusUnprocessableEntity, err)
	}

	if p.DocNumber == "" {
		return fiberError(c, fiber.StatusBadRequest, "O campo 'document_number' é obrigatório")
	}

	id, err := DBCreateAccount(p)
	if err != nil {
		return fiberError(c, fiber.StatusInternalServerError, err)
	}

	// Status Created + Location pelo Header para fácil identificação
	c.Status(201)
	c.Location(c.Context().URI().String() + "/" + id.String())

	return nil
}

// Handler pra buscar uma account.
// Chama DBGetAccount pra retornar as informações
func HandlerGetAccount(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("accountId"))

	if err != nil {
		return fiberError(c, fiber.StatusUnprocessableEntity, err)
	}

	g := GetAccount{
		AccountID: id,
	}

	docNum, err := DBGetAccount(g.AccountID)

	// Verifica se encontramos alguma coisa
	if err == pgx.ErrNoRows {
		return fiberError(c, fiber.StatusNotFound, "id não encontrado")
	}

	// Precisamos garantir que DocNum seja obrigatorio/não-vazio
	g.DocNumber = docNum

	return c.JSON(g)

}

// Handler pra criar uma transação.
// Chama DBCreateTransaction pra armazenar as informações
func HandlerCreateTransaction(c *fiber.Ctx) error {
	/* Insert transaction */

	p := PostTransaction{}
	if err := c.BodyParser(&p); err != nil {
		return fiberError(c, fiber.StatusUnprocessableEntity, err)
	}

	// Compara com um UUID em branco/vazio
	if p.AccountID == (uuid.UUID{}) {
		return fiberError(c, fiber.StatusBadRequest, "O campo 'account_id' é obrigatório")
	}

	if p.Amount <= 0 {
		return fiberError(c, fiber.StatusBadRequest, "O campo 'amount' deve ser maior que zero")
	}

	if p.OpeType <= 0 || p.OpeType > 4 {
		return fiberError(c, fiber.StatusBadRequest, "O campo 'operation_type_id' deve estar entre 1 e 4")
	}

	// Consultar os erros individuais

	err := DBCreateTransaction(p)

	if err != nil {
		// Significa que não encontrou a account
		if err == pgx.ErrNoRows {
			return fiberError(c, fiber.StatusNotFound, "'account_id' inexistente")
		}

		return fiberError(c, fiber.StatusInternalServerError, err)
	}

	c.Status(201)

	return nil
}

type PostAccount struct {
	DocNumber string `json:"document_number"`
}

type GetAccount struct {
	AccountID uuid.UUID `json:"account_id"`
	PostAccount
}

type PostTransaction struct {
	AccountID uuid.UUID `json:"account_id"`
	OpeType   int       `json:"operation_type_id"`
	Amount    float64   `json:"amount"`
}

// Helper pra facilitar o envio dos erros e diminuir duplicação
func fiberError(c *fiber.Ctx, errCode int, errMessage interface{}) error {
	c.Status(errCode)
	return c.JSON(&fiber.Map{
		"code":    errCode,
		"message": errMessage,
	})
}
