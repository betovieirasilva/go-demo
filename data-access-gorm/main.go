//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"

	"example/data-access/conf"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/data-access/controller"
	"example/data-access/service"
)

var dbSql *sql.DB
var dbGorm *gorm.DB

func main() {

	dbGorm = conf.SetupDatabaseGorm()
	dbSql = conf.SetupDatabaseSqlDB()

	albumService := service.NewAlbumService(dbSql)
	albumController := controller.NewAlbumController(albumService)

	//Routers da aplicação
	router := gin.Default()
	router.GET("/albums", albumController.GetAlbums)
	router.GET("/albums/:id", albumController.GetAlbumById)
	router.DELETE("/albums/:id", albumController.DeleteAlbumById)
	router.POST("/albums", albumController.PostAlbums)

	router.Run("localhost:9001")
}
