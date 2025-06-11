package usecase

import (
	"github.com/stretchr/testify/assert"
	"money-stat/internal/model"
	"testing"
)

func TestUpdateSummByAccountType(t *testing.T) {
	var statDto AccountStatDto

	currencyUsdAccount := model.Instrument{
		Rate: 10,
	}

	rubleAccount := model.Account{
		Instrument: 2,
		Balance:    10000,
	}

	usdAccount := model.Account{
		Instrument: 1,
		Balance:    1000,
		Currency:   currencyUsdAccount,
	}

	statDto.updateSummByAccountType(rubleAccount)
	statDto.updateSummByAccountType(usdAccount)

	assert.Equal(t, statDto.SummRuble, rubleAccount.Balance)
	assert.Equal(t, statDto.SummDollar, usdAccount.Balance)
	assert.Equal(t, statDto.RateDollar, currencyUsdAccount.Rate)

}
