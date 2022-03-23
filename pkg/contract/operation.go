package contract

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	OperationRequest struct {
		Description string `json:"description"`
		Debit       string `json:"debit"`
	}
)

func (a OperationRequest) Validate() error {
	return validation.ValidateStruct(
		&a,
		validation.Field(&a.Description, validation.Required),
		validation.Field(&a.Debit, validation.Required),
	)
}
