package storage

import (
	pg "github.com/go-pg/pg/v9"
)

type StorageHandler struct {
	Db *pg.DB
}

func (h *StorageHandler) Save(v interface{}) error {
	return h.Db.Insert(v)
}
