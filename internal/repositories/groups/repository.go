package groups

import (
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	groupsModels "sad/internal/models/groups"
	usersModels "sad/internal/models/users"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(c *fiber.Ctx, group groupsModels.Group) error {
	query := "INSERT INTO groups (id, number) VALUES (@id, @number)"
	log.Printf("Creating group: %#v", group)

	args := pgx.NamedArgs{
		"id":     group.Id,
		"number": group.Number,
	}
	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error creating group: %#v, error: %v", group, err)
	} else {
		log.Printf("Group created successfully: %#v", group)
	}

	return err
}

func (r *repository) GetById(c *fiber.Ctx, groupId string) (*groupsModels.GroupRepoModel, error) {
	query := "SELECT id, number FROM groups WHERE id=$1"
	log.Printf("Fetching group by id: %s", groupId)

	row := r.db.QueryRow(c.Context(), query, groupId)

	group := &groupsModels.GroupRepoModel{}
	err := row.Scan(&group.Id, &group.Number)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Group not found with id: %s", groupId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching group by id: %s, error: %v", groupId, err)
		return nil, err
	}

	log.Printf("Group fetched successfully by id: %s", groupId)
	return group, nil
}

func (r *repository) GetAll(c *fiber.Ctx) ([]groupsModels.GroupRepoModel, error) {
	query := "SELECT id, number FROM groups"
	log.Printf("Fetching all groups")

	rows, err := r.db.Query(c.Context(), query)
	if err != nil {
		log.Printf("Error fetching all groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []groupsModels.GroupRepoModel
	for rows.Next() {
		var group groupsModels.GroupRepoModel
		if err := rows.Scan(&group.Id, &group.Number); err != nil {
			log.Printf("Error scanning group: %v", err)
			continue
		}
		if group.Id.Valid {
			groups = append(groups, group)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over groups: %v", err)
		return nil, err
	}

	log.Printf("All groups fetched successfully")
	return groups, nil
}

func (r *repository) AddUserToGroup(c *fiber.Ctx, groupId string, userId string) error {
	query := "INSERT INTO users_groups (user_id, group_id) VALUES (@user_id, @group_id)"
	args := pgx.NamedArgs{
		"user_id":  userId,
		"group_id": groupId,
	}
	log.Printf("Adding user %s to group %s", userId, groupId)

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error adding user to group: user_id=%s, group_id=%s, error: %v", userId, groupId, err)
		return err
	}

	log.Printf("User added to group successfully: user_id=%s, group_id=%s", userId, groupId)
	return nil
}

func (r *repository) DeleteUserFromGroup(c *fiber.Ctx, groupId string, userId string) error {
	query := "DELETE FROM users_groups WHERE user_id=@user_id AND group_id=@group_id"
	args := pgx.NamedArgs{
		"user_id":  userId,
		"group_id": groupId,
	}
	log.Printf("Deleting user %s from group %s", userId, groupId)

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error deleting user from group: user_id=%s, group_id=%s, error: %v", userId, groupId, err)
		return err
	}

	log.Printf("User deleted from group successfully: user_id=%s, group_id=%s", userId, groupId)
	return nil
}

func (r *repository) IsUserInGroup(c *fiber.Ctx, groupId, userId string) (bool, error) {
	log.Printf("Checking if user with ID '%s' is in group with ID '%s'\n", userId, groupId)

	query := "SELECT COUNT(*) FROM users_groups WHERE user_id=@user_id and group_id=@group_id"
	args := pgx.NamedArgs{
		"user_id":  userId,
		"group_id": groupId,
	}

	var count int
	err := r.db.QueryRow(c.Context(), query, args).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user '%s' is in group '%s': %v\n", userId, groupId, err)
		return false, err
	}

	if count > 0 {
		log.Printf("User '%s' is in group '%s'\n", userId, groupId)
	} else {
		log.Printf("User '%s' is not in group '%s'\n", userId, groupId)
	}

	return count > 0, nil
}

func (r *repository) GetByIdWithUsers(c *fiber.Ctx, groupId string) (*groupsModels.GroupWithUsersRepo, error) {
	query := `
		SELECT g.id, g.number, u.uuid, u.login, u.name
		FROM groups g
		LEFT JOIN users_groups ug ON ug.group_id = g.id
		LEFT JOIN users u ON u.uuid = ug.user_id
		WHERE g.id = @group_id;
	`
	args := pgx.NamedArgs{
		"group_id": groupId,
	}

	rows, err := r.db.Query(c.Context(), query, args)
	if err != nil {
		log.Printf("Error fetching group with users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var group groupsModels.GroupWithUsersRepo
	var users []usersModels.UserInfoRepoModel
	for rows.Next() {
		var groupId, groupNumber sql.NullString
		var user usersModels.UserInfoRepoModel

		if err := rows.Scan(&groupId, &groupNumber, &user.Id, &user.Login, &user.Name); err != nil {
			log.Printf("Error scanning group: %v", err)
			continue
		}

		if groupId.Valid {
			group.Id = groupId
			group.Number = groupNumber
		}

		if user.Id.Valid {
			users = append(users, user)
		}
	}

	group.Users = users

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over users_groups: %v", err)
		return nil, err
	}

	log.Printf("Group with users fetched successfully")
	return &group, nil
}

func (r *repository) DeleteGroup(c *fiber.Ctx, groupId string) error {
	tx, err := r.db.Begin(c.Context())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(c.Context()); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	args := pgx.NamedArgs{
		"group_id": groupId,
	}

	query := "DELETE FROM users_groups WHERE group_id=@group_id"

	if _, err = tx.Exec(c.Context(), query, args); err != nil {
		log.Printf("Error deleting users from group with id %s , err: %v", groupId, err)
		return err
	}
	log.Printf("Delete users from group with id %s successfully", groupId)

	query = "DELETE FROM groups_subjects WHERE group_id=@group_id"
	if _, err = tx.Exec(c.Context(), query, args); err != nil {
		log.Printf("Error deleting subjects from group with id %s , err: %v", groupId, err)
		return err
	}
	log.Printf("Delete subjects with id %s successfully", groupId)

	query = "DELETE FROM groups WHERE id=@group_id"

	if _, err = tx.Exec(c.Context(), query, args); err != nil {
		log.Printf("Error deleting group with id %s , err: %v", groupId, err)
		return err
	}

	if err = tx.Commit(c.Context()); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateGroup(c *fiber.Ctx, group groupsModels.Group) error {
	query := "UPDATE groups SET number=@number WHERE id=@group_id"
	args := pgx.NamedArgs{
		"group_id": group.Id,
		"number":   group.Number,
	}

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error updating group with id %s, err: %v", group.Id, err)
	} else {
		log.Printf("Update group with id %s successfully", group.Id)
	}

	return err
}

func (r *repository) CheckGroupExists(c *fiber.Ctx, groupId string) (bool, error) {
	query := "SELECT COUNT(*) FROM groups WHERE id=$1"
	var count int
	err := r.db.QueryRow(c.Context(), query, groupId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) GetAvailableNewUsers(c *fiber.Ctx, groupId, login string) ([]usersModels.UserInfoRepoModel, error) {
	query := `
	SELECT uuid, name, login 
	FROM users u 
	WHERE login LIKE '%' || @login || '%'
	AND NOT EXISTS (
		SELECT 1 
		FROM users_groups ug 
		WHERE u.uuid = ug.user_id 
	)
	`

	args := pgx.NamedArgs{
		"login":   login,
		"groupId": groupId,
	}

	rows, err := r.db.Query(c.Context(), query, args)
	if err != nil {
		log.Printf("Error fetching group with users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []usersModels.UserInfoRepoModel
	for rows.Next() {
		var user usersModels.UserInfoRepoModel

		if err := rows.Scan(&user.Id, &user.Login, &user.Name); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}

		if user.Id.Valid {
			users = append(users, user)
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over users_groups: %v", err)
		return nil, err
	}

	log.Printf("Available users fetched successfully")
	return users, nil
}
