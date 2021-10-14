package albumdao

import (
	"database/sql"
	"fmt"

	"example/data-access/model"
)

func FindAlbumById(db *sql.DB, id int64) (model.Album, error) {
	var album model.Album

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

func FindAllAlbums(db *sql.DB) ([]model.Album, error) {
	var albums []model.Album

	rows, err := db.Query("SELECT * FROM album WHERE 1 = ? ORDER BY id ASC", 1) //param = 1 apenas para facilitar os testes com lista vazia (informe 2 para lista vazia)
	if err != nil {
		return nil, fmt.Errorf("Erro ao buscar os registros de Albums. %v", err)
	}

	defer rows.Close() //executa quando a function terminar a execução
	empty := true
	for rows.Next() { //o mesmo que while em outras linguagens
		empty = false
		var album model.Album
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

func SaveAlbum(db *sql.DB, album model.Album) (int64, error) {
	if album.ID != 0 { //primitive value is zero by default
		return UpdateAlbum(db, album)
	}
	return InsertAlbum(db, album)
}

func InsertAlbum(db *sql.DB, album model.Album) (int64, error) {
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

func UpdateAlbum(db *sql.DB, album model.Album) (int64, error) {
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

func RemoveAlbum(db *sql.DB, id int64) (bool, error) {
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
