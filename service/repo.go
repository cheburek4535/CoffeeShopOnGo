package service

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB

	Product   *DefaultCrud[Product, int]
	Staff     *DefaultCrud[Staff, int]
	Sale      *DefaultCrud[Sale, int]
	Warehouse *DefaultCrud[Warehouse, int]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:        db,
		Product:   NewDefaultCrud[Product, int](db),
		Staff:     NewDefaultCrud[Staff, int](db),
		Sale:      NewDefaultCrud[Sale, int](db),
		Warehouse: NewDefaultCrud[Warehouse, int](db),
	}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(
		&Product{},
		&Staff{},
		&Sale{},
		&Warehouse{},
	)
}