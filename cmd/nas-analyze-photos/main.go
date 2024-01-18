package main

import (
	"errors"
	"flag"
	"os"

	"github.com/yankeguo/nas-tools/model/dao"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	allowedExtensions = []string{
		".jpg",
		".png",
		".gif",
		".mov",
		".heic",
		".mp4",
	}
)

var (
	db *dao.Query

	optGroup string
	optDir   string
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	db = dao.Use(rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{})))

	flag.StringVar(&optGroup, "group", "", "group")
	flag.StringVar(&optDir, "dir", "", "dir")
	flag.Parse()

	if optGroup == "" {
		err = errors.New("group is required")
		return
	}

	if optDir == "" {
		err = errors.New("dir is required")
		return
	}
}
