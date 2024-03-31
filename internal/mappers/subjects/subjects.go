package subjectsMappers

import (
	subjectsModels "sad/internal/models/subjects"
)

func FromSubjectRepoModelToEntity(repoModel subjectsModels.SubjectRepoModel) subjectsModels.Subject {
	return subjectsModels.Subject{
		Id:   repoModel.Id.String,
		Name: repoModel.Name.String,
	}
}

func FromSubjectsRepoModelToEntity(repoModel []subjectsModels.SubjectRepoModel) []subjectsModels.Subject {
	groups := make([]subjectsModels.Subject, 0)
	for _, groupRepo := range repoModel {
		group := FromSubjectRepoModelToEntity(groupRepo)
		groups = append(groups, group)
	}
	return groups
}
