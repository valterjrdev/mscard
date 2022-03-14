package contract

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	AccountRequest struct {
		Document string `json:"document_number"`
	}
)

func (a AccountRequest) Validate() error {
	return validation.ValidateStruct(
		&a,
		validation.Field(&a.Document, validation.Required),
	)
}
