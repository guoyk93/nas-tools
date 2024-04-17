package archivestore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yankeguo/nas-tools/model"
	"github.com/yankeguo/nas-tools/model/dao"
	"gorm.io/gen"

	"github.com/yankeguo/nas-tools/utils"
)

var (
	Debug = false

	Ignores = map[string]struct{}{
		"@eaDir":       {},
		".DS_Store":    {},
		"venv":         {},
		".venv":        {},
		"node_modules": {},
		"Thumbs.db":    {},
	}
	IgnorePrefixes = []string{
		"._",
	}
)

func ShouldIgnoreFullPath(fullPath string) bool {
	if fullPath == "" || fullPath == "." || fullPath == ".." || fullPath == "/" {
		return false
	}
	if dir, name := filepath.Split(fullPath); ShouldIgnore(name) {
		return true
	} else {
		return ShouldIgnoreFullPath(filepath.Clean(dir))
	}
}

func ShouldIgnore(name string) bool {
	if _, ok := Ignores[name]; ok {
		return true
	}

	for _, prefix := range IgnorePrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

func buildID(year, bundle, name string) string {
	digest := sha256.Sum256([]byte(filepath.Join(year, bundle, name)))
	return hex.EncodeToString(digest[:])
}

func checksumFile(file string, sym bool) (sum string, size int64, err error) {
	h := crc32.NewIEEE()
	if sym {
		var dst string
		if dst, err = os.Readlink(file); err != nil {
			return
		}

		var nsize int
		if nsize, err = io.WriteString(h, dst); err != nil {
			return
		}

		size = int64(nsize)
	} else {
		var f *os.File

		if f, err = os.OpenFile(file, os.O_RDONLY, 0644); err != nil {
			return
		}
		defer f.Close()

		if size, err = io.Copy(h, f); err != nil {
			return
		}
	}
	sum = hex.EncodeToString(h.Sum(nil))
	return
}

type Item struct {
	Name    string
	Symlink bool
	CRC32   string
	Size    int64
}

type Store struct {
	db       *dao.Query
	items    []Item
	data     map[string]Item
	checking map[string]struct{}
	year     string
	bundle   string

	modeWrite *bool
}

func New(db *dao.Query, year, bundle string) (store *Store, err error) {
	store = &Store{
		db:       db,
		data:     map[string]Item{},
		checking: map[string]struct{}{},
		year:     year,
		bundle:   bundle,
	}
	return
}

func (st *Store) mustCheckMode() {
	if st.modeWrite == nil {
		st.modeWrite = utils.Ptr(false)
	} else {
		if *st.modeWrite {
			panic("unexpected write mode")
		}
	}
}

func (st *Store) mustWriteMode() {
	if st.modeWrite == nil {
		st.modeWrite = utils.Ptr(true)
	} else {
		if !*st.modeWrite {
			panic("unexpected check mode")
		}
	}
}

func (st *Store) SampleNotChecked() (out []string) {
	st.mustCheckMode()
	for name := range st.checking {
		if len(out) >= 5 {
			break
		}
		out = append(out, name)
	}
	return
}

func (st *Store) CountDB() (out int64, err error) {
	return st.db.ArchivedFile.Where(
		st.db.ArchivedFile.Year.Eq(st.year),
		st.db.ArchivedFile.Bundle.Eq(st.bundle),
	).Count()
}

func (st *Store) CountChecking() int {
	st.mustCheckMode()
	return len(st.checking)
}

func (st *Store) Add(dirBundle string, name string, sym bool) (err error) {
	st.mustWriteMode()
	var (
		checksum string
		size     int64
	)
	if checksum, size, err = checksumFile(filepath.Join(dirBundle, name), sym); err != nil {
		return
	}

	item := Item{Name: name, Symlink: sym, CRC32: checksum, Size: size}

	if _, existed := st.data[item.Name]; existed {
		err = errors.New("file existed in checksum: " + item.Name)
		return
	}

	if Debug {
		log.Printf("checksum for: %s, sym: %v, checksum: %s", name, sym, checksum)
	}

	st.items = append(st.items, item)
	st.data[item.Name] = item
	st.checking[item.Name] = struct{}{}
	return
}

func (st *Store) Check(dirBundle string, name string, sym bool) (err error) {
	st.mustCheckMode()
	if _, ok := st.data[name]; !ok {
		err = errors.New("file not found in checksum: " + name)
		return
	}

	var (
		checksum string
		size     int64
	)

	if checksum, size, err = checksumFile(filepath.Join(dirBundle, name), sym); err != nil {
		return
	}

	if size != st.data[name].Size {
		err = errors.New("size not match: " + name)
		return
	}

	if checksum != st.data[name].CRC32 {
		err = errors.New("checksum not match: " + name)
		return
	}

	if Debug {
		log.Printf("checksum for: %s, sym: %v, checksum: %s", name, sym, checksum)
	}

	delete(st.checking, name)
	return
}

func (st *Store) Load() (err error) {
	st.mustCheckMode()
	var records []*model.ArchivedFile

	if err = st.db.ArchivedFile.Where(
		st.db.ArchivedFile.Year.Eq(st.year),
		st.db.ArchivedFile.Bundle.Eq(st.bundle),
	).FindInBatches(&records, 1000, func(tx gen.Dao, batch int) error {
		for _, record := range records {
			if ShouldIgnoreFullPath(record.Name) {
				continue
			}
			item := Item{
				Name:    record.Name,
				Symlink: *record.Symlink,
				CRC32:   record.CRC32,
				Size:    *record.Size,
			}
			st.data[record.Name] = item
			st.items = append(st.items, item)
			st.checking[record.Name] = struct{}{}
		}
		return nil
	}); err != nil {
		return
	}

	return
}

func (st *Store) Save() error {
	st.mustWriteMode()
	return st.db.Transaction(func(db *dao.Query) error {
		for _, item := range st.items {
			if err := db.ArchivedFile.Create(&model.ArchivedFile{
				ID:      buildID(st.year, st.bundle, item.Name),
				Year:    st.year,
				Bundle:  st.bundle,
				Name:    item.Name,
				Size:    utils.Ptr(item.Size),
				Symlink: utils.Ptr(item.Symlink),
				CRC32:   item.CRC32,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
