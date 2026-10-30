package main

import "sync"

// ==================================================
// MEMORY REPOSITORY
// ==================================================
type MemoryUserRepository struct {
	mutex  sync.RWMutex
	users  []User
	nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: []User{
			{
				ID:    1,
				Name:  "Imam",
				Email: "imam@example.com",
			},
			{
				ID:    2,
				Name:  "Syafii",
				Email: "syafii@example.com",
			},
		},
		nextID: 3,
	}
}

func (repository *MemoryUserRepository) List() []User {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	users := make([]User, len(repository.users))
	copy(users, repository.users)

	return users
}

func (repository *MemoryUserRepository) FindByID(id int) (User, bool) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	for _, user := range repository.users {
		if user.ID == id {
			return user, true
		}
	}

	return User{}, false
}

func (repository *MemoryUserRepository) Create(name string, email string) User {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	user := User{
		ID:    repository.nextID,
		Name:  name,
		Email: email,
	}
	repository.nextID++
	repository.users = append(repository.users, user)

	return user
}

func (repository *MemoryUserRepository) UpdateByID(id int, name string, email string) (User, bool) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for index, user := range repository.users {
		if user.ID == id {
			user.Name = name
			user.Email = email

			repository.users[index] = user
			return user, true
		}
	}

	return User{}, false
}

func (repository *MemoryUserRepository) DeleteByID(id int) bool {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for index, user := range repository.users {
		if user.ID == id {
			repository.users = append(
				repository.users[:index],
				repository.users[index+1:]...,
			)
			return true
		}
	}
	return false
}
