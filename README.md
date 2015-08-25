## gv
A simple Go vendoring tool. `gv` uses `go get -d` to download packages and moves them into a `vendor` folder for use with [Go 1.5 vendoring](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo/edit). 

## Requirements
In order to use this tool, you must have Go 1.5 installed and enable vendoring in your `.bash_profile` or `.bashrc`:
```
export GO15VENDOREXPERIMENT=1
```

## Installation
`go get -u github.com/forestgiant/gv`

## Usage
Must be used from project root directory.

`gv get [-f] [-fix] [-insecure] [-t] [-u] [packages]` 

## Go get wrapper - (Bash only)
We've added a bash script to wrap `go get` and add a `-v` flag. If you pass this flag to go get it will automatically use `gv` for vendoring

Ex. `go get -v github.com/forestgiant/gv`

To enable this you must edit your `.bash_profile` or `.bashrc` and add:

```
# gv: Used to wrap go get to add -v flag 
export GOCOMMANDLOCATION=/usr/local/bin
source $GOPATH/src/github.com/forestgiant/gv/go_to_gv
```

## Versioning
`gv` only vendors. For versioning checkout https://github.com/forestgiant/version a fork of [skelterjohn/vendor](https://github.com/skelterjohn/vendor)
