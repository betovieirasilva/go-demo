package service

import (
	"example/data-access/entity"
)

func (service *albumServiceGorm) FindById(id int64) (entity.Album, error) {
	return service.albumDao.FindById(id)
}

func (service *albumServiceGorm) FindAll() ([]entity.Album, error) {
	return service.albumDao.FindAll()
}

func (service *albumServiceGorm) Save(album entity.Album) (int64, error) {
	return service.albumDao.Save(album)
}

func (service *albumServiceGorm) Remove(id int64) (bool, error) {
	return service.albumDao.Remove(id)
}
