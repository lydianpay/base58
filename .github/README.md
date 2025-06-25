# base58
---
<div align="center">

[![Go Report Card](https://goreportcard.com/badge/lydianpay/base58g)](https://goreportcard.com/report/lydianpay/base58)

</div>

Written in Go ('Golang' for search engines) with zero external dependencies, this package implements encoding for 
Base58. Base58 encoding is used for Bitcoin (BTC) addresses, along with other cryptocurrencies such as Tron. Similar to
Base64 encoding, Base58 removes characters humans often confuse (0, I, O, l).

Please note: at this time, decoding a Base58 encoded string is not included in this package.

## Installation & Usage

1. Once confirming you have [Go](https://go.dev/doc/install) installed, the command below will add
   `base58` as a dependency to your Go program.
```shell
go get -u github.com/lydianpay/base58
```
2. Import the package into your code
```go
package main

import (
    "github.com/lydianpay/base58"
)
```
3. Pass a string you would like to encode
```go
encoded := base58.Encode("Your String")
```
