package main

import (
	"os"

	"github.com/yankeguo/nas-tools/model"
	"github.com/yankeguo/nas-tools/utils"
	"github.com/yankeguo/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var allModels = []any{
	model.ArchivedBundle{},
	model.ArchivedFile{},
	model.ArchivedFileIgnore{},
	model.MirroredGitRepo{},
}

func main() {
	var err error
	defer utils.Exit(&err)
	defer rg.Guard(&err)

	mysqlDSN := os.Getenv("MYSQL_DSN")

	if mysqlDSN == "" {
		mysqlDSN = "root:qwerty@tcp(127.0.0.1:3306)/automata?charset=utf8mb4&parseTime=True&loc=Local"
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./model/dao",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
	})

	db := rg.Must(gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})).Debug()

	rg.Must0(db.AutoMigrate(allModels...))

	g.UseDB(db)

	g.ApplyBasic(allModels...)

	g.Execute()
}
