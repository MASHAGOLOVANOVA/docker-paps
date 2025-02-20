package interfaces

import (
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-meetings/inputdata"
	"mvp-2-spms/services/manage-meetings/outputdata"
)

type IMeetingInteractor interface {
	AddMeeting(input inputdata.AddMeeting, planner interfaces.IPlannerService) (outputdata.AddMeeting, error)
	GetProfessorMeetings(input inputdata.GetProfessorMeetings, planner interfaces.IPlannerService) (outputdata.GetProfesorMeetings, error)
}
