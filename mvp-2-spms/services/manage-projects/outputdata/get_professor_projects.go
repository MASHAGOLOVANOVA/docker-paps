package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetProfessorProjects struct {
	Projects []getProfProjProjectData `json:"projects"`
}

func MapToGetProfessorProjects(projectEntities []GetProfessorProjectsEntities) GetProfessorProjects {
	outputProjects := []getProfProjProjectData{}
	for _, projectEntitiy := range projectEntities {
		id, _ := strconv.Atoi(projectEntitiy.Project.Id)
		outputProjects = append(outputProjects,
			getProfProjProjectData{
				Id:          id,
				Theme:       projectEntitiy.Project.Theme,
				Status:      projectEntitiy.Project.Status.String(),
				Stage:       projectEntitiy.Project.Stage.String(),
				Year:        int(projectEntitiy.Project.Year),
				StudentName: projectEntitiy.Student.FullNameToString(),
				Cource:      int(projectEntitiy.Student.Cource),
			})
	}
	return GetProfessorProjects{
		Projects: outputProjects,
	}
}

type GetProfessorProjectsEntities struct {
	Project entities.Project
	Student entities.Student
}

type getProfProjProjectData struct {
	Id          int    `json:"id"`
	Theme       string `json:"theme"`
	StudentName string `json:"student_name"`
	Cource      int    `json:"cource"`
	Status      string `json:"status"`
	Stage       string `json:"stage"`
	Year        int    `json:"year"`
}
