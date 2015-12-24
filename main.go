package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
)

const GO_PATH_ENV_NAME = "GOPATH"
const GO_15_VENDOR_EXPERIMENT = "GO15VENDOREXPERIMENT"

func main() {
	fmt.Println("Vendoring packages...")

	// Check that we received the expected format
	var args = os.Args[1:]
	if os.Getenv(GO_15_VENDOR_EXPERIMENT) != "1" {
		fmt.Println("The gv command expects the", GO_15_VENDOR_EXPERIMENT, "environment variable to be set to", 1)
		os.Exit(1)
	} else if len(args) == 0 {
		fmt.Println("The gv command expects the format of 'go get'.")
		os.Exit(1)
	} else if args[0] != "get" {
		fmt.Println("The only command currently supported is 'get'.")
		os.Exit(1)
	}

	// Insert -d flag after go get. This instructs get to stop after downloading the package
	args = append(args[:1], append([]string{"-d"}, args[1:]...)...)

	// Set PATH to the current working directory
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Setenv(GO_PATH_ENV_NAME, path)

	// Set up some path vars
	vendorPath := filepath.Join(path, "vendor")
	srcPath := filepath.Join(path, "src")

	// Instantiate our 'go get' command
	goGetCommand := exec.Command("go", args...)
	goGetCommand.Stdin = os.Stdin
	goGetCommand.Stdout = os.Stdout
	goGetCommand.Stderr = os.Stderr

	// Establish our exit channels
	exit := make(chan bool)
	success := make(chan bool)

	// Run the primary routine
	go func() {
		// Run the 'go get' command and rename src to vendor
		if err = goGetCommand.Run(); err == nil {
			if err = mergeVendors(srcPath, vendorPath); err == nil {
				success <- true
				return
			}
		}

		//Clean up if there was an error
		fmt.Println(err)
		fmt.Println("Cleaning up...")
		if err = os.RemoveAll(srcPath); err != nil {
			fmt.Println(err)
		}

		exit <- true
	}()

	// Listen for interrupts, and if received, cancel 'go get'
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Println("\nCancelling...")
		if err = goGetCommand.Process.Kill(); err != nil {
			fmt.Println("Unable to kill running process:", err)
			exit <- true
		}
	}()

	// Listen for exits
	select {
	case <-success:
		fmt.Println("Done.")
		os.Exit(0)
	case <-exit:
		fmt.Println("Unable to vendor packages.")
		os.Exit(1)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func mergeVendors(src string, dst string) error {
	// Loop through all orgs in the src directory
	domains, _ := ioutil.ReadDir(src)
	for _, domain := range domains {
		if domain.IsDir() {
			srcDomainPath := filepath.Join(src, domain.Name())
			dstDomainPath := filepath.Join(dst, domain.Name())

			// Ensure that the dst domain folder exists
			if !fileExists(dstDomainPath) {
				if err := os.MkdirAll(dstDomainPath, 0700); err != nil {
					return err
				}
			}

			// Loop through the orgs for this domain
			orgs, _ := ioutil.ReadDir(srcDomainPath)
			for _, org := range orgs {
				srcOrgPath := filepath.Join(srcDomainPath, org.Name())
				dstOrgPath := filepath.Join(dstDomainPath, org.Name())

				// Ensure that the dst domain folder exists
				if !fileExists(dstOrgPath) {
					if err := os.MkdirAll(dstOrgPath, 0700); err != nil {
						return err
					}
				}

				// Loop through the repos for this org
				repos, _ := ioutil.ReadDir(srcOrgPath)
				for _, repo := range repos {
					srcRepoPath := filepath.Join(srcOrgPath, repo.Name())
					dstRepoPath := filepath.Join(dstOrgPath, repo.Name())

					// Overwrite (Remove) any content that exists at dstRepoPath
					if fileExists(dstRepoPath) {
						if err := os.RemoveAll(dstRepoPath); err != nil {
							return err
						}
					}

					// Copy each repo to corresponding dst directory
					if err := os.Rename(srcRepoPath, dstRepoPath); err != nil {
						return err
					}

					vendorPath := filepath.Join(domain.Name(), org.Name(), repo.Name())
					fmt.Println("Vendored: ", vendorPath)
				}
			}
		}
	}

	// Remove the src directory
	if err := os.RemoveAll(src); err != nil {
		return err
	}

	return nil
}
