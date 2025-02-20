package interfaces

import (
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/models"
	"time"
)

type IPlannerService interface {
	IIntegration
	AddMeeting(meeting entities.Meeting, plannerInfo models.PlannerIntegration) (models.PlannerMeeting, error)
	FindMeetingById(meetId string, plannerInfo models.PlannerIntegration) (bool, error)
	GetScheduleMeetingIds(from time.Time, plannerInfo models.PlannerIntegration) ([]string, error)
	GetAllPlanners() ([]models.PlannerData, error)
}
