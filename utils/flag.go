package utils

import (
	"errors"
	"flag"
)

func SoloFlag(name string) (output string, err error) {
	flag.StringVar(&output, name, "", "target "+name)
	flag.Parse()

	if output == "" {
		err = errors.New("option '-" + name + "' is required")
		return
	}
	return
}
