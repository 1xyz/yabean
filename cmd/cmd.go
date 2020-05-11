package cmd

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
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
