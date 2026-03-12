package repositories

import (
	"financialcontrol/internal/store"
)

type Repository struct {
	store store.Store
}

func NewRepository(store store.Store) Repository {
	return Repository{store: store}
}
