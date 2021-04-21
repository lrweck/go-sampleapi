package main

import (
	"encoding/json"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func HandlerCreateAccount(c *fiber.Ctx) error {

	p := PostAccount{}
	err := json.Unmarshal(c.Body(), &p)
	if err != nil {
		unmarshalError(c, err)
		return nil
	}

	if err := DBCreateAccount(p); err != nil {
		dbError(c, err)
		return nil
	}

	c.Status(201)

	return nil
}

func HandlerGetAccount(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("accountId"))

	if err != nil {
		unmarshalError(c, err)
	}

	g := GetAccount{
		AccountID: id,
	}

	docNum, err := DBGetAccount(g.AccountID)

	if err != nil {
		dbError(c, err)
		return err
	}

	g.DocNumber = docNum

	c.JSON(g)

	return nil
}

func HandlerCreateTransaction(c *fiber.Ctx) error {
	/* Insert transaction */

	p := PostTransaction{}
	err := json.Unmarshal(c.Body(), &p)
	if err != nil {
		unmarshalError(c, err)
		return err
	}

	if err := DBCreateTransaction(p); err != nil {
		dbError(c, err)
		return err
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

func dbError(c *fiber.Ctx, err error) {
	c.Status(500)
	c.JSON(&fiber.Map{
		"code":    500,
		"message": err,
	})
}

func unmarshalError(c *fiber.Ctx, err error) {
	c.Status(400)
	c.JSON(&fiber.Map{
		"code":    400,
		"message": err,
	})
}
