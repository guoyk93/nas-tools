package main

import (
	"errors"
	"flag"
	"github.com/yankeguo/nas-tools/model"
	"github.com/yankeguo/nas-tools/model/dao"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/nas-tools/utils/archivestore"
	"github.com/yankeguo/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	dirArchives = "/volume1/archives"
	dirTapeRoot = "/volume1/tape"

	sizeThreshold = 1400 * 1000 * 1000 * 1000
)

var (
	optDebug bool

	db *dao.Query
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	var (
		optTape string
	)
	flag.StringVar(&optTape, "tape", "", "tape name")
	flag.Parse()

	if optTape == "" {
		panic("tape name is required")
		return
	}

	// create db
	{
		client := rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{}))
		if optDebug {
			client = client.Debug()
		}
		db = dao.Use(client)
	}

	var (
		candidates []*model.ArchivedBundle
	)

	// select candidates
	{
		bundles := rg.Must(db.ArchivedBundle.Where(
			db.ArchivedBundle.Tape.Eq(""),
		).Order(db.ArchivedBundle.ID.Asc()).Find())

		var totalSize int64
		for _, bundle := range bundles {
			var record *model.ArchivedFile

			record, err = db.ArchivedFile.Where(
				db.ArchivedFile.Bundle.Eq(bundle.ID),
			).Select(
				db.ArchivedFile.Bundle,
				db.ArchivedFile.Size.Sum().As("size"),
			).Group(db.ArchivedFile.Bundle).Take()

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = nil
					continue
				} else {
					rg.Must0(err)
				}
			}

			if *record.Size == 0 {
				rg.Must0(errors.New("bundle size is 0"))
			}

			totalSize += *record.Size

			if totalSize > sizeThreshold {
				break
			}

			candidates = append(candidates, bundle)
		}
	}

	// create workspace
	dirTape := filepath.Join(dirTapeRoot, optTape)
	rg.Must0(os.MkdirAll(dirTape, 0755))

	// create list
	func(fileList string) {
		f := rg.Must(os.OpenFile(fileList, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644))
		defer f.Close()

		for _, bundle := range candidates {
			var names []string
			var batch []*model.ArchivedFile

			rg.Must0(db.ArchivedFile.Where(
				db.ArchivedFile.Bundle.Eq(bundle.ID),
			).FindInBatches(&batch, 10000, func(tx gen.Dao, b int) (err error) {
				for _, record := range batch {
					names = append(names, filepath.Join(record.Year, record.Bundle, record.Name))
				}
				return
			}))

			sort.Strings(names)

			for _, name := range names {
				rg.Must(f.Write([]byte(name + "\r\n")))
			}
		}
	}(filepath.Join(dirTape, "LIST.txt"))

	// create tar
	{
		args := []string{
			"tar",
			"--record-size", "1m",
			"-cvf", "archive.tar",
			"--owner", "yanke:1000",
			"--group", "yanke:1000",
			"LIST.txt",
		}
		log.Println(strings.Join(args, " "))
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dirTape
		rg.Must0(cmd.Run())
	}

	// append tape
	{
		args := []string{
			"tar",
			"--record-size", "1m",
			"-rvf", "archive.tar",
			"--owner", "yanke:1000",
			"--group", "yanke:1000",
		}
		for name := range archivestore.Ignores {
			args = append(args, "--exclude", name)
		}
		for _, name := range archivestore.IgnorePrefixes {
			args = append(args, "--exclude", name+"*")
		}
		args = append(args, "-C", dirArchives)
		for _, candidate := range candidates {
			args = append(args, filepath.Join(candidate.Year, candidate.ID))
		}
		log.Println(strings.Join(args, " "))
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dirTape
		rg.Must0(cmd.Run())
	}
}
