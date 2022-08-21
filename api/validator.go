package api

import (
	"github.com/emon46/bank-application/util"
	"github.com/go-playground/validator/v10"
)

var validateCurrency validator.Func = func(level validator.FieldLevel) bool {
	if currency, ok := level.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
