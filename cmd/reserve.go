package cmd

import (
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func CmdReserve(addr string, argv []string) error {
	usage := `usage: reserve [--timeout=<timeout>] [--tubes=<tubes>] [--no-delete] [--string]
options:
    -h, --help
    --timeout=<timeout>   reservation timeout in seconds [default: 0]
    --tubes=<tubes>       csv of tubes [default: default]
    --no-delete           do not delete (aka. ACK) the job once reserved [default: false]
	--string              display job's body content as a string [default: false]  

example:
    watch for reservations on default tube (topic)
    reserve

    watch for reservations on tubes foo & bar with timeout of 10 seconds
    reserve --timeout 10 --tubes=foo,bar`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	log.Debugf("args:...%v", opts)
	timeout, err := opts.Int("--timeout")
	if err != nil {
		return err
	}

	tubes, err := opts.String("--tubes")
	if err != nil {
		return err
	}

	noDel, err := opts.Bool("--no-delete")
	if err != nil {
		return err
	}

	displayStr, err := opts.Bool("--string")
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	tubeNames := strings.Split(tubes, ",")
	log.Infof("c.reserve() timeout=%v sec tubes=%v no-delete=%v", timeout, tubeNames, noDel)
	ts := beanstalk.NewTubeSet(c, tubeNames...)
	id, body, err := ts.Reserve(time.Duration(timeout) * time.Second)
	if err != nil {
		log.Errorf("reserve. err=%v", err)
		return err
	}

	log.Infof("reserved job id=%v body=%v", id, len(body))
	if displayStr {
		log.Infof("body = %v", string(body))
	}

	if !noDel {
		if err := c.Delete(id); err != nil {
			log.Errorf("delete. err=%v", err)
			return err
		}

		log.Infof("deleted job %v", id)
	}

	return nil
}
