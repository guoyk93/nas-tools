package main

import (
	"bytes"
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
	dirSource     = "/volume1/archives"
	dirTargetRoot = "/volume1/tape"

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
	dirTarget := filepath.Join(dirTargetRoot, optTape)
	rg.Must0(os.MkdirAll(dirTarget, 0755))

	// create list file
	func(fileList string) {
		f := rg.Must(os.OpenFile(fileList, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644))
		defer f.Close()

		for _, candidate := range candidates {
			var names []string

			// build names
			var batch []*model.ArchivedFile
			rg.Must0(db.ArchivedFile.Where(
				db.ArchivedFile.Bundle.Eq(candidate.ID),
			).FindInBatches(&batch, 10000, func(tx gen.Dao, b int) (err error) {
				for _, record := range batch {
					if archivestore.ShouldIgnoreFullPath(record.Name) {
						continue
					}
					names = append(names, filepath.Join(record.Year, record.Bundle, record.Name))
				}
				return
			}))

			// sort names
			sort.Strings(names)

			// write names
			for _, name := range names {
				rg.Must(f.Write([]byte(name + "\r\n")))
			}
		}
	}(filepath.Join(dirTarget, "00-INDEX.txt"))

	// create archives
	for _, candidate := range candidates {

		var (
			fileArchive = filepath.Join(dirTarget, candidate.ID+".7z")
			fileIndex   = filepath.Join(dirTarget, candidate.ID+".7z"+".txt")
		)

		_ = os.RemoveAll(fileArchive)
		_ = os.RemoveAll(fileIndex)

		// create 7z archive
		{
			// build args
			args := []string{"7z", "a", "-mx=0"}
			for ex := range archivestore.Ignores {
				args = append(args, "-xr!"+ex)
			}
			for _, ex := range archivestore.IgnorePrefixes {
				args = append(args, "-xr!"+ex+"*")
			}
			args = append(
				args,
				fileArchive,
				filepath.Join(candidate.Year, candidate.ID),
			)

			// run command
			log.Println(strings.Join(args, " "))
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = dirSource
			rg.Must0(cmd.Run())
		}

		// create archive index
		{
			// run command
			buf := &bytes.Buffer{}
			cmd := exec.Command("7z", "l", fileArchive)
			cmd.Stdout = buf
			cmd.Stderr = os.Stderr
			rg.Must0(cmd.Run())

			// save output
			rg.Must0(os.WriteFile(fileIndex, buf.Bytes(), 0644))
		}

		// update bundle set tape
		rg.Must(db.ArchivedBundle.Where(db.ArchivedBundle.ID.Eq(candidate.ID)).UpdateSimple(db.ArchivedBundle.Tape.Value(optTape)))
	}

}
