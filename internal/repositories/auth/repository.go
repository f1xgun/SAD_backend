package auth

import (
	userModels "sad/internal/models/user"
	def "sad/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

var _ def.AuthRepository = (*repository)(nil)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetByUUID(c *fiber.Ctx, userUUID string) (*userModels.UserRepoModel, error) {
	query := "SELECT uuid, login, password FROM users WHERE uuid=$1"

	row := r.db.QueryRow(c.Context(), query, userUUID)

	userCredentials := &userModels.UserRepoModel{}
	err := row.Scan(&userCredentials.UUID, &userCredentials.Login, &userCredentials.Password)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return userCredentials, nil
}

func (r *repository) GetByLogin(c *fiber.Ctx, login string) (*userModels.UserRepoModel, error) {
	query := "SELECT uuid, login, password FROM users WHERE login=$1"

	row := r.db.QueryRow(c.Context(), query, login)

	userCredentials := &userModels.UserRepoModel{}
	err := row.Scan(&userCredentials.UUID, &userCredentials.Login, &userCredentials.Password)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return userCredentials, nil
}

func (r *repository) Create(c *fiber.Ctx, user userModels.User) error {
	query := "INSERT INTO users (uuid, name, login, password, role) VALUES (@uuid, @name, @login, @password, @role)"
	args := pgx.NamedArgs{
		"uuid":     user.UUID,
		"name":     user.Name,
		"login":    user.Login,
		"password": user.Password,
		"role":     user.Role,
	}
	_, err := r.db.Exec(c.Context(), query, args)

	return err
}
