package main

import (
	"github.com/guoyk93/nas-tools/model"
	"github.com/guoyk93/nas-tools/model/dao"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/guoyk93/nas-tools/utils"
	"github.com/guoyk93/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	q *dao.Query
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	client := rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{}))
	rg.Must0(client.AutoMigrate(&model.MirroredGitRepo{}))

	q = dao.Use(client)

	analyzeDir("/volume1/mirrors/Git", "")
}

func checkGitDirEntries(entries []os.DirEntry) bool {
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == ".git" {
			return true
		}
	}
	var (
		foundRefs    bool
		foundObjects bool
	)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "refs" {
			foundRefs = true
		} else if entry.Name() == "objects" {
			foundObjects = true
		}
	}
	return foundRefs && foundObjects
}

func recordGitDir(dirRoot string, dirRel string) {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("failed to record dir:", dirRoot, dirRel, err.Error())
	}()
	defer rg.Guard(&err)

	repo := rg.Must(git.PlainOpen(filepath.Join(dirRoot, dirRel)))

	var (
		lastCommitAt      time.Time
		lastCommitBy      string
		lastCommitMessage string
	)

	commits := rg.Must(repo.CommitObjects())
	rg.Must0(commits.ForEach(func(commit *object.Commit) error {
		if commit.Author.When.After(lastCommitAt) {
			lastCommitAt = commit.Author.When
			lastCommitBy = commit.Author.String()
			lastCommitMessage = commit.Message
		}
		return nil
	}))

	if lastCommitAt.IsZero() {
		lastCommitAt = time.UnixMilli(0)
	}
	if lastCommitBy == "" {
		lastCommitBy = "unknown <unknown@unknown.com>"
	}
	if lastCommitMessage == "" {
		lastCommitMessage = "unknown"
	}

	record := rg.Must(q.MirroredGitRepo.Where(
		q.MirroredGitRepo.Key.Eq(dirRel),
	).Assign(
		q.MirroredGitRepo.LastCommitAt.Value(lastCommitAt),
		q.MirroredGitRepo.LastCommitBy.Value(lastCommitBy),
		q.MirroredGitRepo.LastCommitMessage.Value(lastCommitMessage),
	).FirstOrCreate())

	log.Println("recorded:", dirRoot, dirRel, lastCommitAt, lastCommitBy, lastCommitMessage)

	now := time.Now()

	if now.Sub(record.LastGCAt) > time.Hour*24*7 {
		// exec: git gc
		cmd := exec.Command("git", "-C", filepath.Join(dirRoot, dirRel), "gc")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		rg.Must0(cmd.Run())

		// record
		rg.Must(q.MirroredGitRepo.Where(
			q.MirroredGitRepo.Key.Eq(dirRel),
		).UpdateSimple(
			q.MirroredGitRepo.LastGCAt.Value(now),
		))
	}

	runtime.GC()
}

func analyzeDir(dirRoot string, dirRel string) {
	entries := rg.Must(os.ReadDir(filepath.Join(dirRoot, dirRel)))

	if checkGitDirEntries(entries) {
		recordGitDir(dirRoot, dirRel)
	} else {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if entry.Name() == "@eaDir" {
				continue
			}
			analyzeDir(dirRoot, filepath.Join(dirRel, entry.Name()))
		}
	}
}
