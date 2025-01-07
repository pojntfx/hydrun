package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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
	archFlag := pflag.StringP("arch", "a", runtime.GOARCH, "Comma-separated list of architectures to run on")
	osFlag := pflag.StringP("os", "o", "debian", "Comma-separated list of operating systems (Docker images) to run on")
	jobFlag := pflag.Int64P("jobs", "j", 1, "Maximum amount of parallel jobs")
	itFlag := pflag.BoolP("it", "i", false, "Attach stdin and setup a TTY")
	contextFlag := pflag.StringP("context", "c", "", "Directory to use in the container (default is the current working directory)")
	extraArgs := pflag.StringP("extra-args", "e", "", "Extra arguments to pass to the Docker command")
	pullFlag := pflag.BoolP("pull", "p", false, "Always pull the specified tags of the operating systems (Docker images)")
	quietFlag := pflag.BoolP("quiet", "q", false, "Disable logging executed commands")
	mountFlag := pflag.BoolP("mount", "m", true, "Enable mounting the directory specified with the context flag")
	readOnlyFlag := pflag.BoolP("readyOnly", "r", false, "Mount the directory specified as read-only")

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
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("could not get working directory:", err)
	}
	if *contextFlag != "" {
		pwd = *contextFlag
	}

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

	// Pull and tag images
	for _, target := range targets {
		existsCmd := exec.Command("docker", strings.Split(fmt.Sprintf(`inspect %v`, getImageNameWithSuffix(target.OS, target.Architecture)), " ")...)

		if !*quietFlag {
			log.Println(existsCmd)
		}

		shouldPullAndTag := true
		if *pullFlag {
			shouldPullAndTag = false
		} else {
			if output, err := existsCmd.CombinedOutput(); err != nil {
				// Titlecase-insensitive comparison since `docker` and `podman-docker` differ in their casing here
				if strings.Contains(strings.ToLower(string(output)), strings.ToLower("Error: No such object")) {
					shouldPullAndTag = false
				} else {
					log.Fatalln("could not check if image already exists:", err)
				}
			}
		}

		if !shouldPullAndTag {
			runCmd := exec.Command("docker", strings.Split(fmt.Sprintf(`pull --platform linux/%v %v`, target.Architecture, target.OS), " ")...)

			if !*quietFlag {
				log.Println(runCmd)
			}

			if output, err := runCmd.CombinedOutput(); err != nil {
				log.Fatalln("could not pull image:", err.Error()+":", string(output))
			}

			tagCmd := exec.Command("docker", strings.Split(fmt.Sprintf(`tag %v %v`, target.OS, getImageNameWithSuffix(target.OS, target.Architecture)), " ")...)

			if !*quietFlag {
				log.Println(tagCmd)
			}

			if output, err := tagCmd.CombinedOutput(); err != nil {
				log.Fatalln("could not tag image:", err.Error()+":", string(output))
			}
		}
	}

	// Setup concurrency
	sem := semaphore.NewWeighted(*jobFlag)
	ctx := context.Background()

	for _, target := range targets {
		// Aquire lock
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatalln("could not acquire lock:", err)
		}

		go func(t Target) {
			defer sem.Release(1)

			// Construct the arguments
			dockerArgs := fmt.Sprintf(
				`run %v%v--platform linux/%v%v %v /bin/sh -c`,
				func() string {
					// Attach stdin and setup a TTY
					if *itFlag {
						return "-it "
					}

					return ""
				}(),
				func() string {
					if *mountFlag {
						if *readOnlyFlag {
							return fmt.Sprintf("-v %v:/data:ro ", pwd)
						}

						return fmt.Sprintf("-v %v:/data:z ", pwd)
					}

					return ""
				}(),
				t.Architecture,
				func() string {
					args := *extraArgs
					if args != "" {
						args = " " + args
					}

					return args
				}(),
				getImageNameWithSuffix(t.OS, t.Architecture),
			)
			commandArgs := t.Command
			if *mountFlag {
				commandArgs = fmt.Sprintf(`cd /data && %v`, t.Command)
			}

			// Construct the command
			cmd := exec.Command("docker", append(strings.Split(dockerArgs, " "), commandArgs)...)

			if !*quietFlag {
				log.Println(cmd)
			}

			// Attach stdout and stderr
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Attach stdin
			if *itFlag {
				cmd.Stdin = os.Stdin
			}

			// Run the command
			if err := cmd.Run(); err != nil {
				log.Fatalln("could not run command:", err)
			}
		}(target)
	}

	// Wait till all targets have run
	sem.Acquire(ctx, *jobFlag)
}

func getImageNameWithSuffix(image, architecture string) string {
	return image + "-" + strings.Replace(architecture, "/", "-", -1)
}
