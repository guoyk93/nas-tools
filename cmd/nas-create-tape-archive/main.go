package main

import (
	"github.com/yankeguo/nas-tools/model/dao"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

const (
	dirArchives = "/volume1/archives"
	dirTape     = "/volume1/tape"
)

var (
	optDebug bool

	db *dao.Query
)

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	// create db
	{
		client := rg.Must(gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")), &gorm.Config{}))
		if optDebug {
			client = client.Debug()
		}
		db = dao.Use(client)
	}
}
