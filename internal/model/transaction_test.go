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
