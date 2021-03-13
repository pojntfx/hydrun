package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	pflag.Usage = func() {
		fmt.Printf("Run a command for the current directory on multiple processor architectures and operating systems.\n\n")

		fmt.Printf("Usage: %s [options...] <commands...>\n", os.Args[0])

		pflag.PrintDefaults()

		fmt.Printf("\nSee https://github.com/pojntfx/hydrun for more information.\n")
	}

	architecture := pflag.StringP("arch", "a", "amd64,arm64v8", "Processor architecture(s) to run on. Separate multiple values with commas.")
	os := pflag.StringP("os", "o", "debian", "Operating system(s) to run on. Separate multiple values with commas.")

	pflag.Parse()

	log.Println(*architecture, *os)
}
