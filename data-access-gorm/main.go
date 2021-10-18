//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"

	"example/data-access/conf"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	albumapi "example/data-access/api"
)

var db *sql.DB
var dbGorm *gorm.DB

func main() {

	dbGorm = conf.SetupDatabase()

	db = conf.Connection()
	albumapi.SetConnection(db)

	//Routers da aplicação
	router := gin.Default()
	router.GET("/albums", albumapi.GetAlbums)
	router.GET("/albums/:id", albumapi.GetAlbumById)
	router.DELETE("/albums/:id", albumapi.DeleteAlbumById)
	router.POST("/albums", albumapi.PostAlbums)

	router.Run("localhost:9001")
}
