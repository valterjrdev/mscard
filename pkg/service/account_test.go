package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"ms/card/pkg/persistence/entity"
	"ms/card/pkg/persistence/repository"
	"testing"
	"time"
)

func TestAccount_UpdateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		input    *entity.Account
		negative bool
		amount   int64
		expected int64
	}{
		{
			input: &entity.Account{
				Limit: 2000,
			},
			negative: true,
			amount:   int64(100),
			expected: int64(1900),
		},
		{
			input: &entity.Account{
				Limit: 2000,
			},
			negative: false,
			amount:   int64(100),
			expected: int64(2100),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run("", func(t *testing.T) {
			mockAccountRepository := repository.NewMockAccounts(ctrl)
			mockAccountRepository.EXPECT().UpdateLimit(gomock.Any(), tt.input).Return(nil)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			accountService := NewAccount(AccountOpts{
				AccountRepository: mockAccountRepository,
			})
			err := accountService.UpdateLimit(ctx, tt.input, tt.amount, tt.negative)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, tt.input.Limit)
		})
	}
}

func TestAccount_UpdateLimit_Exceeded_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockAccountEntity := &entity.Account{
		Limit: 50,
	}
	accountService := NewAccount(AccountOpts{})
	err := accountService.UpdateLimit(ctx, mockAccountEntity, 100, true)
	assert.EqualError(t, err, ErrLimitExceeded.Error())
}
