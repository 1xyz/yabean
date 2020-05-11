package cmd

import (
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func CmdStatsJob(addr string, argv []string) error {
	usage := `usage: stats-job [--job-id=<id>]
options:
    -h, --help
    --job-id=<id>   Job identifier 

example:
    retrieve statistics for a job with identifier 100
    stats-job --job-id 100`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	jobIDStr, err := opts.String("--job-id")
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

	log.Infof("CmdStatsJob jobId=%v", jobID)
	s, err := c.StatsJob(jobID)
	if err != nil {
		return err
	}
	printMap(s)
	return nil
}

func printMap(s map[string]string) {
	for k, v := range s {
		log.Infof("(%v => %v)", k, v)
	}
}

func CmdStatsTube(addr string, argv []string) error {
	usage := `usage: stats-tube [--tube=<tube>]
options:
    -h, --help
    --tube=<tube>   name of the tube [default: default]

example:
    retrieve statistics for a specific tube with name foobar
    stats-tube --tube foobar`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	tube, err := opts.String("--tube")
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	log.Infof("StatsTube tube=%s", tube)
	t := beanstalk.Tube{Conn: c, Name: tube,}
	s, err := t.Stats()
	if err != nil {
		return err
	}
	printMap(s)
	return nil
}
