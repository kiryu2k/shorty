package validator

import "github.com/go-playground/validator/v10"

type RequestValidator struct {
	v *validator.Validate
}

func NewValidator() *RequestValidator {
	return &RequestValidator{validator.New()}
}

func (v *RequestValidator) Validate(instance any) error {
	return v.v.Struct(instance)
}
