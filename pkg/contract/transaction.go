package contract

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	TransactionRequest struct {
		Account   uint  `json:"account_id"`
		Operation uint  `json:"operation_id"`
		Amount    int64 `json:"amount"`
	}
)

func (t TransactionRequest) Validate() error {
	return validation.ValidateStruct(
		&t,
		validation.Field(&t.Account, validation.Required),
		validation.Field(&t.Operation, validation.Required),
		validation.Field(&t.Amount, validation.Required),
	)
}
