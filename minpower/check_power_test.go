package minpower

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinPower(t *testing.T) {
	min, ok := new(big.Int).SetString("10995116277760", 10) // 10TiB = 10 * 1024^4
	assert.True(t, ok)

	cases := []struct {
		in   string
		want bool
	}{
		{"f01000", false},
		{"f02620", true},
	}
	for _, c := range cases {
		ok, err := MinQualityPowerOk(context.Background(), c.in, min)
		assert.Equal(t, c.want, ok)
		assert.Nil(t, err)
	}
}
