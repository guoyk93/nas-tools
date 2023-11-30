package archivestore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/guoyk93/nas-tools/model"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/guoyk93/nas-tools/utils"
	"gorm.io/gorm"
)

var (
	Debug = false
)

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
	items      []Item
	data       map[string]Item
	notChecked map[string]struct{}
	client     *gorm.DB
	year       string
	bundle     string
	ignores    []string
}

func New(client *gorm.DB, year, bundle string) (store *Store, err error) {
	store = &Store{
		data:       map[string]Item{},
		notChecked: map[string]struct{}{},
		client:     client,
		year:       year,
		bundle:     bundle,
	}
	var items []model.ArchivedFileIgnore
	if err = client.Where(&model.ArchivedFileIgnore{Year: year, Bundle: bundle}).Find(&items).Error; err != nil {
		return
	}
	for _, item := range items {
		store.ignores = append(store.ignores, item.Dir)
	}
	return
}

func (st *Store) SampleNotChecked() (out []string) {
	for name := range st.notChecked {
		if len(out) >= 5 {
			break
		}
		out = append(out, name)
	}
	return
}

func (st *Store) CountDB() (out int64, err error) {
	err = st.client.Where(&model.ArchivedFile{
		Year:   st.year,
		Bundle: st.bundle,
	}).Model(&model.ArchivedFile{}).Count(&out).Error
	return
}

func (st *Store) CountNotChecked() int {
	return len(st.notChecked)
}

func (st *Store) Add(dirBundle string, name string, sym bool) (err error) {
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
	st.notChecked[item.Name] = struct{}{}
	return
}

func (st *Store) Check(dirBundle string, name string, sym bool) (err error) {
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

	delete(st.notChecked, name)

	for i, item := range st.items {
		if item.Name == name {
			st.items = append(st.items[:i], st.items[i+1:]...)
			break
		}
	}
	return
}

func (st *Store) Load() (err error) {
	var records []model.ArchivedFile

	if err = st.client.Where(model.ArchivedFile{
		Year:   st.year,
		Bundle: st.bundle,
	}).FindInBatches(
		&records,
		1000,
		func(tx *gorm.DB, batch int) error {
			for _, record := range records {
				item := Item{
					Name:    record.Name,
					Symlink: *record.Symlink,
					CRC32:   record.CRC32,
					Size:    *record.Size,
				}
				st.data[record.Name] = item
				st.notChecked[record.Name] = struct{}{}
				st.items = append(st.items, item)
			}
			return nil
		},
	).Error; err != nil {
		return err
	}

	return
}

func (st *Store) Save() error {
	return st.client.Transaction(func(tx *gorm.DB) error {
		for _, item := range st.items {
			if err := tx.Create(&model.ArchivedFile{
				ID:      buildID(st.year, st.bundle, item.Name),
				Year:    st.year,
				Bundle:  st.bundle,
				Name:    item.Name,
				Size:    utils.Ptr(item.Size),
				Symlink: utils.Ptr(item.Symlink),
				CRC32:   item.CRC32,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
