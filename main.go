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
	fmt.Println("Vendoring packages...")

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
	}

	//Get the PWD
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Set GOPATH to the current working directory
	os.Setenv(GO_PATH_ENV_NAME, path)

	// insert -d flag after go get. This instructs get to stop after downloading the package
	args = append(args[:1], append([]string{"-d"}, args[1:]...)...)

	// Run go get
	goGetCommand := exec.Command("go", args...)
	goGetCommand.Stdin = os.Stdin
	goGetCommand.Stdout = os.Stdout
	goGetCommand.Stderr = os.Stderr
	err = goGetCommand.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Create the vendor directory
	vendorPath := filepath.Join(path, "vendor")
	srcPath := filepath.Join(path, "src")
	err = os.Rename(srcPath, vendorPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done.")
}
