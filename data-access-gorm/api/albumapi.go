package api

import (
	"database/sql"
	"net/http"
	"strconv"

	albumdao "example/data-access/dao"
	"example/data-access/model"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func SetConnection(_db *sql.DB) {
	db = _db
}

func GetAlbums(c *gin.Context) {
	albums, err := albumdao.FindAllAlbums(db)
	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func GetAlbumById(c *gin.Context) {
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

func DeleteAlbumById(c *gin.Context) {
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

func PostAlbums(c *gin.Context) {
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
