package random

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	type testT struct {
		length int
		base   int
	}
	tests := []testT{
		{
			length: 30,
			base:   8,
		},
		{
			length: 30,
			base:   10,
		},
		{
			length: 30,
			base:   16,
		},
		{
			length: 30,
			base:   36,
		},
		{
			length: 30,
			base:   62,
		},
		{
			length: 30,
			base:   64,
		},
	}

	for _, test := range tests {
		rand.Seed(0)
		s, e := RandomGenerator.String(test.length, test.base)
		assert.Nil(t, e)
		t.Log(s)
	}
}

func TestSecureGenerator(t *testing.T) {
	type testT struct {
		length int
		base   int
	}
	tests := []testT{
		{
			length: 30,
			base:   8,
		},
		{
			length: 30,
			base:   10,
		},
		{
			length: 30,
			base:   16,
		},
		{
			length: 30,
			base:   36,
		},
		{
			length: 30,
			base:   62,
		},
		{
			length: 30,
			base:   64,
		},
	}

	for _, test := range tests {
		rand.Seed(0)
		s, e := SecureRandomGenerator.String(test.length, test.base)
		assert.Nil(t, e)
		t.Log(s)
	}
}

func TestBits(t *testing.T) {
	tests := map[uint64]uint64{
		0:         0,
		1:         1,
		2:         2,
		3:         2,
		4:         3,
		7:         3,
		8:         4,
		15:        4,
		16:        5,
		1<<63 - 1: 63,
		1 << 63:   64,
		1<<64 - 1: 64,
	}

	for i, b := range tests {
		assert.Equal(t, b, bits(i))
	}
}
