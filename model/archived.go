package model

import "time"

type ArchivedFile struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Year      string    `gorm:"column:year;index;not null"`
	Bundle    string    `gorm:"column:bundle;index;not null"`
	Name      string    `gorm:"column:name;type:text"`
	Symlink   *bool     `gorm:"column:symlink;index;not null;default:false"`
	Size      *int64    `gorm:"column:size;index;not null;default:0"`
	CRC32     string    `gorm:"column:crc32;index;not null"`
	CreatedAt time.Time `gorm:"column:created_at;index;not null;autoCreateTime"`
}

func (ArchivedFile) TableName() string {
	return "archived_files"
}

type ArchivedFileIgnore struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Year      string    `gorm:"column:year;index;not null"`
	Bundle    string    `gorm:"column:bundle;index;not null"`
	Dir       string    `gorm:"column:dir;index;not null"`
	CreatedAt time.Time `gorm:"column:created_at;index;not null;autoCreateTime"`
}

func (ArchivedFileIgnore) TableName() string {
	return "archived_file_ignores"
}

type ArchivedBundle struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	Year      string    `json:"year" gorm:"column:year;not null;index"`
	CreatedAt time.Time `gorm:"column:created_at;index;not null;autoCreateTime"`
}

func (ArchivedBundle) TableName() string {
	return "archived_bundles"
}
