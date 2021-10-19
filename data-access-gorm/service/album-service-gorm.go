//https://gorm.io/docs/error_handling.html
package service

import (
	"example/data-access/entity"
	"fmt"

	"gorm.io/gorm"
)

func (service *albumServiceGorm) FindById(id int64) (entity.Album, error) {
	var album entity.Album
	result := service.conn.First(&album, id)
	if result.Error != nil {
		return album, result.Error
	}
	return album, nil

}

func (service *albumServiceGorm) FindAll() ([]entity.Album, error) {
	var albums []entity.Album
	result := service.conn.Find(&albums)
	if result.Error != nil {
		return nil, result.Error
	}
	return albums, nil
}

func (service *albumServiceGorm) Save(album entity.Album) (int64, error) {
	var result *gorm.DB
	if album.ID == 0 { //default em Go é sempre zero
		result = service.conn.Create(&album)
	} else {
		result = service.conn.Model(&album).Updates(album)
	}

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("Nenhum registro atualizado")
	}
	return (&album).ID, nil
}

func (service *albumServiceGorm) Remove(id int64) (bool, error) {
	//exclusão pelo ID
	result := service.conn.Delete(&entity.Album{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
