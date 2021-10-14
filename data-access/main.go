//https://golang.org/doc/tutorial/web-service-gin
//https://golang.org/doc/tutorial/database-access
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

//vars
type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

/*
var albums = []Album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry", Price: 56.99},
	{ID: 3, Title: "Sarah", Artist: "Sarah", Price: 56.99},
}*/

var db *sql.DB

//REST functions
func getAlbums(c *gin.Context) {
	albums, err := albumsAll()
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

	album, errDb := albumById(id)
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

	deleted, err := deleteAlbum(id)

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

func removeIndex(s []Album, index int) []Album {
	return append(s[:index], s[index+1:]...)
}

func postAlbums(c *gin.Context) {
	var newAlbum Album

	//faz o paser do Json e alimenta na variável newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	id, err := saveAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	newAlbum.ID = id //atualiza o ID dentro do próprio objeto garantindo que no retorno o JSON será retornado com ele, no caso de INSERT
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

//DB functions
func albumById(id int64) (Album, error) {
	var album Album

	row := db.QueryRow("SELECT * from album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("Não existe um album com o ID %d", id) //fmt.Errorf => transforma a mensagem em um erro
		}
		return album, fmt.Errorf("Erro ao buscar o algum com ID %d: %v", id, err)
	}
	//retorna o registro encontrato
	return album, nil
}

func albumsAll() ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE 1 = ? ORDER BY id ASC", 1) //param = 1 apenas para facilitar os testes com lista vazia (informe 2 para lista vazia)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar os registros de Albums. %v", err)
	}

	defer rows.Close() //executa quando a function terminar a execução
	empty := true
	for rows.Next() { //o mesmo que while em outras linguagens
		empty = false
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("Erro ao buscar a lista de albums %v", err)
		}
		albums = append(albums, album)
	}

	if empty {
		return nil, fmt.Errorf("Não existem registros a serem retornados!")
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Erro ao buscar a lista de albums %v", err)
	}

	return albums, nil
}

func insertAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES(?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("Erro ao inserir um registro em album %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Erro ao inserir um registro em album %v", err)
	}
	return id, nil
}

func updateAlbum(album Album) (int64, error) {
	result, err := db.Exec("UPDATE album set title = ?, artist = ?, price = ? WHERE id = ?", album.Title, album.Artist, album.Price, album.ID)
	if err != nil {
		return 0, fmt.Errorf("Erro ao atualizar um registro em album %v", err)
	}

	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Erro ao atualizar um registro em album %v", err)
	}

	if rowsUpdated == 0 { //se nenhum informação mudou, não dispara o UPDATE
		return 0, fmt.Errorf("Nenhum registro atualizado")
	}

	return album.ID, nil //retorn album.ID por padrão para facilitar seu uso no save
}

func saveAlbum(album Album) (int64, error) {
	if album.ID != 0 { //primitive value is zero by default
		return updateAlbum(album)
	}
	return insertAlbum(album)
}

func deleteAlbum(id int64) (bool, error) {
	result, err := db.Exec("DELETE FROM album WHERE id = ?", id)
	if err != nil {
		return false, fmt.Errorf("Erro ao excluir o registro do album %v", err)
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erro ao excluir um registro do album %v", err)
	}

	if rowsDeleted == 0 {
		return false, fmt.Errorf("Nenhum registro localizado para exclusão")
	}

	return true, nil
}

func main() {
	// Create connection
	//From => https://pkg.go.dev/github.com/go-sql-driver/mysql#Config
	cfg := mysql.Config{
		// @see export properties, get with os.Getenv("DBUSER"), etc
		//$ export DBUSER=root
		//$ export DBPASS=1
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		//User:   "root",
		//Passwd: "1",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
		// connector.go:95: could not use requested auth plugin 'mysql_native_password': this user requires mysql native password authentication.
		// 2021/10/14 09:40:00 this user requires mysql native password authentication.
		AllowNativePasswords: true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err) //exit
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr) //exit
		return
	}
	fmt.Println("Connected!")

	//Routers da aplicação
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.DELETE("/albums/:id", deleteAlbumById)
	router.POST("/albums", postAlbums)

	router.Run("localhost:9001")
}
