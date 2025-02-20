package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetProjectById struct {
	Id              int                       `json:"id"`
	Theme           string                    `json:"theme"`
	Student         GetProjectByIdStudentData `json:"student"`
	Status          string                    `json:"status"`
	Stage           string                    `json:"stage"`
	Year            int                       `json:"year"`
	CloudFolderLink string                    `json:"cloud_folder_link"`
}

func MapToGetProjectsById(project entities.Project, student entities.Student, edProgramme entities.EducationalProgramme, folderLink string) GetProjectById {
	pId, _ := strconv.Atoi(project.Id)
	sId, _ := strconv.Atoi(student.Id)
	return GetProjectById{
		Id:    pId,
		Theme: project.Theme,
		Student: GetProjectByIdStudentData{
			Id:          sId,
			Name:        student.Name,
			Surname:     student.Surname,
			Middlename:  student.Middlename,
			Cource:      int(student.Cource),
			EdProgramme: edProgramme.Name,
		},
		Status:          project.Status.String(),
		Stage:           project.Stage.String(),
		Year:            int(project.Year),
		CloudFolderLink: folderLink,
	}
}

type GetProjectByIdStudentData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Middlename  string `json:"middlename"`
	Cource      int    `json:"cource"`
	EdProgramme string `json:"education_programme"`
}
