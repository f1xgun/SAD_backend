package groups

import (
	"log"
	groupsModels "sad/internal/models/groups"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(c *fiber.Ctx, group groupsModels.Group) error {
	query := "INSERT INTO groups (id, number) VALUES (@id, @number)"
	log.Printf("Creating group: %v", group)

	args := pgx.NamedArgs{
		"id":     group.Id,
		"number": group.Number,
	}
	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error creating group: %v, error: %v", group, err)
	} else {
		log.Printf("Group created successfully: %v", group)
	}

	return err
}

func (r *repository) GetById(c *fiber.Ctx, groupId string) (*groupsModels.Group, error) {
	query := "SELECT id, number FROM groups WHERE id=$1"
	log.Printf("Fetching group by id: %s", groupId)

	row := r.db.QueryRow(c.Context(), query, groupId)

	group := &groupsModels.Group{}
	err := row.Scan(&group.Id, &group.Number)
	if err == pgx.ErrNoRows {
		log.Printf("Group not found with id: %s", groupId)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching group by id: %s, error: %v", groupId, err)
		return nil, err
	}

	log.Printf("Group fetched successfully by id: %s", groupId)
	return group, nil
}

func (r *repository) GetAll(c *fiber.Ctx) ([]groupsModels.Group, error) {
	query := "SELECT id, number FROM groups"
	log.Printf("Fetching all groups")

	rows, err := r.db.Query(c.Context(), query)
	if err != nil {
		log.Printf("Error fetching all groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []groupsModels.Group
	for rows.Next() {
		var group groupsModels.Group
		if err := rows.Scan(&group.Id, &group.Number); err != nil {
			log.Printf("Error scanning group: %v", err)
			continue
		}
		groups = append(groups, group)
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

func (r *repository) GetGroupUsers(c *fiber.Ctx, groupId string) ([]string, error) {
	query := "SELECT user_id FROM users_groups WHERE group_id=@group_id"
	args := pgx.NamedArgs{
		"group_id": groupId,
	}

	rows, err := r.db.Query(c.Context(), query, args)
	if err != nil {
		log.Printf("Error fetching group's users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var usersId []string
	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		usersId = append(usersId, userId)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over users_groups: %v", err)
		return nil, err
	}

	log.Printf("All users id fetched successfully")
	return usersId, nil
}

func (r *repository) DeleteGroup(c *fiber.Ctx, groupId string) error {
	query := "DELETE FROM users_groups WHERE group_id=@group_id"
	args := pgx.NamedArgs{
		"group_id": groupId,
	}

	_, err := r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error deleting users from group with id %s , err: %v", groupId, err)
	} else {
		log.Printf("Delete users from group with id %s successfully", groupId)
	}

	query = "DELETE FROM groups WHERE id=@group_id"
	_, err = r.db.Exec(c.Context(), query, args)
	if err != nil {
		log.Printf("Error deleting group with id %s , err: %v", groupId, err)
	} else {
		log.Printf("Delete group with id %s successfully", groupId)
	}

	return err
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
