package cmd

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func newConn(addr string) (*beanstalk.Conn, error) {
	c, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		log.Errorf("error dial beanstalkd %v", err)
		return nil, err
	}

	return c, nil
}

type cmdFunc func(string, []string) error
type cmdFuncMap map[string]cmdFunc

var cmdFuncs = cmdFuncMap{
	"del":        CmdDelete,
	"kick":       CmdKick,
	"list":       CmdListTubes,
	"peek":       CmdPeek,
	"peek-tube":  CmdPeekTube,
	"put":        CmdPut,
	"reserve":    CmdReserve,
	"stats":      CmdStats,
	"stats-job":  CmdStatsJob,
	"stats-tube": CmdStatsTube,
}

func RunCommand(addr string, cmd string, args []string) (err error) {
	f, ok := cmdFuncs[cmd]
	if !ok {
		return fmt.Errorf("%s is not a valid command", cmd)
	}

	argv := append([]string{cmd}, args...)
	return f(addr, argv)
}

func getJobID(usage string, argv []string) (uint64, error) {
	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return 0, err
	}

	id, err := opts.Int("<job-id>")
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}
