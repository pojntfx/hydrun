package main

import (
	"flag"
	"log"
)

func main() {
	architecture := flag.String("arch", "amd64", "Processor architecture to run")
	os := flag.String("os", "debian:10", "Operating system to run")
	cmd := flag.String("cmd", "make", "Command to run")

	flag.Parse()

	log.Println(*architecture, *os, *cmd)
}
