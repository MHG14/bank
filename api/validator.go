package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/mhg14/simplebank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// check if currency is supported
		return util.IsSupportedCurrency(currency)
	}
	return false
}
