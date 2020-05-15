package cmd

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func CmdReserve(addr string, argv []string) error {
	usage := `usage: reserve [--del|--bury|--release] [options]
options:
    -h, --help
    --timeout=<seconds>   reservation timeout in seconds [default: 0]
    --tubes=<tubes>       csv of tubes [default: default]
    --string              display job's body content as a string [default: false]
  Post reserve actions:
    --bury                bury the job once a job is reserved
    --del                 delete the job (similar to ACK) once a job is reserved
    --release             release the job (similar to NACK) once a job is reserved
  Post reserve action options:
    --pri=<int>           new priority if the job is buried or released [default: 1024]
    --delay=<seconds>     new delay if the job is release [default: 10]
  Other reserve options:
    --touch=<int>         touch (aka renew TTR) the reserved job n times prior to either burying, 
                          deleting, releasing or timeout [default: 0]

example:
    watch for reservations on default tube (topic)
    reserve

    watch for reservations on tubes foo & bar with timeout of 10 seconds
    reserve --timeout 10 --tubes=foo,bar

    delete the job after it is reserved frpm the default tube 
    reserve --del

    bury the job with a priority 123 after it is reserved from the foo tube
    reserve --tubes=foo --bury --priority 123

    release the job immediately after it is reserved from the bar tube
    reserve --tube=bar --release

    touch a job 5 times and bury it after it is reserved from the foobar tube
    reserve --tube=foobar --touch 5 --bury`

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
	displayStr, err := opts.Bool("--string")
	if err != nil {
		return err
	}
	del, err := opts.Bool("--del")
	if err != nil {
		return err
	}
	bury, err := opts.Bool("--bury")
	if err != nil {
		return err
	}
	release, err := opts.Bool("--release")
	if err != nil {
		return err
	}
	pri, err := opts.Int("--pri")
	if err != nil {
		return err
	}
	delay, err := opts.Int("--delay")
	if err != nil {
		return err
	}
	touch, err := opts.Int("--touch")
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	tubeNames := strings.Split(tubes, ",")
	ts := beanstalk.NewTubeSet(c, tubeNames...)
	id, body, err := ts.Reserve(time.Duration(timeout) * time.Second)
	if err != nil {
		log.Errorf("reserve. err=%v", err)
		return err
	}

	fmt.Printf("reserved job id=%v body=%v\n", id, len(body))
	if displayStr {
		fmt.Printf("body = %v\n", string(body))
	}

	for i := 0; i < touch; i++ {
		log.Infof("Try to touch a job %v (%d)/(%d)", id, i+1, touch)
		stats, err := c.StatsJob(id)
		if err != nil {
			log.Errorf("cannot get job statistics err = %v")
			return err
		}
		log.Infof("job state = %v", stats["state"])
		ttr, err := strconv.Atoi(stats["ttr"])
		if err != nil {
			return err
		}

		if _, _, err := ts.Reserve(time.Duration(ttr) * time.Second); err != nil {
			if err.Error() != "reserve-with-timeout: deadline soon" {
				log.Errorf("reserve: err=%v. expected deadline-soon", err)
				return err
			}
			log.Infof("DeadlineSoon error returned")
		} else {
			log.Error("expected to get an error")
			return fmt.Errorf("reserve: expected a deadline-soon error")
		}

		if err := c.Touch(id); err != nil {
			log.Errorf("reserve: c.Touch err=%v", err)
			return err
		}
		log.Infof("Successfully touched job %v (%d)/(%d)", id, i+1, touch)
	}

	if del {
		if err := c.Delete(id); err != nil {
			log.Errorf("delete. err=%v", err)
			return err
		}
		fmt.Printf("deleted job %v\n", id)
		return nil
	}
	if bury {
		if err := c.Bury(id, uint32(pri)); err != nil {
			log.Errorf("bury err = %v", err)
			return err
		}
		fmt.Printf("buried job %v, pri = %v\n", id, pri)
		return nil
	}
	if release {
		d := time.Duration(delay) * time.Second
		if err := c.Release(id, uint32(pri), d); err != nil {
			log.Errorf("release err = %v", err)
			return err
		}
		fmt.Printf("release job %v, pri = %v, delay = %v\n", id, pri, d)
		return nil
	}
	log.Infof("job allowed to timeout without delete, bury or release actions")
	return nil
}
