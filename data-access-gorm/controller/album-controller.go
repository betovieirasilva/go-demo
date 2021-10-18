package api

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

func (controller *albumController) GetAlbums(c *gin.Context) {
	albums, err := albumDao.FindAllAlbums(controller.db)
	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func (controller *albumController) GetAlbumById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)

	if err != nil {
		//TODO: => log.Fatal vai derrubar a aplciação, então o correto não é lançlar o log.Fatal, mas tratar a exceção corretamente
		//log.Fatal(err)
		c.IndentedJSON(http.StatusFound, gin.H{"message": "Informe um ID válido para realizar a busca do album"})
		return
	}

	album, errDb := albumDao.FindAlbumById(controller.db, id)
	if errDb != nil {
		msgErrorStr := errDb.Error() //pega a mensagem de erro retornada
		c.IndentedJSON(http.StatusFound, gin.H{"message": msgErrorStr})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func (controller *albumController) DeleteAlbumById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)

	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": "Informe um ID válido para realizar a exclusão do album"})
		return
	}

	deleted, err := albumDao.RemoveAlbum(controller.db, id)

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

func (controller *albumController) PostAlbums(c *gin.Context) {
	var newAlbum entity.Album

	//faz o paser do Json e alimenta na variável newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	id, err := albumDao.SaveAlbum(controller.db, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	newAlbum.ID = id //atualiza o ID dentro do próprio objeto garantindo que no retorno o JSON será retornado com ele, no caso de INSERT
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
