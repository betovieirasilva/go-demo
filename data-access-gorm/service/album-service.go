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
type albumService struct {
	albumDao dao.AlbumDao
}

func NewAlbumService(_db *sql.DB) AlbumService {
	return &albumService{
		albumDao: dao.NewAlbumDao(_db),
	}
}

func (service *albumService) FindById(id int64) (entity.Album, error) {
	return service.albumDao.FindById(id)
}

func (service *albumService) FindAll() ([]entity.Album, error) {
	return service.albumDao.FindAll()
}

func (service *albumService) Save(album entity.Album) (int64, error) {
	return service.albumDao.Save(album)
}

func (service *albumService) Remove(id int64) (bool, error) {
	return service.albumDao.Remove(id)
}
