package entityvalidator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type EntityValidator interface {
	Register()
}

type entityValidator struct {
}

func NewEntityValidator() EntityValidator {
	return &entityValidator{}
}

var stockValidator validator.Func = func(fl validator.FieldLevel) bool {
	stock, ok := fl.Field().Interface().(int64)
	if ok {
		stockMin := 10 //fixo apenas para testes
		if stock < int64(stockMin) {
			return false
		}
	}
	return true
}

func (ev *entityValidator) Register() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("stockValidator", stockValidator)
	}
}
