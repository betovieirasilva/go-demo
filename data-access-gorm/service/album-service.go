package service

import (
	"database/sql"

	"example/data-access/dao"
	"example/data-access/entity"
)

type AlbumService interface {
	FindById(id int64) (entity.Album, error)
	FindAll() ([]entity.Album, error)
	Save(album entity.Album) (int64, error)
	Remove(id int64) (bool, error)
}

//TODO: [Giba] Refatorar movendo para outro arquivo
type albumServiceSql struct {
	albumDao dao.AlbumDao
}

func NewAlbumServiceSql(_db *sql.DB) AlbumService {
	return &albumServiceSql{
		albumDao: dao.NewAlbumDao(_db),
	}
}

func (service *albumServiceSql) FindById(id int64) (entity.Album, error) {
	return service.albumDao.FindById(id)
}

func (service *albumServiceSql) FindAll() ([]entity.Album, error) {
	return service.albumDao.FindAll()
}

func (service *albumServiceSql) Save(album entity.Album) (int64, error) {
	return service.albumDao.Save(album)
}

func (service *albumServiceSql) Remove(id int64) (bool, error) {
	return service.albumDao.Remove(id)
}
