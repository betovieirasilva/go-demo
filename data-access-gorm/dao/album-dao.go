package dao

import (
	"database/sql"
	"fmt"

	"example/data-access/entity"
)

type AlbumDao interface {
	FindAlbumById(id int64) (entity.Album, error)
	FindAllAlbums() ([]entity.Album, error)
	SaveAlbum(album entity.Album) (int64, error)
	RemoveAlbum(id int64) (bool, error)
}

type albumDao struct {
	db *sql.DB
}

func NewAlbumDao(_db *sql.DB) AlbumDao {
	return &albumDao{
		db: _db,
	}
}

func (dao *albumDao) FindAlbumById(id int64) (entity.Album, error) {
	var album entity.Album

	row := dao.db.QueryRow("SELECT * from album WHERE id = ?", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("Não existe um album com o ID %d", id) //fmt.Errorf => transforma a mensagem em um erro
		}
		return album, fmt.Errorf("Erro ao buscar o algum com ID %d: %v", id, err)
	}
	//retorna o registro encontrato
	return album, nil
}

func (dao *albumDao) FindAllAlbums() ([]entity.Album, error) {
	var albums []entity.Album

	rows, err := dao.db.Query("SELECT * FROM album WHERE 1 = ? ORDER BY id ASC", 1) //param = 1 apenas para facilitar os testes com lista vazia (informe 2 para lista vazia)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar os registros de Albums. %v", err)
	}

	defer rows.Close() //executa quando a function terminar a execução
	empty := true
	for rows.Next() { //o mesmo que while em outras linguagens
		empty = false
		var album entity.Album
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

func (dao *albumDao) SaveAlbum(album entity.Album) (int64, error) {
	if album.ID != 0 { //primitive value is zero by default
		return updateAlbum(dao.db, album)
	}
	return insertAlbum(dao.db, album)
}

func insertAlbum(db *sql.DB, album entity.Album) (int64, error) {
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

func updateAlbum(db *sql.DB, album entity.Album) (int64, error) {
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

func (dao *albumDao) RemoveAlbum(id int64) (bool, error) {
	result, err := dao.db.Exec("DELETE FROM album WHERE id = ?", id)
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
