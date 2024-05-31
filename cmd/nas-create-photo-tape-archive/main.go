package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
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

	fileDst := filepath.Join(dirDst, "TAPE-PHOTO.7z")
	fileIdx := filepath.Join(dirDst, "TAPE-PHOTO.7z.txt")

	{
		args := []string{
			"7z",
			"a",
			"-v100g",
			"-mx=0",
			"-mhe=on",
			"-p" + archivePassword,
			fileDst,
		}

		for _, user := range os.Args[1:] {
			args = append(args, filepath.Join(user, "Photos"))
		}

		// run command
		log.Println(strings.Join(args, " "))
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dirHomes
		rg.Must0(cmd.Run())
	}

	// create archive index
	{
		// run command
		buf := &bytes.Buffer{}
		cmd := exec.Command("7z", "l", "-p"+archivePassword, fileDst+".001")
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr
		rg.Must0(cmd.Run())

		// save output
		rg.Must0(os.WriteFile(fileIdx, buf.Bytes(), 0644))
	}
}
