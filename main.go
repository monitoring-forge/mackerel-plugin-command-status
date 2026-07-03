package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
)

// Version by Makefile
var version string

type opts struct {
	OptArgs    []string
	OptCommand string
	Timeout    time.Duration `long:"timeout" default:"30s" description:"Timeout to wait for command finished"`
	Name       string        `short:"n" long:"name" description:"Metrics name" required:"true"`
	Quiet      bool          `short:"q" long:"quiet" description:"Suppress error output of sub command"`
	Version    bool          `short:"v" long:"version" description:"Show version"`
}

func runCmd(opts opts) int {
	now := time.Now().Unix()
	start := time.Now()
	cmd := exec.Command(opts.OptCommand, opts.OptArgs...)
	cmd.Stdout = os.Stderr
	cmd.Start()
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()
	var status int
	select {
	case <-ctx.Done():
		cmd.Process.Kill()
		if !opts.Quiet {
			log.Printf("Command %s timeout. killed", opts.OptCommand)
		}
		status = 137
	case err := <-done:
		if err != nil {
			if !opts.Quiet {
				log.Printf("Command %s exit with err: %v", opts.OptCommand, err)
			}
		}
		status = cmd.ProcessState.ExitCode()
	}
	duration := time.Since(start)
	fmt.Printf("command-status.time-taken.%s\t%f\t%d\n", opts.Name, duration.Seconds(), now)
	fmt.Printf("command-status.exit-code.%s\t%d\t%d\n", opts.Name, status, now)
	return 0
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := opts{}
	psr := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	psr.Usage = "mackerel-plugin-command-status [OPTIONS] -- command args1 args2 ..."
	args, err := psr.Parse()
	if opts.Version {
		fmt.Fprintf(os.Stderr, "Version: %s\nCompiler: %s %s\n",
			version,
			runtime.Compiler,
			runtime.Version())
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if len(args) == 0 {
		psr.WriteHelp(os.Stderr)
		os.Exit(1)
	}
	opts.OptCommand = args[0]
	if len(args) > 1 {
		opts.OptArgs = args[1:]
	}

	return runCmd(opts)
}
