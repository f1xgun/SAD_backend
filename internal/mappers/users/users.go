package usersMapper

import (
	"sad/internal/models/users"
)

func UserInfoFromRepoToService(repoModel usersModels.UserInfoRepoModel) usersModels.UserInfo {
	return usersModels.UserInfo{
		Id:    repoModel.Id.String,
		Name:  repoModel.Name.String,
		Login: repoModel.Login.String,
	}
}
