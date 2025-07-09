package zenmoney

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstrument_IsDollar(t *testing.T) {
	inst := Instrument{
		ShortTitle: "USD",
	}

	assert.True(t, inst.IsDollar())
}
