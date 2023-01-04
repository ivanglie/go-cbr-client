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

	"golang.org/x/text/encoding/charmap"
)

const (
	baseURL    = "http://www.cbr.ru/scripts/XML_daily_eng.asp"
	dateFormat = "02/01/2006"
)

// Debug mode
// If this variable is set to true, debug mode activated for the package
var Debug = false

// Cache for requests
var cache map[string]Result
var cacheHits int

// Currency is a currency item
type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  uint   `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nom      uint   `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

// Result is a result representation
type Result struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Currencies []Currency `xml:"Valute"`
}

func getRate(currency string, t time.Time, fetch fetchFunction, useCache bool) (float64, error) {
	if Debug {
		log.Printf("Fetching the currency rate for %s at %v\n", currency, t.Format("02.01.2006"))
	}

	result, err := getCurrenciesCacheOrRequest(t, fetch, useCache)
	if err != nil {
		return 0, err
	}

	for _, v := range result.Currencies {
		if v.CharCode == currency {
			return getCurrencyRateValue(v)
		}
	}
	return 0, fmt.Errorf("Unknown currency: %s", currency)
}

func getCurrencyRateValue(cur Currency) (float64, error) {
	var res float64 = 0
	properFormattedValue := strings.Replace(cur.Value, ",", ".", -1)
	res, err := strconv.ParseFloat(properFormattedValue, 64)
	if err != nil {
		return res, err
	}
	return res / float64(cur.Nom), nil
}

func getCurrenciesCacheOrRequest(t time.Time, fetch fetchFunction, useCache bool) (Result, error) {
	formatedDate := t.Format(dateFormat)

	result := Result{}

	// if currencies were already requested for this date - return from cache, if it is used
	if cachedResult, exist := cache[formatedDate]; exist && useCache {
		log.Printf("Get from cache!")
		result = cachedResult
		cacheHits += 1
	} else {
		err := getCurrencies(&result, t, fetch)
		if err != nil {
			return result, err
		}
		if useCache {
			// if cache is used - put result to cache
			if len(cache) == 0 {
				// if cache is empty - initialize
				cache = make(map[string]Result)
			}
			cache[formatedDate] = result
		}
	}
	return result, nil
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
