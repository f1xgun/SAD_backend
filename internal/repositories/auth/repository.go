package auth

import (
	"sync"

	userModels "sad/internal/models/user"
	def "sad/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

var _ def.AuthRepository = (*repository)(nil)

type repository struct {
	data map[string]*userModels.User
	m    sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]*userModels.User),
	}
}

func (r *repository) GetByUUID(c *fiber.Ctx, userUUID string) (*userModels.UserCredentials, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	user, ok := r.data[userUUID]
	if !ok {
		return nil, nil
	}

	userCredentials := &userModels.UserCredentials{
		Login:    user.Login,
		Password: user.Password,
	}

	return userCredentials, nil
}

func (r *repository) GetByLogin(c *fiber.Ctx, login string) (*userModels.UserCredentials, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	var user *userModels.User
	for _, value := range r.data {
		if value.Login == login {
			user = value
			break
		}
	}

	userCredentials := &userModels.UserCredentials{
		Login:    user.Login,
		Password: user.Password,
	}

	return userCredentials, nil
}

func (r *repository) Create(c *fiber.Ctx, user userModels.User) error {
	r.m.Lock()
	defer r.m.Unlock()

	r.data[user.UUID] = &user

	return nil
}
