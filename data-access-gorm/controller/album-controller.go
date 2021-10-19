package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	albumDao "example/data-access/dao"
	entity "example/data-access/entity"

	"github.com/gin-gonic/gin"
)

type AlbumController interface {
	GetAlbums(c *gin.Context)
	GetAlbumById(c *gin.Context)
	DeleteAlbumById(c *gin.Context)
	PostAlbums(c *gin.Context)
}

type albumController struct {
	db *sql.DB
}

func NewAlbumController(_db *sql.DB) AlbumController {
	return &albumController{
		db: _db,
	}
}

func (c *albumController) GetAlbums(context *gin.Context) {
	albums, err := albumDao.FindAllAlbums(c.db)
	if err != nil {
		context.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, albums)
}

func (c *albumController) GetAlbumById(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 0, 64)

	if err != nil {
		//TODO: => log.Fatal vai derrubar a aplciação, então o correto não é lançlar o log.Fatal, mas tratar a exceção corretamente
		//log.Fatal(err)
		context.IndentedJSON(http.StatusFound, gin.H{"message": "Informe um ID válido para realizar a busca do album"})
		return
	}

	album, errDb := albumDao.FindAlbumById(c.db, id)
	if errDb != nil {
		msgErrorStr := errDb.Error() //pega a mensagem de erro retornada
		context.IndentedJSON(http.StatusFound, gin.H{"message": msgErrorStr})
		return
	}
	context.IndentedJSON(http.StatusOK, album)
}

func (c *albumController) DeleteAlbumById(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 0, 64)

	if err != nil {
		context.IndentedJSON(http.StatusFound, gin.H{"message": "Informe um ID válido para realizar a exclusão do album"})
		return
	}

	deleted, err := albumDao.RemoveAlbum(c.db, id)

	if err != nil {
		context.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}

	if !deleted {
		context.IndentedJSON(http.StatusFound, gin.H{"message": "Nenhum registro encontrado!"})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"message": "Registro excluído com sucesso!"})
}

func (c *albumController) PostAlbums(context *gin.Context) {
	var newAlbum entity.Album

	//faz o paser do Json e alimenta na variável newAlbum
	if err := context.BindJSON(&newAlbum); err != nil {
		return
	}

	id, err := albumDao.SaveAlbum(c.db, newAlbum)
	if err != nil {
		context.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	newAlbum.ID = id //atualiza o ID dentro do próprio objeto garantindo que no retorno o JSON será retornado com ele, no caso de INSERT
	context.IndentedJSON(http.StatusCreated, newAlbum)
}
