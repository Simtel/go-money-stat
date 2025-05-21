package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction_GetTagsTitle(t *testing.T) {

	tag1 := Tag{Title: "tag1"}
	tag2 := Tag{Title: "tag2"}
	transaction := Transaction{
		Tag: []Tag{tag1, tag2},
	}

	assert.Equal(t, "tag1 tag2 ", transaction.GetTagsTitle())
}

func TestTransaction_IsTransfer(t *testing.T) {
	transaction := Transaction{
		Income:  1000,
		Outcome: 1000,
	}
	assert.True(t, transaction.IsTransfer())
	assert.False(t, transaction.IsIncome())
	assert.False(t, transaction.IsOutcome())
}

func TestTransaction_IsIncome(t *testing.T) {
	transaction := Transaction{
		Income: 1000,
	}
	assert.True(t, transaction.IsIncome())
	assert.False(t, transaction.IsTransfer())
	assert.False(t, transaction.IsOutcome())
}

func TestTransaction_IsOutcome(t *testing.T) {
	transaction := Transaction{
		Outcome: 1000,
	}
	assert.True(t, transaction.IsOutcome())
	assert.False(t, transaction.IsTransfer())
	assert.False(t, transaction.IsIncome())
}
