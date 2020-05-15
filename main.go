// Demo client CLI using the go-beanstalkd library
package main

import (
	command "github.com/1xyz/yabean/cmd"
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
	usage := `usage: yabean [--version] [--addr=<addr>] <command> [<args>...]
options:
   --addr=<addr>  Beanstalkd Address [default: :11300].
   -h, --help
The commands are:
   del        Delete a specific job.
   kick       Kick a buried job (Note: see reserve command to bury a job).
   list       List tubes.
   peek       Peek at a specific job.
   peek-tube  Peek into a specific tube.
   put        Put a job into a beanstalkd tube.
   reserve    Reserve a job from one or more tubes.
   stats      Retrieve serve statistics.	
   stats-job  Retrieve statistics for a specific job.
   stats-tube Retrieve statistics for a specific tube.
`
	parser := &docopt.Parser{OptionsFirst: true}
	args, err := parser.ParseArgs(usage, nil, "beanstalk cli version 0.1")
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
