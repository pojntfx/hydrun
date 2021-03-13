package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/sync/semaphore"
)

type target struct {
	architecture    string
	operatingSystem string
	command         string
}

func main() {
	// Usage info
	pflag.Usage = func() {
		fmt.Printf("Run a command for the current directory on multiple processor architectures and operating systems.\n\n")

		fmt.Printf("Usage: %s [options...] \"<commands...>\"\n", os.Args[0])

		pflag.PrintDefaults()

		fmt.Printf("\nSee https://github.com/pojntfx/hydrun for more information.\n")
	}

	// Parse flags
	archFlag := pflag.StringP("arch", "a", "amd64,arm64", "Processor architecture(s) to run on. Separate multiple values with commas.")
	osFlag := pflag.StringP("os", "o", "debian", "Operating system(s) to run on. Separate multiple values with commas.")
	jobsFlag := pflag.Int64P("jobs", "j", 1, "Max amount of arch/os combinations to run in parallel")

	pflag.Parse()

	// Validate the flags
	if pflag.NArg() == 0 {
		pflag.Usage()

		fmt.Println("needs an argument: 'command' in <commands...>")

		return
	}

	// Interpret flags
	architectures := strings.Split(*archFlag, ",")
	operatingSystems := strings.Split(*osFlag, ",")
	command := strings.Join(pflag.Args(), " ")
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Create target matrix
	targets := []target{}
	for _, architecture := range architectures {
		for _, operatingSystem := range operatingSystems {
			targets = append(targets, target{
				architecture:    architecture,
				operatingSystem: operatingSystem,
				command:         command,
			})
		}
	}

	// Run the targets
	sem := semaphore.NewWeighted(*jobsFlag)
	ctx := context.Background()
	for _, t := range targets {
		if err := sem.Acquire(ctx, 1); err != nil {
			panic(err)
		}

		go func(it target) {
			// Construct run command
			dockerArgs := fmt.Sprintf(`run -v %v:/data --platform linux/%v %v /bin/sh -c`, pwd, it.architecture, it.operatingSystem)
			commandArgs := fmt.Sprintf(`cd /data && %v`, it.command)

			cmd := exec.Command("docker", append(strings.Split(dockerArgs, " "), commandArgs)...)

			// Capture stdout and stderr
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				panic(err)
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				panic(err)
			}

			// Start the command
			if err := cmd.Start(); err != nil {
				panic(err)
			}

			// Print stdout and stderr
			stdoutScanner := bufio.NewScanner(stdout)
			stderrScanner := bufio.NewScanner(stderr)

			stdoutScanner.Split(bufio.ScanLines)
			stderrScanner.Split(bufio.ScanLines)

			go func() {
				for stdoutScanner.Scan() {
					fmt.Println(stdoutScanner.Text())
				}
			}()

			go func() {
				for stderrScanner.Scan() {
					fmt.Println(stderrScanner.Text())
				}
			}()

			// Wait till the command has finished
			cmd.Wait()

			sem.Release(1)
		}(t)
	}

	// Wait till all targets have run
	sem.Acquire(ctx, *jobsFlag)
}
