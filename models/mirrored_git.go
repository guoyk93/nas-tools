package models

import "time"

type MirroredGitRepo struct {
	Key               string    `gorm:"column:key;primaryKey"`
	LastGCAt          time.Time `gorm:"column:last_gc_at;index;not null;default:1970-01-01 00:00:00"`
	LastCommitAt      time.Time `gorm:"column:last_commit_at;index;not null"`
	LastCommitBy      string    `gorm:"column:last_commit_by;index;not null"`
	LastCommitMessage string    `gorm:"column:last_commit_message;type:text"`
	UpdatedAt         time.Time `gorm:"column:updated_at;index;not null;autoUpdateTime"`
}

func (MirroredGitRepo) TableName() string {
	return "mirrored_git_repos"
}
