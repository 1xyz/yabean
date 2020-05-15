package cmd

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func CmdPeek(addr string, argv []string) error {
	usage := `usage: peek <job-id> [--string]
options:
    -h, --help
    --string    display job's body content as a string [default: false]  

example:
    Peek a job with identifier 1
    peek 1`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	id, err := opts.Int("<job-id>")
	if err != nil {
		return err
	}
	show, err := opts.Bool("--string")
	if err != nil {
		return err
	}

	log.Debugf("c.Peek() id=%d", id)

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	body, err := c.Peek(uint64(id))
	if err != nil {
		log.Errorf("Peek(id=%v), error %v", id, err)
		return err
	}

	fmt.Printf("job with id = %v job size = %v bytes\n", id, len(body))
	if show {
		fmt.Printf("content = [%v]\n", string(body))
	}
	return nil
}

func CmdPeekTube(addr string, argv []string) error {
	usage := `usage: peek-tube <tube> [options]
options:
    -h, --help
    --string    display job's body content as a string [default: false]  

example:
    Peek a job with identifier 1
    peek 1`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	tube, err := opts.String("<tube>")
	if err != nil {
		return err
	}
	show, err := opts.Bool("--string")
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}
	t := beanstalk.Tube{
		Conn: c,
		Name: tube,
	}

	peekTube("ready", show, t.PeekReady)
	peekTube("buried", show, t.PeekBuried)
	peekTube("delayed", show, t.PeekDelayed)
	return nil
}

func peekTube(queueType string, show bool, pf func() (uint64, []byte, error)) {
	id, body, err := pf()
	if err != nil {
		log.Errorf("peek-%v err = %v", queueType, err)
		return
	}
	fmt.Printf("peek=%v, job id = %v size = %v bytes\n", queueType, id, len(body))
	if show {
		fmt.Printf("content = [%v]\n", string(body))
	}
}
