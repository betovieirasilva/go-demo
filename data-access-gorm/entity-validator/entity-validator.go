package entityvalidator

import (
	"example/data-access/entity"

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

var albumStockValidator validator.Func = func(fl validator.FieldLevel) bool {
	stock, ok := fl.Field().Interface().(int64)
	if ok {
		stockMin := 10 //fixo apenas para testes
		if stock < int64(stockMin) {
			return false
		}
	}
	return true
}

func albumValidation(sl validator.StructLevel) {
	album := sl.Current().Interface().(entity.Album)
	if len(album.Artist) == 0 {
		sl.ReportError(album.Artist, "Artist", "", "AlbumValidation", "")
	} else if len(album.Title) == 0 {
		sl.ReportError(album.Title, "Title", "", "AlbumValidation", "")
	} else if album.Price < 1 {
		sl.ReportError(album.Price, "Price", "", "AlbumValidation", "")
	}
}

func (ev *entityValidator) Register() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("stockValidator", albumStockValidator)
		v.RegisterStructValidation(albumValidation, entity.Album{})
	}
}
