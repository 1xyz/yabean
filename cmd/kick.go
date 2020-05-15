package cmd

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func CmdKick(addr string, argv []string) error {
	usage := `usage: kick <job-id>
options:
    -h, --help

example:
    bury a job with id 1
    bury 1`

	opts, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	id, err := opts.Int("<job-id>")
	if err != nil {
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}

	err = c.KickJob(uint64(id))
	if err != nil {
		log.Errorf("Kick(id=%v), error %v", id, err)
		return err
	}

	fmt.Printf("job with id = %v Kicked.\n", id)
	return nil
}
