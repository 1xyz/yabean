package cmd

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func CmdListTubes(addr string, argv []string) error {
	usage := `usage: list
options:
    -h, --help

example:
    list`

	_, err := docopt.ParseArgs(usage, argv[1:], "version")
	if err != nil {
		log.Errorf("error parsing arguments. err=%v", err)
		return err
	}

	c, err := newConn(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	tubes, err := c.ListTubes()
	if err != nil {
		return err
	}
	for i, t := range tubes {
		fmt.Printf("tube [%d] => %s\n", i+1, t)
	}

	return nil
}
