package teststore

import (
	"RestApi/internal/app/model"
	"RestApi/internal/app/store"
)

type Store struct {
	UserRepository *UserRepository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}
	return s.UserRepository
}
