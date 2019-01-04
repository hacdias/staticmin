package users

import (
	"github.com/filebrowser/filebrowser/errors"
)

// StorageBackend is the interface to implement for a users storage.
type StorageBackend interface {
	GetByID(uint) (*User, error)
	GetByUsername(string) (*User, error)
	Gets() ([]*User, error)
	Save(u *User) error
	Update(u *User, fields ...string) error
	DeleteByID(uint) error
	DeleteByUsername(string) error
}

// Storage is a users storage.
type Storage struct {
	back StorageBackend
}

// NewStorage creates a users storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{back: back}
}

// Get allows you to get a user by its name or username. The provided
// id must be a string for username lookup or a uint for id lookup. If id
// is neither, a ErrInvalidDataType will be returned.
func (s *Storage) Get(id interface{}) (*User, error) {
	var (
		user *User
		err  error
	)

	switch id.(type) {
	case string:
		user, err = s.back.GetByUsername(id.(string))
	case uint:
		user, err = s.back.GetByID(id.(uint))
	default:
		return nil, errors.ErrInvalidDataType
	}

	if err != nil {
		return nil, err
	}

	user.Clean()
	return user, err
}

// Gets gets a list of all users.
func (s *Storage) Gets() ([]*User, error) {
	users, err := s.back.Gets()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.Clean()
	}

	return users, err
}

// Update updates a user in the database.
func (s *Storage) Update(user *User, fields ...string) error {
	err := user.Clean(fields...)
	if err != nil {
		return err
	}

	return s.back.Update(user, fields...)
}

// Save saves the user in a storage.
func (s *Storage) Save(user *User) error {
	if err := user.Clean(); err != nil {
		return err
	}

	return s.back.Save(user)
}

// Delete allows you to delete a user by its name or username. The provided
// id must be a string for username lookup or a uint for id lookup. If id
// is neither, a ErrInvalidDataType will be returned.
func (s *Storage) Delete(id interface{}) (err error) {
	switch id.(type) {
	case string:
		err = s.back.DeleteByUsername(id.(string))
	case uint:
		err = s.back.DeleteByID(id.(uint))
	default:
		err = errors.ErrInvalidDataType
	}

	return
}
