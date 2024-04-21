package grades

import (
	gradesModels "sad/internal/models/grades"
)

func FromGradeRepoModelToEntity(repoModel gradesModels.GradeRepoModel) gradesModels.Grade {
	return gradesModels.Grade{
		Id:         repoModel.Id.String,
		StudentId:  repoModel.StudentId.String,
		SubjectId:  repoModel.SubjectId.String,
		CreatedAt:  repoModel.CreatedAt.Time,
		Evaluation: int(repoModel.Evaluation.Int16),
		TeacherId:  repoModel.TeacherId.String,
	}
}

func FromGradesRepoModelToEntity(repoModel []gradesModels.GradeRepoModel) []gradesModels.Grade {
	grades := make([]gradesModels.Grade, 0)
	for _, gradeRepo := range repoModel {
		grade := FromGradeRepoModelToEntity(gradeRepo)
		grades = append(grades, grade)
	}
	return grades
}

func FromGradeInfoRepoModelToEntity(repoModel gradesModels.GradeInfoRepoModel) gradesModels.GradeInfo {
	return gradesModels.GradeInfo{
		Id:          repoModel.Id.String,
		SubjectName: repoModel.SubjectName.String,
		TeacherName: repoModel.TeacherName.String,
		Evaluation:  int(repoModel.Evaluation.Int16),
		CreatedAt:   repoModel.CreatedAt.Time,
	}
}

func FromGradesInfoRepoModelToEntity(repoModel []gradesModels.GradeInfoRepoModel) []gradesModels.GradeInfo {
	grades := make([]gradesModels.GradeInfo, 0)
	for _, gradeInfoRepo := range repoModel {
		grade := FromGradeInfoRepoModelToEntity(gradeInfoRepo)
		grades = append(grades, grade)
	}
	return grades
}
