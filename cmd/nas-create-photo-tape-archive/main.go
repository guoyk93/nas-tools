package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
)

const (
	dirDst   = "/volume1/tape/TAPE-PHOTO"
	dirHomes = "/volume1/homes"
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	archivePassword := strings.TrimSpace(os.Getenv("ARCHIVE_PASSWORD"))
	if archivePassword == "" {
		err = errors.New("ARCHIVE_PASSWORD is required")
		return
	}

	rg.Must0(os.MkdirAll(dirDst, 0755))

	args := []string{
		"7z",
		"a",
		"-v100g",
		"-mx=0",
		"-mhe=on",
		"-p" + archivePassword,
		filepath.Join(dirDst, "TAPE-PHOTO.7z"),
	}

	for _, user := range os.Args[1:] {
		args = append(args, filepath.Join(user, "Photos"))
	}

	log.Println(strings.Join(args, " "))
}
