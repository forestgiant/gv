package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const GO_PATH_ENV_NAME = "GOPATH"
const GO_15_VENDOR_EXPERIMENT = "GO15VENDOREXPERIMENT"

func main() {
	if os.Getenv(GO_15_VENDOR_EXPERIMENT) != "1" {
		fmt.Println("The gv command expects the", GO_15_VENDOR_EXPERIMENT, "environment variable to be set to", 1)
		os.Exit(0)
	}

	var args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("The gv command expects the format of 'go get'.")
		os.Exit(0)
	} else {
		if args[0] != "get" {
			fmt.Println("The only command currently supported is 'get'.")
			os.Exit(0)
		}
		fmt.Println(args)
	}

	//Get the PWD
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//Set the GOPATH to PWD
	fmt.Println("Temporarily overriding GOPATH to", path)
	os.Setenv(GO_PATH_ENV_NAME, path)

	//Issue 'go get' command
	fmt.Println("Running go with commands=", args)
	goGetCommand := exec.Command("go", args...)
	goGetCommand.Stdin = os.Stdin
	goGetCommand.Stdout = os.Stdout
	goGetCommand.Stderr = os.Stderr
	err = goGetCommand.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//Making vendor folder
	fmt.Println("Making vendor folder")
	vendorPath := filepath.Join(path, "vendor")
	err = os.Mkdir(vendorPath, 0700)
	if err != nil {
		fmt.Println(err)
	}

	srcPath := filepath.Join(path, "src")

	err = os.Rename(srcPath, vendorPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Removing src folder (created by go get command)")
	err = os.Remove(srcPath)
	if err == nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
