package grades

import (
	usersMapper "sad/internal/mappers/users"
	gradesModels "sad/internal/models/grades"
)

func FromGradeRepoModelToEntity(repoModel gradesModels.GradeRepoModel) gradesModels.Grade {
	isFinal := false
	if repoModel.IsFinal.Valid {
		isFinal = repoModel.IsFinal.Bool
	}

	var comment string
	if repoModel.Comment.Valid {
		comment = repoModel.Comment.String
	}

	return gradesModels.Grade{
		Id:         repoModel.Id.String,
		StudentId:  repoModel.StudentId.String,
		SubjectId:  repoModel.SubjectId.String,
		CreatedAt:  repoModel.CreatedAt.Time,
		Evaluation: int(repoModel.Evaluation.Int16),
		TeacherId:  repoModel.TeacherId.String,
		IsFinal:    &isFinal,
		Comment:    &comment,
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
	isFinal := false
	if repoModel.IsFinal.Valid {
		isFinal = repoModel.IsFinal.Bool
	}

	var comment string
	if repoModel.Comment.Valid {
		comment = repoModel.Comment.String
	}

	return gradesModels.GradeInfo{
		Id:          repoModel.Id.String,
		SubjectName: repoModel.SubjectName.String,
		TeacherName: repoModel.TeacherName.String,
		Evaluation:  int(repoModel.Evaluation.Int16),
		CreatedAt:   repoModel.CreatedAt.Time,
		IsFinal:     &isFinal,
		Comment:     &comment,
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

func FromUserWithGradesRepoModelToEntity(repoModel []gradesModels.UserSubjectGradesRepoModel) []gradesModels.UserSubjectGrades {
	usersWithGrades := make([]gradesModels.UserSubjectGrades, 0)
	for _, userWithGradesRepo := range repoModel {
		userWithGrades := gradesModels.UserSubjectGrades{
			Student: usersMapper.UserInfoFromRepoToService(userWithGradesRepo.Student),
			Grades:  FromGradesInfoRepoModelToEntity(userWithGradesRepo.Grades),
		}
		usersWithGrades = append(usersWithGrades, userWithGrades)
	}

	return usersWithGrades
}
