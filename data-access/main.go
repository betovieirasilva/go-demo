//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"

	"example/data-access/conf"

	"github.com/gin-gonic/gin"

	albumapi "example/data-access/api"
)

var db *sql.DB

func main() {
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
