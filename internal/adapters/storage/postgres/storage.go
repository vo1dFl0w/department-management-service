package postgres

import (
	"github.com/vo1dFl0w/department-management-service/internal/repository"
	"gorm.io/gorm"
)

type Storage struct {
	DB                 *gorm.DB
	postgresRepo 		repository.Postgres
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Postgres() repository.Postgres {
	return s.postgresRepo
}