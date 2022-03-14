package contract

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	OperationTypeRequest struct {
		Description string `json:"description"`
		Negative    string `json:"negative" `
	}
)

func (a OperationTypeRequest) Validate() error {
	return validation.ValidateStruct(
		&a,
		validation.Field(&a.Description, validation.Required),
		validation.Field(&a.Negative, validation.Required),
	)
}
