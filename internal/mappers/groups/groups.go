package groupsMapper

import (
	subjectsMappers "sad/internal/mappers/subjects"
	"sad/internal/mappers/users"
	"sad/internal/models/groups"
	subjectsModels "sad/internal/models/subjects"
	"sad/internal/models/users"
)

func FromGroupDetailsRepoModelToEntity(repoModel groupsModels.GroupDetailsRepo) groupsModels.GroupDetails {
	users := make([]usersModels.UserInfo, 0)
	for _, userRepo := range repoModel.Users {
		user := usersMapper.UserInfoFromRepoToService(userRepo)
		users = append(users, user)
	}

	return groupsModels.GroupDetails{
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

func GroupsWithSubjectsFromRepoModelToEntity(repoModel []subjectsModels.GroupsWithSubjectsRepoModel) []subjectsModels.GroupsWithSubjects {
	groupsWithSubjects := make([]subjectsModels.GroupsWithSubjects, 0)
	for _, groupWithSubjectRepo := range repoModel {
		groupWithSubject := subjectsModels.GroupsWithSubjects{
			Group:    FromGroupRepoModelToEntity(groupWithSubjectRepo.Group),
			Subjects: subjectsMappers.FromSubjectsRepoModelToEntity(groupWithSubjectRepo.Subjects),
		}

		groupsWithSubjects = append(groupsWithSubjects, groupWithSubject)
	}
	return groupsWithSubjects
}
