package conf

import (
	"example/data-access/entity"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabaseGorm() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:1@tcp(127.0.0.1:3306)/recordings?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), //habilita o log para exibir todos os SQLs executados
	})

	if err != nil {
		log.Fatal(err) //exit
	}

	db.AutoMigrate(&entity.Album{})

	return db
}
