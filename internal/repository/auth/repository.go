package auth

import (
	"sync"

	def "sad/internal/repository"
	user "sad/internal/repository/models/user"

	"github.com/gofiber/fiber/v2"
)

var _ def.AuthRepository = (*repository)(nil)

type repository struct {
	data map[string]*user.User
	m    sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]*user.User),
	}
}

func (r *repository) GetByUUID(c *fiber.Ctx, userUUID string) (*user.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	user, ok := r.data[userUUID]
	if !ok {
		return nil, nil
	}

	return user, nil
}

func (r *repository) GetByLogin(c *fiber.Ctx, login string) (*user.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	var user *user.User
	for _, value := range r.data {
		if value.Login == login {
			user = value
			break
		}
	}

	return user, nil
}

func (r *repository) Create(c *fiber.Ctx, user *user.User) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.data[user.UUID] = user

	return nil
}
