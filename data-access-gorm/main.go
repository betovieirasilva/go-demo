//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"

	"example/data-access/conf"
	"example/data-access/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/data-access/service"
)

var dbSqlConnection *sql.DB
var dbGormConnection *gorm.DB
var albumService service.AlbumService

func main() {

	dbGormConnection = conf.SetupDatabaseGorm()

	dbSqlConnection = conf.SetupDatabaseSqlDB()

	//TODO: [Giba] Implementar um provider de teste para enrtregar uma ou outra implementação do service
	//albumService := service.NewAlbumServiceSql(dbSqlConnection)
	albumService = service.NewAlbumServiceGorm(dbGormConnection)

	albumController := controller.NewAlbumController(albumService)

	//Routers da aplicação
	router := gin.Default()
	router.GET("/albums", albumController.GetAlbums)
	router.GET("/albums/:id", albumController.GetAlbumById)
	router.DELETE("/albums/:id", albumController.DeleteAlbumById)
	router.POST("/albums", albumController.PostAlbums)

	router.Run("localhost:9001")
}
