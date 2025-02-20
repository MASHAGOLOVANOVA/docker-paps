package inputdata

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/models"
	"time"
)

type GetProfessorProjects struct {
	ProfessorId  uint
	FilterStatus *int
}

type GetProjectCommits struct {
	ProfessorId uint
	ProjectId   uint
	From        time.Time
}

type GetProjectById struct {
	ProfessorId uint
	ProjectId   uint
}

type GetProjectSupReport struct {
	ProfessorId uint
	ProjectId   uint
	Comment     string `json:"comment"`
}

type GetProjectStatsById struct {
	ProfessorId uint
	ProjectId   uint
}

type AddProject struct {
	ProfessorId         uint
	Theme               string
	StudentId           uint
	Year                uint
	RepositoryOwnerName string
	RepositoryName      string
}

type UpdateProject struct {
	Id                  int
	ProfessorId         *int
	Theme               *string
	StudentId           *int
	Year                *int
	RepositoryOwnerName *string
	RepositoryName      *string
	Status              *int
	Stage               *int
}

func (as UpdateProject) UpdateProjectEntity(p *entities.Project) error {
	if as.ProfessorId != nil {
		p.SupervisorId = fmt.Sprint(*as.ProfessorId)
	}
	if as.Stage != nil {
		p.Stage = entities.ProjectStage(*as.Stage)
	}
	if as.Status != nil {
		p.Status = entities.ProjectStatus(*as.Status)
	}
	if as.StudentId != nil {
		p.StudentId = fmt.Sprint(*as.StudentId)
	}
	if as.Theme != nil {
		p.Theme = *as.Theme
	}
	if as.Year != nil {
		p.Year = uint(*as.Year)
	}
	return nil
}

func (as AddProject) MapToProjectEntity() entities.Project {
	return entities.Project{
		Theme:        as.Theme,
		SupervisorId: fmt.Sprint(as.ProfessorId),
		StudentId:    fmt.Sprint(as.StudentId),
		Year:         as.Year,
		Stage:        entities.ProjectStage(entities.Analysis),
		Status:       entities.ProjectStatus(entities.ProjectInProgress),
	}
}

func (as AddProject) MapToRepositoryEntity() models.Repository {
	return models.Repository{
		RepoId:    as.RepositoryName,
		OwnerName: as.RepositoryOwnerName,
	}
}

type UpdateProjectGrading struct {
	ProjId           int
	ProfessorId      *int
	DefenctGrade     *float32 `json:"defence_grade,omitempty"`
	SupervisorReview *SupRew  `json:"supervisor_review,omitempty"`
}

type SupRew struct {
	Id           *int      `json:"id,omitempty"`
	Criterias    *[]Crit   `json:"criterias,omitempty"`
	CreationDate time.Time `json:"created"`
}

type Crit struct {
	Criteria string   `json:"criteria"`
	Grade    *float32 `json:"grade,omitempty"`
	Weight   float32  `json:"weight"`
}

func (as UpdateProjectGrading) MapToProjectGrading() entities.ProjectGrading {
	result := entities.ProjectGrading{}
	if as.DefenctGrade != nil {
		result.DefenceGrade = *as.DefenctGrade
	}
	if as.SupervisorReview != nil {
		result.SupervisorReview = entities.SupervisorReview{}
		if as.SupervisorReview.Id != nil {
			result.SupervisorReview.Id = uint(*as.SupervisorReview.Id)
		}
		if as.SupervisorReview.Criterias != nil {
			result.SupervisorReview.Criterias = []entities.Criteria{}
			for _, c := range *as.SupervisorReview.Criterias {
				cr := entities.Criteria{
					Description: c.Criteria,
					Weight:      c.Weight,
				}
				if c.Grade != nil {
					cr.Grade = *c.Grade
				}
				result.SupervisorReview.Criterias = append(result.SupervisorReview.Criterias, cr)
			}
		}
		result.SupervisorReview.CreationDate = as.SupervisorReview.CreationDate
	}
	return result
}
