package hyperutils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
)

var validate = validator.New()

func BodyParser(c *fiber.Ctx, out any) error {
	if err := c.BodyParser(out); err != nil {
		return err
	} else if err := validate.Struct(out); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func ErrorParser(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("could not found data with your provide condition. %s", err.Error()))
	} else if errors.Is(err, gorm.ErrDuplicatedKey) || errors.Is(err, gorm.ErrPrimaryKeyRequired) {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("your request contains duplicated field. %s", err.Error()))
	} else if errors.Is(err, gorm.ErrInvalidData) || errors.Is(err, gorm.ErrInvalidField) {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("your request contains invalid field. %s", err.Error()))
	} else {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
