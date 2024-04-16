package main

import (
	"bytes"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	tapeSizeThreshold = 1500 * 1024 * 1024 * 1024 // 1.5TB

	dirTapeOrig = "/volume1/archives"
	dirTapePack = "/volume1/tape/PACK"
	dirTapeMark = "/volume1/tape/MARK"

	extMark = ".mark"
)

var (
	patternYear   = regexp.MustCompile(`^\d{4}$`)
	patternBundle = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}.+$`)

	excludes = []string{
		"node_modules",
		".DS_Store",
		".venv",
		"venv",
		"@eaDir",
		"Thumbs.db",
	}
)

func isPackedSizeExceeded() bool {
	var existedSize int64

	for _, entry := range rg.Must(os.ReadDir(dirTapePack)) {
		if entry.IsDir() {
			continue
		}
		info := rg.Must(entry.Info())
		existedSize += info.Size()
	}

	return existedSize > tapeSizeThreshold
}

func isMarked(bundle string) bool {
	_, err := os.Stat(filepath.Join(dirTapeMark, bundle+extMark))
	return err == nil
}

func markBundle(bundle string) {
	rg.Must0(os.WriteFile(filepath.Join(dirTapeMark, bundle+extMark), nil, 0644))
}

func packBundle(year, bundle string) {
	fileTarget := filepath.Join(dirTapePack, bundle+".tar.zst")

	{
		args := []string{
			"tar", "--zstd", "-cvf", fileTarget,
			"--owner", "yanke:1000", "--group", "yanke:1000",
		}
		for _, ex := range excludes {
			args = append(args, "--exclude", ex)
		}
		args = append(args, "-C", dirTapeOrig, filepath.Join(year, bundle))

		log.Println(strings.Join(args, " "))

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		rg.Must0(cmd.Run())
	}

	{
		args := []string{
			"tar", "-tvf", fileTarget,
		}

		log.Println(strings.Join(args, " "))

		buf := &bytes.Buffer{}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = buf
		rg.Must0(cmd.Run())

		rg.Must0(os.WriteFile(fileTarget+".txt", buf.Bytes(), 0644))
	}
}

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	log.Println("calculating existing tape packs size...")

	if isPackedSizeExceeded() {
		log.Println("existing tape packs size exceeds threshold, skip preparing tape pack")
		return
	}

	for _, entryYear := range rg.Must(os.ReadDir(dirTapeOrig)) {
		if !entryYear.IsDir() {
			continue
		}
		if !patternYear.MatchString(entryYear.Name()) {
			continue
		}
		for _, entryBundle := range rg.Must(os.ReadDir(filepath.Join(dirTapeOrig, entryYear.Name()))) {
			if !entryBundle.IsDir() {
				continue
			}
			if !patternBundle.MatchString(entryBundle.Name()) {
				continue
			}
			if isMarked(entryBundle.Name()) {
				continue
			}
			packBundle(entryYear.Name(), entryBundle.Name())
			//markBundle(entryBundle.Name())
			if isPackedSizeExceeded() {
				log.Println("existing tape packs size exceeds threshold, skip preparing tape pack")
				return
			}
		}
	}
}