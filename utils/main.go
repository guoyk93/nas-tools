package utils

import (
	"log"
	"os"
)

func Exit(err *error) {
	if *err == nil {
		return
	}
	log.Println("exited with error:", (*err).Error())
	os.Exit(1)
}
