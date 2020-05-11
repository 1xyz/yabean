// Demo client CLI using the go-beanstalkd library
package main

import (
	command "github.com/1xyz/beanstalk-cli/cmd"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	usage := `usage: bs-client [--version] [--addr=<addr>] <command> [<args>...]
options:
   --addr=<addr>  Beanstalkd Address [default: :11300].
   -h, --help
The commands are:
   put        Put a job into a beanstalkd tube.
   reserve    Reserve a job from one or more tubes.
   stats-job  Retrieve statistics for a specific job.
`
	parser := &docopt.Parser{OptionsFirst: true}
	args, err := parser.ParseArgs(usage, nil, "bs-demo-client version 0.1")
	if err != nil {
		log.Errorf("err = %v", err)
		os.Exit(1)
	}

	cmd := args["<command>"].(string)
	cmdArgs := args["<args>"].([]string)

	addr, err := args.String("--addr")
	if err != nil {
		log.Errorf("args.String(--addr). err=%v", err)
		os.Exit(1)
	}

	if err := command.RunCommand(addr, cmd, cmdArgs); err != nil {
		log.Errorf("run command, err = %v", err)
		os.Exit(1)
	}
}
