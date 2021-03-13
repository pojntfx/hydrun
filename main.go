package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	osutils "os"
	"os/exec"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/sync/semaphore"
)

// Target is a target in the build matrix
type Target struct {
	Architecture string
	OS           string
	Command      string
}

func main() {
	// Define usage
	pflag.Usage = func() {
		fmt.Printf(`Execute a command for the current directory on multiple architectures and operating systems.

See https://github.com/pojntfx/hydrun for more information.

Usage: %s [OPTION...] "<COMMAND...>"
`, os.Args[0])

		pflag.PrintDefaults()
	}

	// Parse flags
	archFlag := pflag.StringP("arch", "a", "amd64", "Comma-separated list of architectures to run on")
	osFlag := pflag.StringP("os", "o", "debian", "Comma-separated list of operating systems to run on")
	jobFlag := pflag.Int64P("jobs", "j", 1, "Maximum amount of parallel jobs")
	itFlag := pflag.BoolP("it", "i", false, "Attach stdin and setup a TTY")

	pflag.Parse()

	// Validate arguments
	if pflag.NArg() == 0 {
		help := `command needs an argument: 'COMMAND' in "<COMMAND...>"`

		fmt.Println(help)

		pflag.Usage()

		fmt.Println(help)

		os.Exit(2)
	}

	// Interpret arguments
	arches := strings.Split(*archFlag, ",")
	oses := strings.Split(*osFlag, ",")
	command := strings.Join(pflag.Args(), " ")
	pwd := osutils.Getenv("PWD")

	// Create build matrix
	targets := []Target{}
	for _, arch := range arches {
		for _, os := range oses {
			targets = append(targets, Target{
				Architecture: arch,
				OS:           os,
				Command:      command,
			})
		}
	}

	// Setup concurrency
	sem := semaphore.NewWeighted(*jobFlag)
	ctx := context.Background()

	for _, target := range targets {
		// Aquire lock
		if err := sem.Acquire(ctx, 1); err != nil {
			panic(err)
		}

		go func(t Target) {
			// Construct the arguments
			dockerArgs := fmt.Sprintf(`run %v %v:/data:z --platform linux/%v %v /bin/sh -c`, func() string {
				// Attach stdin and setup a TTY
				if *itFlag {
					return "-it -v"
				}

				return "-v"
			}(), pwd, t.Architecture, t.OS)
			commandArgs := fmt.Sprintf(`cd /data && %v`, t.Command)

			// Construct the command
			cmd := exec.Command("docker", append(strings.Split(dockerArgs, " "), commandArgs)...)

			// Handle interactivity
			if *itFlag {
				// Attach stdin, stdout and stderr
				cmd.Stdin = osutils.Stdin
				cmd.Stdout = osutils.Stdout
				cmd.Stderr = osutils.Stderr
			} else {
				// Get stdout and stderr pipes
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					panic(err)
				}
				stderr, err := cmd.StderrPipe()
				if err != nil {
					panic(err)
				}

				// Read from stderr and stdout
				stdoutScanner := bufio.NewScanner(stdout)
				stderrScanner := bufio.NewScanner(stderr)

				// Split into lines
				stdoutScanner.Split(bufio.ScanLines)
				stderrScanner.Split(bufio.ScanLines)

				// Print to stdout with prefix
				prefix := fmt.Sprintf("%v/%v/%v", t.Architecture, t.OS, t.Command)
				go func() {
					for stdoutScanner.Scan() {
						fmt.Println(prefix+"/stdout\t", stdoutScanner.Text())
					}
				}()

				go func() {
					for stderrScanner.Scan() {
						fmt.Println(prefix+"/stderr\t", stderrScanner.Text())
					}
				}()
			}

			// Run the command
			if err := cmd.Run(); err != nil {
				panic(err)
			}

			// Release lock
			sem.Release(1)
		}(target)
	}

	// Wait till all targets have run
	sem.Acquire(ctx, *jobFlag)
}
