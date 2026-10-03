package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/monitoring-forge/flagrun"
)

var version string

const (
	TimeoutStatus        = 137
	UnknownCommandStatus = 127
)

type Opt struct {
	Args    []string
	Command string
	Timeout time.Duration `long:"timeout" default:"30s" description:"Timeout to wait for command finished"`
	Name    string        `short:"n" long:"name" description:"Metrics name" required:"true"`
	Quiet   bool          `short:"q" long:"quiet" description:"Suppress error output of sub command"`
	Version bool          `short:"v" long:"version" description:"Show version"`
}

func (opt *Opt) QuietLogf(format string, v ...any) {
	if !opt.Quiet {
		log.Printf(format, v...)
	}
}

func (opt *Opt) cmd() (int, time.Duration) {
	start := time.Now()
	cmd := exec.Command(opt.Command, opt.Args...)
	cmd.Stdout = os.Stderr
	if err := cmd.Start(); err != nil {
		opt.QuietLogf("Command %s start failed: %v", opt.Command, err)
		return UnknownCommandStatus, time.Since(start)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()
	var status int
	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		opt.QuietLogf("Command %s timeout. killed", opt.Command)
		status = TimeoutStatus
	case err := <-done:
		if err != nil {
			opt.QuietLogf("Command %s exit with err: %v", opt.Command, err)
		}
		status = cmd.ProcessState.ExitCode()
	}
	duration := time.Since(start)
	if status < 0 {
		status = UnknownCommandStatus
	}
	return status, duration
}

func (opt *Opt) Run(_ []string) (any, int) {
	now := time.Now().Unix()
	status, duration := opt.cmd()
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "command-status.time-taken.%s\t%f\t%d\n", opt.Name, duration.Seconds(), now)
	fmt.Fprintf(buf, "command-status.exit-code.%s\t%d\t%d\n", opt.Name, status, now)
	return buf.String(), flagrun.OK
}

func (opt *Opt) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}
	opt.Command = args[0]
	if len(args) > 1 {
		opt.Args = args[1:]
	}
	return nil
}

func main() {
	opt := &Opt{}
	os.Exit(flagrun.Go(opt, flagrun.Version(version), flagrun.ArgsRequired(), flagrun.Validator(opt.Validate)))
}
