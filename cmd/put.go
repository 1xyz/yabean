package cmd

import (
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func CmdPut(addr string, argv []string) error {
	usage := `usage: put [--body=<body>] [--pri=<pri>] [--ttr=<ttr>] [--delay=<delay>] [--tube=<tube>]
options:
    -h, --help
    --body=<body>     body [default: hello]
    --pri=<pri>       job priority [default: 1]
    --ttr=<ttr>       ttr in seconds [default: 10]
    --delay=<delay>   job delay in seconds [default: 0]
    --tube=<tube>     tube (topic) to put the job [default: default]

example:
    put --body "hello world"
    put --body "hello world" --tube foo`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	tube, err := opts.String("--tube")
	if err != nil {
		return err
	}

	log.Debugf("args:...%v", opts)
	body, err := opts.String("--body")
	if err != nil {
		return err
	}

	pri, err := opts.Int("--pri")
	if err != nil {
		return err
	}

	ttr, err := opts.Int("--ttr")
	if err != nil {
		return err
	}

	delay, err := opts.Int("--delay")
	if err != nil {
		return err
	}

	log.Infof("c.Put() body=%v, pri=%v, delay=%v sec, ttr=%v sec tube=%v",
		body, pri, ttr, delay, tube)

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	var t *beanstalk.Tube = nil
	if tube != "default" {
		t = &beanstalk.Tube{Conn: c, Name: tube}
	}

	var id uint64
	if t == nil {
		// t == nil; indicates no specific tube is used the put call is made to the default tube (implicitly)
		id, err = c.Put([]byte(body), uint32(pri), time.Duration(delay)*time.Second, time.Duration(ttr)*time.Second)
	} else {
		id, err = t.Put([]byte(body), uint32(pri), time.Duration(delay)*time.Second, time.Duration(ttr)*time.Second)
	}

	if err != nil {
		log.Errorf("Put(...), error %v", err)
		return err
	}

	log.Infof("c.Put() returned id = %v", id)
	return nil
}
