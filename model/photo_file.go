package model

type PhotoFile struct {
	// shasum256(group + path)
	ID string `gorm:"column:id;primaryKey"`
	// group
	Group string `gorm:"column:group;index;not null"`
	// path
	Path string `gorm:"column:path;type:text"`
	// MD5
	Md5 string `gorm:"column:md5;index;not null"`
}

func (PhotoFile) TableName() string {
	return "photo_files"
}
