package service

import (
	"example/data-access/entity"
)

func (service *albumServiceGorm) FindById(id int64) (entity.Album, error) {
	var album entity.Album
	service.conn.First(&album, id)
	return album, nil

}

func (service *albumServiceGorm) FindAll() ([]entity.Album, error) {
	var albums []entity.Album
	service.conn.Find(&albums)
	return albums, nil
}

func (service *albumServiceGorm) Save(album entity.Album) (int64, error) {
	if album.ID == 0 {
		service.conn.Create(&album)
	} else {
		service.conn.Model(&album).Updates(album)
	}
	return 0, nil
}

func (service *albumServiceGorm) Remove(id int64) (bool, error) {
	album, error := service.FindById(id)
	if error != nil {
		return false, error
	}
	service.conn.Delete(&album, id)
	return true, nil
}
