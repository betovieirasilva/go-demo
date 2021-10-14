//https://golang.org/doc/tutorial/web-service-gin
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `jso:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry", Price: 56.99},
	{ID: "3", Title: "Sarah", Artist: "Sarah", Price: 56.99},
}

//retorna uma lista de albuns para testes iniciais do Go
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	//busca a informação na lista que sá em memória
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusFound, gin.H{"message": "ambum not found"})
}

//novo exemplo, não está no tutorial
func deleteAlbumById(c *gin.Context) {
	id := c.Param("id")
	index := 0
	//busca a informação na lista que sá em memória
	for _, a := range albums {
		if a.ID == id {
			albums = removeIndex(albums, index)
			c.IndentedJSON(http.StatusCreated, gin.H{"message": "Remove sucess"})
			return
		}
		index++
	}
	c.IndentedJSON(http.StatusFound, gin.H{"message": "Album not found"})
}

func removeIndex(s []album, index int) []album {
	return append(s[:index], s[index+1:]...)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	//faz o paser do Json e alimenta na variável newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	//add newValue into newAlbum
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.POST("/albums", postAlbums)

	router.Run("localhost:9001")
}
