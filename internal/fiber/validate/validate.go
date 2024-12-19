package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

var validate = newValidate()

type RequestInterface any

func CheckParams(ctx *fiber.Ctx, rt RequestInterface) (err error) {
	switch strings.ToUpper(string(ctx.Request().Header.Method())) {
	case fiber.MethodGet:
		err = ctx.QueryParser(rt)
		if err != nil {
			return err
		}
	case fiber.MethodPost, fiber.MethodPut:
		err = ctx.BodyParser(rt)
		if err != nil {
			return err
		}
	}
	err = validate.Struct(rt)
	return
}

func newValidate() *validator.Validate {
	v := validator.New()

	//load custom validate
	v.RegisterValidation("trim", trim)

	return v
}

func trim(fl validator.FieldLevel) bool {
	if fl.Field().CanSet() {
		fl.Field().SetString(strings.TrimSpace(fl.Field().String()))
		return true
	}

	return false
}
