package main

import (
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yankeguo/nas-tools/model"
	"github.com/yankeguo/nas-tools/model/dao"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	allowedExtensions = map[string]struct{}{
		".jpg":  {},
		".png":  {},
		".gif":  {},
		".mov":  {},
		".heic": {},
		".mp4":  {},
		".dng":  {},
		".jpeg": {},
		".arw":  {},
		".cr2":  {},
	}

	alertedExts = map[string]struct{}{}
)

var (
	db *dao.Query

	optBundle string
	optDir    string
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	db = dao.Use(rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{})))

	flag.StringVar(&optBundle, "bundle", "", "bundle")
	flag.StringVar(&optDir, "dir", "", "dir")
	flag.Parse()

	if optBundle == "" {
		err = errors.New("bundle is required")
		return
	}

	if optDir == "" {
		err = errors.New("dir is required")
		return
	}

	rg.Must(
		db.PhotoFile.Where(
			db.PhotoFile.Bundle.Eq(optBundle),
		).Delete(),
	)

	rg.Must0(db.Transaction(func(tx *dao.Query) (err error) {
		return checksumDir(optDir, tx)
	}))

	duplicatedMd5s := rg.Must(listDuplicatedMd5s())

	output := []string{"# deletion script"}

	for _, md5 := range duplicatedMd5s {
		records := rg.Must(db.PhotoFile.Where(
			db.PhotoFile.Bundle.Eq(optBundle),
			db.PhotoFile.Md5.Eq(md5),
		).Find())

		if len(records) < 2 {
			err = errors.New("invalid duplicated md5")
			return
		}

		output = append(output, "# ("+md5+") KEEP: "+records[0].Path)

		for _, record := range records[1:] {
			output = append(output, "rm -vf "+strconv.Quote(record.Path))
		}
	}

	log.Println(strings.Join(output, "\n"))
}

func listDuplicatedMd5s() (md5s []string, err error) {
	var rows *sql.Rows

	if rows, err = db.PhotoFile.Where(
		db.PhotoFile.Bundle.Eq(optBundle),
	).Select(
		db.PhotoFile.Md5,
		db.PhotoFile.Md5.Count().As("count"),
	).Group(
		db.PhotoFile.Md5,
	).Having(
		db.PhotoFile.Md5.Count().Gt(1),
	).Rows(); err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			md5   string
			count int
		)
		if err = rows.Scan(&md5, &count); err != nil {
			return
		}
		md5s = append(md5s, md5)
	}

	err = rows.Err()

	return
}

func checksumFile(fullPath string, tx *dao.Query) (err error) {
	id := sha256sum(optBundle + "::" + fullPath)

	var md5 string
	if md5, err = md5sumFile(fullPath); err != nil {
		return
	}

	if err = tx.PhotoFile.Create(&model.PhotoFile{
		ID:     id,
		Bundle: optBundle,
		Path:   fullPath,
		Md5:    md5,
	}); err != nil {
		return
	}

	return
}

func checksumDir(dir string, tx *dao.Query) (err error) {
	log.Println("checking:", dir)

	var entries []fs.DirEntry

	if entries, err = os.ReadDir(dir); err != nil {
		return
	}

	for _, entry := range entries {
		if entry.Name() == "@eaDir" || entry.Name() == ".DS_Store" {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if err = checksumDir(fullPath, tx); err != nil {
				return
			}
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))

		if _, ok := allowedExtensions[ext]; !ok {
			if _, ok := alertedExts[ext]; !ok {
				log.Println("unsupported extension:", ext)
				alertedExts[ext] = struct{}{}
			}
			continue
		}

		if err = checksumFile(fullPath, tx); err != nil {
			return
		}
	}

	return
}

func md5sumFile(fullPath string) (sum string, err error) {
	var f *os.File
	if f, err = os.OpenFile(fullPath, os.O_RDONLY, 0); err != nil {
		return
	}
	defer f.Close()

	hash := md5.New()
	if _, err = io.Copy(hash, f); err != nil {
		return
	}

	sum = hex.EncodeToString(hash.Sum(nil))
	return
}

func sha256sum(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
