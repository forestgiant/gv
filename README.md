## gv
A simple Go vendoring tool. `gv` uses `go get -d` to download packages and moves them into a `vendor` folder for use with [Go 1.5 vendoring](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo/edit). 

## Installation
`go get -u github.com/forestgiant/gv`

## Usage
`gv get [-f] [-fix] [-insecure] [-t] [-u] [packages]` 

## Versioning
`gv` only vendors. For versioning checkout https://github.com/forestgiant/version a fork of [skelterjohn/vendor](https://github.com/skelterjohn/vendor)
