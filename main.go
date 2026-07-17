package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
)

var version string
var commit string

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

func (opt *Opt) cmd() (int, time.Duration) {
	start := time.Now()
	cmd := exec.Command(opt.Command, opt.Args...)
	cmd.Stdout = os.Stderr
	cmd.Start()
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()
	var status int
	select {
	case <-ctx.Done():
		cmd.Process.Kill()
		if !opt.Quiet {
			log.Printf("Command %s timeout. killed", opt.Command)
		}
		status = TimeoutStatus
	case err := <-done:
		if err != nil {
			if !opt.Quiet {
				log.Printf("Command %s exit with err: %v", opt.Command, err)
			}
		}
		status = cmd.ProcessState.ExitCode()
	}
	duration := time.Since(start)
	if status < 0 {
		status = UnknownCommandStatus
	}
	return status, duration
}

func (opt *Opt) run() int {
	now := time.Now().Unix()
	status, duration := opt.cmd()
	fmt.Printf("command-status.time-taken.%s\t%f\t%d\n", opt.Name, duration.Seconds(), now)
	fmt.Printf("command-status.exit-code.%s\t%d\t%d\n", opt.Name, status, now)
	return 0
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	psr.Usage = "mackerel-plugin-command-status [OPTIONS] -- command args1 args2 ..."
	args, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	if len(args) == 0 {
		psr.WriteHelp(os.Stderr)
		return 1
	}
	opt.Command = args[0]
	if len(args) > 1 {
		opt.Args = args[1:]
	}

	return opt.run()
}
