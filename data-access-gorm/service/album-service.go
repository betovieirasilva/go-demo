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

type albumServiceSql struct {
	albumDao dao.AlbumDao
}

func NewAlbumServiceSql(_db *sql.DB) AlbumService {
	return &albumServiceSql{
		albumDao: dao.NewAlbumDao(_db),
	}
}

type albumServiceGorm struct {
	albumDao dao.AlbumDao
}

func NewAlbumServiceGorm(_db *sql.DB) AlbumService {
	return &albumServiceGorm{
		albumDao: dao.NewAlbumDao(_db),
	}
}
