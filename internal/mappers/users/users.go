package usersMapper

import (
	"sad/internal/models/users"
)

func UserInfoFromRepoToService(repoModel usersModels.UserInfoRepoModel) usersModels.UserInfo {
	var role usersModels.UserRole

	switch repoModel.Role.String {
	case "student":
		role = usersModels.Student
	case "teacher":
		role = usersModels.Teacher
	case "admin":
		role = usersModels.Admin
	default:
		role = ""
	}

	var middleName string
	if repoModel.MiddleName.Valid {
		middleName = repoModel.MiddleName.String
	}

	return usersModels.UserInfo{
		Id:         repoModel.Id.String,
		Name:       repoModel.Name.String,
		Login:      repoModel.Login.String,
		Role:       role,
		LastName:   repoModel.LastName.String,
		MiddleName: &middleName,
	}
}

func UsersInfoFromRepoToService(repoModel []usersModels.UserInfoRepoModel) []usersModels.UserInfo {
	users := make([]usersModels.UserInfo, 0)
	for _, userRepo := range repoModel {
		user := UserInfoFromRepoToService(userRepo)
		users = append(users, user)
	}
	return users
}
