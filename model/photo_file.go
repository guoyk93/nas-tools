package model

type PhotoFile struct {
	// shasum256(bundle + path)
	ID string `gorm:"column:id;primaryKey"`
	// bundle
	Bundle string `gorm:"column:bundle;index;not null"`
	// path
	Path string `gorm:"column:path;type:text"`
	// MD5
	Md5 string `gorm:"column:md5;index;not null"`
}

func (PhotoFile) TableName() string {
	return "photo_files"
}
