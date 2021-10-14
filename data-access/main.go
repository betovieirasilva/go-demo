//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"example/data-access/conf"
	albumdao "example/data-access/dao"
	"example/data-access/model"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

//REST functions
func getAlbums(c *gin.Context) {
	albums, err := albumdao.FindAllAlbums(db)
	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)

	if err != nil {
		//TODO: => log.Fatal vai derrubar a aplciação, então o correto não é lançlar o log.Fatal, mas tratar a exceção corretamente
		//log.Fatal(err)
		c.IndentedJSON(http.StatusFound, gin.H{"message": "Informe um ID válido para realizar a busca do album"})
		return
	}

	album, errDb := albumdao.FindAlbumById(db, id)
	if errDb != nil {
		msgErrorStr := errDb.Error() //pega a mensagem de erro retornada
		c.IndentedJSON(http.StatusFound, gin.H{"message": msgErrorStr})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func deleteAlbumById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)

	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": "Informe um ID válido para realizar a exclusão do album"})
		return
	}

	deleted, err := albumdao.RemoveAlbum(db, id)

	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}

	if !deleted {
		c.IndentedJSON(http.StatusFound, gin.H{"message": "Nenhum registro encontrado!"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Registro excluído com sucesso!"})
}

func removeIndex(s []model.Album, index int) []model.Album {
	return append(s[:index], s[index+1:]...)
}

func postAlbums(c *gin.Context) {
	var newAlbum model.Album

	//faz o paser do Json e alimenta na variável newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	id, err := albumdao.SaveAlbum(db, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	newAlbum.ID = id //atualiza o ID dentro do próprio objeto garantindo que no retorno o JSON será retornado com ele, no caso de INSERT
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func main() {
	db = conf.Connection()

	//Routers da aplicação
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.POST("/albums", postAlbums)

	router.Run("localhost:9001")
}
