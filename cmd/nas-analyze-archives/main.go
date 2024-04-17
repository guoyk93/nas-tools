package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/yankeguo/nas-tools/model"
	"github.com/yankeguo/nas-tools/model/dao"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/nas-tools/utils/archivestore"
	"github.com/yankeguo/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	regexpYear               = regexp.MustCompile(`^\d{4}$`)
	regexpYearMonthDayPrefix = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-`)
)

func checkYearEntryDir(fails *[]string, nameYear string, nameBundle string) {
	var err error
	defer utils.Failed(&err, fails)
	defer rg.Guard(&err)

	rec := rg.Must(archivestore.New(db, nameYear, nameBundle))

	doCreate := rg.Must(rec.CountDB()) == 0

	if doCreate {
		if optSkipCreation {
			log.Println("skip creating checksum:", nameBundle)
			return
		}
		log.Println("creating checksum:", nameBundle)
	} else {
		if optSkipValidation {
			log.Println("skip validating checksum:", nameBundle)
			return
		}
		log.Println("validating checksum:", nameBundle)
	}

	if !doCreate {
		rg.Must0(rec.Load())
	}

	if err = checksumDir(
		rec,
		checksumDirOptions{
			doCreate:  doCreate,
			dirBundle: filepath.Join(optDirRoot, nameYear, nameBundle),
		},
		"",
	); err != nil {
		log.Println("failed:", err)
		return
	}

	if doCreate {
		rg.Must0(rec.Save())
	} else {
		if rec.CountNotChecked() != 0 {
			err = errors.New("missing files: " + strings.Join(rec.SampleNotChecked(), ", "))
		}
	}
}

type checksumDirOptions struct {
	doCreate  bool
	dirBundle string
}

func checksumDir(rec *archivestore.Store, opts checksumDirOptions, current string) (err error) {
	doCreate, dirBundle := opts.doCreate, opts.dirBundle

	var entries []os.DirEntry
	if entries, err = os.ReadDir(filepath.Join(dirBundle, current)); err != nil {
		return
	}
	for _, entry := range entries {
		if archivestore.ShouldIgnore(entry.Name()) {
			continue
		}

		if entry.IsDir() {
			if err = checksumDir(
				rec,
				opts,
				filepath.Join(current, entry.Name()),
			); err != nil {
				return
			}
			continue
		}

		if entry.Type() == 0 || entry.Type() == os.ModeSymlink {
			var (
				name = filepath.Join(current, entry.Name())
				sym  = entry.Type() == os.ModeSymlink
			)

			if doCreate {
				err = rec.Add(dirBundle, name, sym)
			} else {
				err = rec.Check(dirBundle, name, sym)
			}
			if err != nil {
				return
			}
		}
	}

	return err
}

func checkYear(fails *[]string, nameYear string) {
	var err error
	defer utils.Failed(&err, fails)
	defer rg.Guard(&err)

	var namesBundle []string

	for _, entryBundle := range rg.Must(os.ReadDir(filepath.Join(optDirRoot, nameYear))) {
		if !entryBundle.IsDir() {
			continue
		}

		nameBundle := entryBundle.Name()

		if archivestore.ShouldIgnore(nameBundle) {
			continue
		}

		if !strings.HasPrefix(nameBundle, nameYear+"-") {
			*fails = append(*fails, "invalid year prefix: "+nameBundle+" in "+nameYear)
			continue
		}
		if !regexpYearMonthDayPrefix.MatchString(nameBundle) {
			*fails = append(*fails, "invalid year-month-day prefix: "+nameBundle+" in "+nameYear)
			continue
		}

		namesBundle = append(namesBundle, nameBundle)
	}

	for _, nameBundle := range namesBundle {
		checkYearEntryDir(fails, nameYear, nameBundle)

		_ = rg.Must(db.ArchivedBundle.Where(
			db.ArchivedBundle.ID.Eq(nameBundle),
			db.ArchivedBundle.Year.Eq(nameYear),
		).FirstOrCreate())
	}
}

const (
	optDirRoot = "/volume1/archives"
)

var (
	optDebug          bool
	optSkipCreation   bool
	optSkipValidation bool
	optFixMissingSize bool
	optFixSymlinkSize bool

	db *dao.Query
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	optDebug, _ = strconv.ParseBool(os.Getenv("OPT_DEBUG"))
	optSkipCreation, _ = strconv.ParseBool(os.Getenv("OPT_SKIP_CREATION"))
	optSkipValidation, _ = strconv.ParseBool(os.Getenv("OPT_SKIP_VALIDATION"))
	optFixMissingSize, _ = strconv.ParseBool(os.Getenv("OPT_FIX_MISSING_SIZE"))
	optFixSymlinkSize, _ = strconv.ParseBool(os.Getenv("OPT_FIX_SYMLINK_SIZE"))

	log.Println("dirRoot:", optDirRoot, "debug:", optDebug, "skip-creation:", optSkipCreation, "skip-validation:", optSkipValidation)

	// create db
	{
		client := rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{}))
		if optDebug {
			client = client.Debug()
		}
		db = dao.Use(client)
	}

	archivestore.Debug = optDebug

	if optFixMissingSize {
		var records []*model.ArchivedFile

		rg.Must0(db.ArchivedFile.Where(
			db.ArchivedFile.Where(
				db.ArchivedFile.Size.Eq(0),
			).Or(
				db.ArchivedFile.Size.Eq(132),
			).Or(
				db.ArchivedFile.CRC32.Eq("00000000"),
			),
			db.ArchivedFile.Symlink.Is(false),
		).FindInBatches(&records, 10000, func(tx gen.Dao, batch int) error {
			return db.Transaction(func(db *dao.Query) (err error) {
				for _, record := range records {
					// read file size
					file := filepath.Join(optDirRoot, record.Year, record.Bundle, record.Name)
					var info os.FileInfo
					if info, err = os.Stat(file); err != nil {
						return
					}
					// update size
					if _, err = db.ArchivedFile.Where(
						db.ArchivedFile.ID.Eq(record.ID),
					).UpdateSimple(
						db.ArchivedFile.Size.Value(info.Size()),
					); err != nil {
						return
					}
				}
				return
			})
		}))
	}

	if optFixSymlinkSize {
		var records []*model.ArchivedFile

		rg.Must0(db.ArchivedFile.
			Select(db.ArchivedFile.ID).
			Where(
				db.ArchivedFile.Symlink.Is(true),
			).
			FindInBatches(&records, 10000, func(tx gen.Dao, batch int) error {
				return db.Transaction(func(db *dao.Query) (err error) {
					for _, record := range records {
						// read link for symlink
						var link string
						if link, err = os.Readlink(
							filepath.Join(optDirRoot, record.Year, record.Bundle, record.Name),
						); err != nil {
							return
						}
						// update size
						if _, err = db.ArchivedFile.Where(
							db.ArchivedFile.ID.Eq(record.ID),
						).UpdateSimple(
							db.ArchivedFile.Size.Value(int64(len(link))),
						); err != nil {
							return err
						}
					}
					return nil
				})
			}),
		)
	}

	var fails []string

	for _, entryYear := range rg.Must(os.ReadDir(optDirRoot)) {
		if !entryYear.IsDir() {
			continue
		}

		nameYear := entryYear.Name()

		if archivestore.ShouldIgnore(nameYear) {
			continue
		}

		switch nameYear {
		case "MAIL", "TAPE", "SUMS", "CRED":
			log.Println("skipping " + nameYear)
		default:
			if !regexpYear.MatchString(nameYear) {
				fails = append(fails, "invalid year name: "+nameYear)
				continue
			}
			log.Println("checking year:", nameYear)
			checkYear(&fails, nameYear)
		}
	}

	if len(fails) > 0 {
		err = errors.New("failures:\n" + strings.Join(fails, "\n"))
	}
}
