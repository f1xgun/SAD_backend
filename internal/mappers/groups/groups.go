package groupsMapper

import (
	usersMapper "sad/internal/mappers/users"
	groupsModels "sad/internal/models/groups"
	usersModels "sad/internal/models/users"
)

func FromGroupWithUsersRepoModelToEntity(repoModel groupsModels.GroupWithUsersRepo) groupsModels.GroupWithUsers {
	users := make([]usersModels.UserInfo, 0)
	for _, userRepo := range repoModel.Users {
		user := usersMapper.UserInfoFromRepoToService(userRepo)
		users = append(users, user)
	}

	return groupsModels.GroupWithUsers{
		Group: groupsModels.Group{
			Id:     repoModel.Id.String,
			Number: repoModel.Number.String,
		},
		Users: users,
	}
}

func FromGroupRepoModelToEntity(repoModel groupsModels.GroupRepoModel) groupsModels.Group {
	return groupsModels.Group{
		Id:     repoModel.Id.String,
		Number: repoModel.Number.String,
	}
}

func FromGroupsRepoModelToEntity(repoModel []groupsModels.GroupRepoModel) []groupsModels.Group {
	groups := make([]groupsModels.Group, 0)
	for _, groupRepo := range repoModel {
		group := FromGroupRepoModelToEntity(groupRepo)
		groups = append(groups, group)
	}
	return groups
}
