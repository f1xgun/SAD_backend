package subjectsMappers

import (
	"sad/internal/models/subjects"
)

func FromSubjectRepoModelToEntity(repoModel subjectsModels.SubjectRepoModel) subjectsModels.Subject {
	return subjectsModels.Subject{
		Id:   repoModel.Id.String,
		Name: repoModel.Name.String,
	}
}

func FromSubjectsRepoModelToEntity(repoModel []subjectsModels.SubjectRepoModel) []subjectsModels.Subject {
	subjects := make([]subjectsModels.Subject, 0)
	for _, subjectRepo := range repoModel {
		subject := FromSubjectRepoModelToEntity(subjectRepo)
		subjects = append(subjects, subject)
	}
	return subjects
}
