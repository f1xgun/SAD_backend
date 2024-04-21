package users

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	usersModels "sad/internal/models/users"
	def "sad/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetById(c *fiber.Ctx, userId string) (*usersModels.UserCredentials, error) {
	query := "SELECT login, password, role FROM users WHERE uuid=$1"
	log.Printf("Fetching user by id: %s", userId)

	row := r.db.QueryRow(c.Context(), query, userId)

	userCredentials := &usersModels.UserCredentials{}
	err := row.Scan(&userCredentials.Login, &userCredentials.Password, &userCredentials.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("User not found with id: %s", userId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching user by id: %s, error: %v", userId, err)
		return nil, err
	}

	log.Printf("User fetched successfully by id: %s", userId)
	return userCredentials, nil
}

func (r *repository) GetByLogin(c *fiber.Ctx, login string) (*usersModels.UserRepoModel, error) {
	query := "SELECT uuid, login, password FROM users WHERE login=$1"
	log.Printf("Fetching user by login: %s", login)

	row := r.db.QueryRow(c.Context(), query, login)

	userCredentials := &usersModels.UserRepoModel{}
	err := row.Scan(&userCredentials.Id, &userCredentials.Login, &userCredentials.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("User not found with login: %s", login)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching user by login: %s, error: %v", login, err)
		return nil, err
	}

	log.Printf("User fetched successfully by login: %s", login)
	return userCredentials, nil
}

func (r *repository) Create(c *fiber.Ctx, user usersModels.User) error {
	query := "INSERT INTO users (uuid, name, login, password, role) VALUES (@uuid, @name, @login, @password, @role)"
	log.Printf("Creating user: %#v", user)

	args := pgx.NamedArgs{
		"uuid":     user.Id,
		"name":     user.Name,
		"login":    user.Login,
		"password": user.Password,
		"role":     user.Role,
	}
	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error creating user: %#v, error: %v", user, err)
	} else {
		log.Printf("User created successfully: %#v", user)
	}

	return err
}

func (r *repository) ChangeUserRole(c *fiber.Ctx, userId string, newRole usersModels.UserRole) error {
	query := "UPDATE users SET role=@role WHERE uuid=@uuid"
	log.Printf("Changing user role for userId: %s to new role: %#v", userId, newRole)

	args := pgx.NamedArgs{
		"role": newRole,
		"uuid": userId,
	}
	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error changing user role for userId: %s, error: %v", userId, err)
	} else {
		log.Printf("User role changed successfully for userId: %s", userId)
	}

	return err
}

func (r *repository) CheckUserExists(c *fiber.Ctx, userId string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE uuid=$1"
	var count int
	err := r.db.QueryRow(c.Context(), query, userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) GetUserInfo(c *fiber.Ctx, userId string) (*usersModels.UserInfoRepoModel, error) {
	query := "SELECT uuid, name, login, role FROM users WHERE uuid=$1"
	log.Printf("Fetching user by id: %s", userId)

	row := r.db.QueryRow(c.Context(), query, userId)

	userInfo := &usersModels.UserInfoRepoModel{}
	err := row.Scan(&userInfo.Id, &userInfo.Name, &userInfo.Login, &userInfo.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("User not found with id: %s", userId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching user by id: %s, error: %v", userId, err)
		return nil, err
	}

	log.Printf("User fetched successfully by id: %s", userId)
	return userInfo, nil
}
