# Golang client for the Central Bank of the Russian Federation currency rates API

[![Go Reference](https://pkg.go.dev/badge/github.com/ivanglie/go-cbr-client.svg)](https://pkg.go.dev/github.com/ivanglie/go-cbr-client)
[![Test](https://github.com/ivanglie/go-cbr-client/actions/workflows/test.yml/badge.svg)](https://github.com/ivanglie/go-cbr-client/actions/workflows/test.yml)
[![Codecov](https://codecov.io/gh/ivanglie/go-cbr-client/branch/master/graph/badge.svg?token=46HUJQAM56)](https://codecov.io/gh/ivanglie/go-cbr-client)

go-cbr-client is a fork of [matperez's](https://github.com/matperez) [client](https://github.com/matperez/go-cbr-client) for [CBRF API](http://www.cbr.ru/development/).

## Example

First, ensure the library is installed and up to date by running ```go get -u github.com/ivanglie/go-cbr-client```.

This is a very simple app that just displays exchange rate of US dollar.

```golang
package main

import (
	"fmt"
	"time"

	cbr "github.com/ivanglie/go-cbr-client"
)

func main() {
	client := cbr.NewClient()
	
	// For float64 value:
	rateFloat64, err := client.GetRate("USD", time.Now())

	// For Decimal value:
	rateDecimal, err := client.GetRateDecimal("USD", time.Now())

	// For String value with dot as decimal separator:
	rateString, err := client.GetRateString("USD", time.Now())

	if err != nil {
		panic(err)
	}
	fmt.Println(rateFloat64)
}
```
See [main.go](./_example/main.go).

## References

For more information check out the following links:

* [CBRF API](http://www.cbr.ru/development/SXML/)
* [CBRF technical resources](http://www.cbr.ru/eng/development/) (EN)
