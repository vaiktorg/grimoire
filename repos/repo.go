package repos

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

type DBRepo struct {
	db *gorm.DB
}

// NewDBRepo returns
func NewDBRepo(db *gorm.DB) *DBRepo {
	return &DBRepo{db: db}
}

// All returns array of result
func (d *DBRepo) All(dst interface{}) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		tx.Find(dst)

		if dst == nil {
			return errors.New("no results found")
		}

		return nil
	}, &sql.TxOptions{})
}

// Find anything by ID in their gorm.Model struct.
func (d *DBRepo) Find(dst, crit interface{}, args ...interface{}) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		tx.Take(dst, crit, args)

		if dst == nil {
			return errors.New("no results found")
		}

		return nil
	}, &sql.TxOptions{})
}

// Persist saves if ID not found, and updates if ID found.
func (d *DBRepo) Persist(dst interface{}) error {
	return d.db.Transaction(func(tx *gorm.DB) error {

		if dst == nil {
			return errors.New("nothing to persist")
		}

		tx.Save(dst)
		return nil
	})
}

func (d *DBRepo) GetDB() *gorm.DB {
	return d.db
}
