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

type ArchivedBundle struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	Year      string    `json:"year" gorm:"column:year;not null;index"`
	CreatedAt time.Time `gorm:"column:created_at;index;not null;autoCreateTime"`
	Tape      string    `json:"tape" gorm:"column:tape;varchar(64);not null;default:'';index"`
}

func (ArchivedBundle) TableName() string {
	return "archived_bundles"
}
