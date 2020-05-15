package cmd

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
)

func CmdStatsJob(addr string, argv []string) error {
	usage := `usage: stats-job <job-id>
options:
    -h, --help

example:
    Retrieve statistics for a job with identifier 100
    stats-job 100`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	jobIDStr, err := opts.String("<job-id>")
	if err != nil {
		return err
	}
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	log.Debugf("CmdStatsJob jobId=%v", jobID)
	s, err := c.StatsJob(jobID)
	if err != nil {
		return err
	}
	printMap(s)
	return nil
}

func printMap(s map[string]string) {
	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("(%v => %v)\n", k, s[k])
	}
}

func CmdStatsTube(addr string, argv []string) error {
	usage := `usage: stats-tube <tube>
options:
    -h, --help

example:
    retrieve statistics for the tube foobar
    stats-tube  foobar`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	tube, err := opts.String("<tube>")
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	log.Infof("StatsTube tube=%s", tube)
	t := beanstalk.Tube{Conn: c, Name: tube}
	s, err := t.Stats()
	if err != nil {
		return err
	}
	printMap(s)
	return nil
}

func CmdStats(addr string, argv []string) error {
	usage := `usage: stats-tube [--tube=<tube>]
options:
    -h, --help

example:
    retrieve statistics for the beanstalk service
    stats`

	_, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	s, err := c.Stats()
	if err != nil {
		return err
	}
	printMap(s)
	return nil
}
