package cbr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Debug(t *testing.T) {
	Debug = true
	getRate("CNY", time.Now(), nil)
	assert.True(t, Debug)

	Debug = false
	getRate("CNY", time.Now(), nil)
	assert.False(t, Debug)
}

func Test_getRate_Error(t *testing.T) {
	rate, err := getRate("CNY", time.Now(), nil)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), rate)
}
