package main

import (
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	dirSrc = "/volume1/mirrors/Git"
	dirDst = "/volume1/tape/TAPE-CODE"
)

func isGitBareRepository(dir string, entries []os.DirEntry) bool {
	var (
		foundRefs    bool
		foundObjects bool
	)

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == "refs" {
			foundRefs = true
		}
		if entry.IsDir() && entry.Name() == "objects" {
			foundObjects = true
		}
	}

	if !foundRefs || !foundObjects {
		return false
	}

	cmd := exec.Command(
		"git",
		"-C", dir,
		"rev-parse",
		"--is-bare-repository",
	)
	if err := cmd.Run(); err != nil {
		log.Println("found refs and objects but not a git repository:", err.Error())
		return false
	}

	return true
}

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	rg.Must0(iterateDir(dirSrc))
}

func doArchive(dir string) (err error) {
	defer rg.Guard(&err)

	rel, base := filepath.Split(rg.Must(filepath.Rel(dirSrc, dir)))

	rg.Must0(os.MkdirAll(filepath.Join(dirDst, rel), 0755))

	file := filepath.Join(dirDst, rel, base) + ".7z"

	rg.Must0(os.RemoveAll(file))

	log.Println("creating:", file)

	cmd := exec.Command(
		"7z",
		"a",
		"-mx=0",
		file,
		base,
	)
	cmd.Dir = filepath.Join(dirSrc, rel)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	return
}

func iterateDir(dir string) (err error) {
	defer rg.Guard(&err)

	entries := rg.Must(os.ReadDir(dir))

	if isGitBareRepository(dir, entries) {
		return doArchive(dir)
	} else {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if entry.Name() == "@eaDir" {
				continue
			}
			rg.Must0(iterateDir(filepath.Join(dir, entry.Name())))
		}
	}

	return
}
