package zenmoney

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_IsDollar(t *testing.T) {
	account := Account{
		Instrument: 1,
	}
	assert.True(t, account.IsDollar())
	assert.False(t, account.IsRuble())
}

func TestAccount_IsRuble(t *testing.T) {
	account := Account{
		Instrument: 2,
	}
	assert.True(t, account.IsRuble())
	assert.False(t, account.IsDollar())
}
