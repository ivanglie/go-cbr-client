package cbr

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/charmap"
)

const (
	baseURL    = "http://www.cbr.ru/scripts/XML_daily_eng.asp"
	dateFormat = "02/01/2006"
)

// Debug mode
// If this variable is set to true, debug mode activated for the package
var Debug = false

// Currency is a currency item
type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  uint   `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nom      uint   `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

// Returns properly formatted currency Value string
func (cur Currency) ValueString() string {
	return strings.Replace(cur.Value, ",", ".", -1)
}

// Returns currency Value in float64, without nominal correction
func (cur Currency) ValueFloatRaw() (float64, error) {
	return strconv.ParseFloat(cur.ValueString(), 64)
}

// Returns currency Value in float64, corrected by nominal
func (cur Currency) ValueFloat() (float64, error) {
	res, err := cur.ValueFloatRaw()
	if err != nil {
		return res, err
	}
	return res / float64(cur.Nom), nil
}

// Returns currency Value in Decimal, without nominal correction
//
// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
func (cur Currency) ValueDecimalRaw() (decimal.Decimal, error) {
	return decimal.NewFromString(cur.ValueString())
}

// Returns currency Value in Decimal, corrected by nominal
//
// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
func (cur Currency) ValueDecimal() (decimal.Decimal, error) {
	res, err := cur.ValueDecimalRaw()
	if err != nil {
		return res, err
	}
	nominal := decimal.NewFromInt(int64(cur.Nom))
	return res.Div(nominal), nil
}

// Result is a result representation
type Result struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Currencies []Currency `xml:"Valute"`
}

func getRate(currency string, t time.Time, fetch fetchFunction) (float64, error) {
	if Debug {
		log.Printf("Fetching the currency rate for %s at %v\n", currency, t.Format("02.01.2006"))
	}
	curr, err := getCurrency(currency, t, fetch)
	if err != nil {
		return 0, err
	}
	return curr.ValueFloat()
}

func getRateDecimal(currency string, t time.Time, fetch fetchFunction) (decimal.Decimal, error) {
	if Debug {
		log.Printf("Fetching the currency rate for %s at %v\n  in Decimal", currency, t.Format("02.01.2006"))
	}
	curr, err := getCurrency(currency, t, fetch)
	if err != nil {
		return decimal.Zero, err
	}
	return curr.ValueDecimal()
}

func getRateString(currency string, t time.Time, fetch fetchFunction) (string, error) {
	if Debug {
		log.Printf("Fetching the currency rate string for %s at %v\n", currency, t.Format("02.01.2006"))
	}
	curr, err := getCurrency(currency, t, fetch)
	if err != nil {
		return "", err
	}
	return curr.ValueString(), nil
}

func getCurrency(currency string, t time.Time, fetch fetchFunction) (Currency, error) {
	var result Result
	err := getCurrencies(&result, t, fetch)
	if err != nil {
		return Currency{}, err
	}
	for _, v := range result.Currencies {
		if v.CharCode == currency {
			return v, nil
		}
	}
	return Currency{}, fmt.Errorf("unknown currency: %s", currency)
}

func getCurrencies(v *Result, t time.Time, fetch fetchFunction) error {
	url := baseURL + "?date_req=" + t.Format(dateFormat)
	if fetch == nil {
		return errors.New("fetch is empty")
	}
	resp, err := fetch(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("Unknown charset: %s", charset)
		}
	}
	err = decoder.Decode(&v)
	if err != nil {
		return err
	}

	return nil
}
