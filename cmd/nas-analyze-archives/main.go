package main

import (
	"errors"
	"github.com/guoyk93/nas-tools/models"
	"github.com/guoyk93/nas-tools/sumstore"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/guoyk93/nas-tools/utils"
	"github.com/guoyk93/rg"
	"gorm.io/driver/mysql"
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

	rec := rg.Must(sumstore.New(client, nameYear, nameBundle))

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

	var ignoreItems []models.ArchivedFileIgnore
	rg.Must0(client.Where(models.ArchivedFileIgnore{
		Year:   nameYear,
		Bundle: nameBundle,
	}).Find(&ignoreItems).Error)

	ignores := map[string]struct{}{}
	for _, item := range ignoreItems {
		ignores[item.Dir] = struct{}{}
	}

	if !doCreate {
		rg.Must0(rec.Load())
	}

	if err = checksumDir(
		rec,
		checksumDirOptions{
			doCreate:  doCreate,
			dirBundle: filepath.Join(optDirRoot, nameYear, nameBundle),
			ignores:   ignores,
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
	ignores   map[string]struct{}
}

func checksumDir(rec *sumstore.Store, opts checksumDirOptions, current string) (err error) {
	doCreate, dirBundle, ignores := opts.doCreate, opts.dirBundle, opts.ignores

	var entries []os.DirEntry
	if entries, err = os.ReadDir(filepath.Join(dirBundle, current)); err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "@eaDir" {
				continue
			}
			if _, ok := ignores[entry.Name()]; ok {
				continue
			}
			if err = checksumDir(
				rec,
				opts,
				filepath.Join(current, entry.Name()),
			); err != nil {
				return
			}
			continue
		}

		if entry.Name() == ".DS_Store" {
			continue
		}

		if strings.HasPrefix(entry.Name(), "._") {
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

		if nameBundle == "@eaDir" {
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

		var b models.ArchivedBundle
		rg.Must0(client.Where(models.ArchivedBundle{ID: nameBundle, Year: nameYear}).FirstOrCreate(&b).Error)

		namesBundle = append(namesBundle, nameBundle)
	}

	for _, nameBundle := range namesBundle {
		checkYearEntryDir(fails, nameYear, nameBundle)
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

	client *gorm.DB
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

	client = rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{}))
	//rg.Must0(client.Debug().AutoMigrate(&sumstore.ArchivedFile{}, &sumstore.ArchivedFileIgnore{}, &sumstore.ArchivedBundle{}))

	if optDebug {
		client = client.Debug()
	}

	sumstore.Debug = optDebug

	if optFixMissingSize {
		var records []models.ArchivedFile
		rg.Must0(client.Where("(size = ? OR size = ? OR crc32 = ?) AND symlink = ?", 0, 132, "00000000", false).FindInBatches(&records, 10000, func(tx *gorm.DB, batch int) error {
			log.Println("fix missing size batch:", batch)
			return tx.Transaction(func(tx *gorm.DB) (err error) {
				for _, record := range records {
					file := filepath.Join(optDirRoot, record.Year, record.Bundle, record.Name)
					var info os.FileInfo
					if info, err = os.Stat(file); err != nil {
						return
					}
					if err = tx.Model(&models.ArchivedFile{}).Where("id = ?", record.ID).Update("size", info.Size()).Error; err != nil {
						return
					}
				}
				return
			})
		}).Error)
	}

	if optFixSymlinkSize {
		var records []models.ArchivedFile
		rg.Must0(client.Select("id").Where("symlink = ?", true).FindInBatches(&records, 10000, func(tx *gorm.DB, batch int) error {
			log.Println("fix symlink size batch:", batch)
			return tx.Transaction(func(tx *gorm.DB) (err error) {
				for _, record := range records {
					var link string
					if link, err = os.Readlink(filepath.Join(optDirRoot, record.Year, record.Bundle, record.Name)); err != nil {
						return
					}
					if err = tx.Model(&models.ArchivedFile{}).Where("id = ?", record.ID).Update("size", len(link)).Error; err != nil {
						return
					}
				}
				return
			})
		}).Error)
	}

	var fails []string

	for _, entryYear := range rg.Must(os.ReadDir(optDirRoot)) {
		if !entryYear.IsDir() {
			continue
		}

		nameYear := entryYear.Name()

		if nameYear == "@eaDir" {
			continue
		}

		switch nameYear {
		case "MAIL", "TAPE", "SUMS":
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
