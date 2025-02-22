package usersMapper

import (
	"sad/internal/models/users"
)

func UserInfoFromRepoToService(repoModel users.UserInfoRepoModel) users.UserInfo {
	var role users.UserRole

	switch repoModel.Role.String {
	case "student":
		role = users.Student
	case "teacher":
		role = users.Teacher
	case "admin":
		role = users.Admin
	default:
		role = ""
	}

	var middleName string
	if repoModel.MiddleName.Valid {
		middleName = repoModel.MiddleName.String
	}

	return users.UserInfo{
		Id:         repoModel.Id.String,
		Name:       repoModel.Name.String,
		Login:      repoModel.Login.String,
		Role:       role,
		LastName:   repoModel.LastName.String,
		MiddleName: &middleName,
	}
}

func UsersInfoFromRepoToService(repoModel []users.UserInfoRepoModel) []users.UserInfo {
	users := make([]users.UserInfo, 0)
	for _, userRepo := range repoModel {
		user := UserInfoFromRepoToService(userRepo)
		users = append(users, user)
	}
	return users
}
