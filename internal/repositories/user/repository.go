package users

import (
	"context"
	"errors"
	"log"
	usersModels "sad/internal/models/users"
	def "sad/internal/repositories"

	"github.com/jackc/pgx/v5/pgxpool"

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

func (r *repository) GetById(ctx context.Context, userId string) (*usersModels.UserCredentials, error) {
	query := "SELECT login, password, role FROM users WHERE uuid=$1"
	log.Printf("Fetching user by id: %s", userId)

	row := r.db.QueryRow(ctx, query, userId)

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

func (r *repository) GetByLogin(ctx context.Context, login string) (*usersModels.UserRepoModel, error) {
	query := "SELECT uuid, login, password FROM users WHERE login=$1"
	log.Printf("Fetching user by login: %s", login)

	row := r.db.QueryRow(ctx, query, login)

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

func (r *repository) Create(ctx context.Context, user usersModels.User) error {
	query := "INSERT INTO users (last_name, name, login, password, role"

	args := pgx.NamedArgs{
		"name":      user.Name,
		"login":     user.Login,
		"password":  user.Password,
		"role":      user.Role,
		"last_name": user.LastName,
	}

	if user.MiddleName != nil {
		query += ", middle_name"
		args["middle_name"] = *user.MiddleName
	}

	query += ") VALUES (@last_name, @name, @login, @password, @role"

	if user.MiddleName != nil {
		query += ", @middle_name"
	}

	query += ")"
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		log.Printf("Error creating user: %#v, error: %v", user, err)
	} else {
		log.Printf("User created successfully: %#v", user)
	}

	return err
}

func (r *repository) ChangeUserInfo(ctx context.Context, userId string, newRole usersModels.UserRole, newName string) error {
	query := "UPDATE users SET role=@role, name=@name WHERE uuid=@uuid"
	log.Printf("Changing user info for userId: %s to new role: %#v and new name: %#v", userId, newRole, newName)

	args := pgx.NamedArgs{
		"role": newRole,
		"uuid": userId,
		"name": newName,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		log.Printf("Error changing user info for userId: %s, error: %v", userId, err)
	} else {
		log.Printf("User info changed successfully for userId: %s", userId)
	}

	return err
}

func (r *repository) CheckUserExists(ctx context.Context, userId string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE uuid=$1"
	var count int
	err := r.db.QueryRow(ctx, query, userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) GetUserInfo(ctx context.Context, userId string) (*usersModels.UserInfoRepoModel, error) {
	query := "SELECT uuid, name, login, role, last_name, middle_name FROM users WHERE uuid=$1"
	log.Printf("Fetching user by id: %s", userId)

	row := r.db.QueryRow(ctx, query, userId)

	userInfo := &usersModels.UserInfoRepoModel{}
	err := row.Scan(
		&userInfo.Id,
		&userInfo.Name,
		&userInfo.Login,
		&userInfo.Role,
		&userInfo.LastName,
		&userInfo.MiddleName,
	)
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

func (r *repository) GetUsersInfo(ctx context.Context) ([]usersModels.UserInfoRepoModel, error) {
	query := "SELECT uuid, name, login, role, last_name, middle_name FROM users"
	log.Printf("Fetching users")

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var usersInfo []usersModels.UserInfoRepoModel
	for rows.Next() {
		var userInfo usersModels.UserInfoRepoModel

		if err := rows.Scan(
			&userInfo.Id,
			&userInfo.Name,
			&userInfo.Login,
			&userInfo.Role,
			&userInfo.LastName,
			&userInfo.MiddleName,
		); err != nil {
			log.Printf("Error scaning user: %v", err)
			continue
		}

		if userInfo.Id.Valid {
			usersInfo = append(usersInfo, userInfo)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over users: %v", err)
		return nil, err
	}

	log.Printf("Users info fetched successfully")
	return usersInfo, nil
}

func (r *repository) DeleteUser(ctx context.Context, userId string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	query := "DELETE FROM users WHERE uuid=$1"

	if _, err = tx.Exec(ctx, query, userId); err != nil {
		log.Printf("Error deleting user with id %s, err %v", userId, err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAvailableTeachers(ctx context.Context, teacherName string) ([]usersModels.UserInfoRepoModel, error) {
	query := `
	SELECT uuid, name, login, last_name, middle_name
	FROM users
	WHERE last_name || name || middle_name LIKE '%' || $1 || '%'
	AND role = 'teacher'
	`
	log.Printf("Fetching teachers by name: %s", teacherName)

	rows, err := r.db.Query(ctx, query, teacherName)
	if err != nil {
		log.Printf("Error fetching teachers by name: %v", err)
		return nil, err
	}
	defer rows.Close()

	var teachers []usersModels.UserInfoRepoModel
	for rows.Next() {
		var teacher usersModels.UserInfoRepoModel

		if err := rows.Scan(
			&teacher.Id,
			&teacher.Name,
			&teacher.Login,
			&teacher.LastName,
			&teacher.MiddleName,
		); err != nil {
			log.Printf("Error scanning teacher: %v", err)
			continue
		}

		if teacher.Id.Valid {
			teachers = append(teachers, teacher)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over users: %v", err)
		return nil, err
	}

	log.Printf("Available teachers fetched successfully")
	return teachers, nil
}
