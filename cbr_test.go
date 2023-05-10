package cbr

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct{}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 500}, nil
}

type MockHttpClientErr struct{}

func (m *MockHttpClientErr) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("error")
}

func TestClient_GetRate(t *testing.T) {
	Debug = true

	client := NewClient()
	rate, err := client.GetRate("USD", time.Now())
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, rate, float64(1))
}

func TestClient_GetRate_Error(t *testing.T) {
	Debug = false

	// unknown currency: _
	client := NewClient()
	rate, err := client.GetRate("_", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "unknown currency: _", err.Error())
	assert.Equal(t, float64(0), rate)

	// status code: 500
	client.httpClient = &MockHttpClient{}
	rate, err = client.GetRate("CNY", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "status code: 500", err.Error())
	assert.Equal(t, float64(0), rate)

	// error
	client.httpClient = &MockHttpClientErr{}
	rate, err = client.GetRate("CNY", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
	assert.Equal(t, float64(0), rate)
}

func Test_currencyRateValue_Error(t *testing.T) {
	c := Currency{}
	c.Value = "0'1"
	rate, err := currencyRateValue(c)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), rate)
}
