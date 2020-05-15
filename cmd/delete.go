package cmd

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func CmdDelete(addr string, argv []string) error {
	usage := `usage: del <job-id>
options:
    -h, --help

example:
    delete a job with id 1
    del-job 1`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	id, err := opts.Int("<job-id>")
	if err != nil {
		return err
	}

	log.Debugf("c.Delete() id=%d", id)

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	err = c.Delete(uint64(id))
	if err != nil {
		log.Errorf("Del(id=%v), error %v", id, err)
		return err
	}

	fmt.Printf("job with id = %v deleted.\n", id)
	return nil
}
