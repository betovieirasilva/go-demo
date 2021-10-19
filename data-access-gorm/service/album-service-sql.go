package service

import (
	"example/data-access/entity"
)

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
