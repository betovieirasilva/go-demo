//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"
	"log"
	"time"

	"example/data-access/conf"
	"example/data-access/controller"
	validator "example/data-access/entity-validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example/data-access/service"
)

var dbSqlConnection *sql.DB
var dbGormConnection *gorm.DB
var albumService service.AlbumService

func printRawJson(c *gin.Context) {
	// Ao pegar o JSON com Raw gera EOF nas próximas fases da requisição quando delamos a conversão para o serviço e não realizar uma cópia do context
	// => https://gin-gonic.com/docs/examples/goroutines-inside-a-middleware/ [Ver melhor, deixar para o final]
	cCopy := c.Copy()
	jsonData, err := cCopy.GetRawData()
	if err != nil {
		log.Println("Error: ", err)
	}
	log.Print("#PATH    => ", cCopy.Request.URL.Path)
	log.Print("#RawJSON => ", string(jsonData))
}

func LoggerHandlerDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("#Before request")

		//printRawJson(c)

		t := time.Now()

		// Set example variable context
		c.Set("example", "12345")

		c.Next()

		log.Println("#After request")
		latency := time.Since(t)
		log.Print("latency=>", latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println("Status: ", status)
	}
}

func LoggerHandlerV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("#HandlerV1 OK")
		c.Next()
	}
}

func LoggerHandlerV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("#HandlerV2 OK")
		c.Next()
	}
}

func main() {

	dbGormConnection = conf.SetupDatabaseGorm()

	dbSqlConnection = conf.SetupDatabaseSqlDB()

	//TODO: [Giba] Implementar um provider de teste para enrtregar uma ou outra implementação do service
	//albumService := service.NewAlbumServiceSql(dbSqlConnection)
	albumService = service.NewAlbumServiceGorm(dbGormConnection)
	albumController := controller.NewAlbumController(albumService)

	//validators
	validator.NewEntityValidator().Register()

	//Routers da aplicação geranciados por versão de API (v1/v2) e com handler geral e por grupo
	router := gin.Default()

	router.Use(LoggerHandlerDefault())

	apiV1 := router.Group("/v1")
	{
		apiV1.Use(LoggerHandlerV1())
		apiV1.GET("/albums", albumController.GetAlbums)
		apiV1.GET("/albums/:id", albumController.GetAlbumById)
		apiV1.DELETE("/albums/:id", albumController.DeleteAlbumById)
		apiV1.POST("/albums", albumController.PostAlbums)
	}

	apiV2 := router.Group("/v2")
	{
		apiV2.Use(LoggerHandlerV2())
		apiV2.GET("/albums", albumController.GetAlbums)
		apiV2.GET("/albums/:id", albumController.GetAlbumById)
		apiV2.DELETE("/albums/:id", albumController.DeleteAlbumById)
		apiV2.POST("/albums", albumController.PostAlbums)
	}
	router.Run("localhost:9001")
}
