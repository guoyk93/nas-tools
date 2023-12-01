package main

import (
	"github.com/guoyk93/nas-tools/model"
	"github.com/guoyk93/nas-tools/utils"
	"github.com/guoyk93/rg"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
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
		mysqlDSN = "root:qwertyqwerty@tcp(127.0.0.1:3306)/automata?charset=utf8mb4&parseTime=True&loc=Local"
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
