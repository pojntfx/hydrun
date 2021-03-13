package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type target struct {
	platform string
	image    string
	command  string
}

func main() {
	// Usage info
	pflag.Usage = func() {
		fmt.Printf("Run a command for the current directory on multiple processor architectures and operating systems.\n\n")

		fmt.Printf("Usage: %s [options...] <commands...>\n", os.Args[0])

		pflag.PrintDefaults()

		fmt.Printf("\nSee https://github.com/pojntfx/hydrun for more information.\n")
	}

	// Parse flags
	arch := pflag.StringP("arch", "a", "amd64,arm64v8", "Processor architecture(s) to run on. Separate multiple values with commas.")
	os := pflag.StringP("os", "o", "debian", "Operating system(s) to run on. Separate multiple values with commas.")
	jobs := pflag.IntP("jobs", "j", 1, "Max amount of arch/os combinations to run in parallel")

	pflag.Parse()

	// Interpret flags
	architectures := strings.Split(*arch, ",")
	operatingSystems := strings.Split(*os, ",")
	command := strings.Join(pflag.Args(), " ")

	// Create target matrix
	targets := []target{}
	for _, architecture := range architectures {
		for _, operatingSystem := range operatingSystems {
			targets = append(targets, target{
				platform: "linux/" + architecture,
				image:    operatingSystem,
				command:  command,
			})
		}
	}

	log.Println(targets, *jobs)
}
