package meetingrepository

import (
	"errors"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
	"time"

	"gorm.io/gorm"
)

type MeetingRepository struct {
	dbContext database.Database
}

func InitMeetingRepository(dbcxt database.Database) *MeetingRepository {
	return &MeetingRepository{
		dbContext: dbcxt,
	}
}

func (r *MeetingRepository) CreateMeeting(meeting entities.Meeting) (entities.Meeting, error) {
	dbmeeting := models.Meeting{}
	dbmeeting.MapEntityToThis(meeting)
	result := r.dbContext.DB.Create(&dbmeeting)
	if result.Error != nil {
		return entities.Meeting{}, result.Error
	}
	return dbmeeting.MapToEntity(), nil
}

func (r *MeetingRepository) DeleteMeeting(id int) error {
	dbmeeting := models.Meeting{Id: uint(id)}

	result := r.dbContext.DB.Delete(&dbmeeting)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MeetingRepository) AssignPlannerMeeting(meeting usecasemodels.PlannerMeeting) error {
	err := r.dbContext.DB.Transaction(func(tx *gorm.DB) error {
		result := r.dbContext.DB.Model(&models.Meeting{}).Where("id = ?", meeting.Meeting.Id).Update("planner_id", meeting.MeetingPlannerId)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return usecasemodels.ErrMeetingNotFound
		}

		return nil
	})

	return err
}

func (r *MeetingRepository) GetProfessorMeetings(profId string, from time.Time, to time.Time) ([]entities.Meeting, error) {
	var meetingsDb []models.Meeting

	query := r.dbContext.DB.Select("*")
	if to.IsZero() {
		query = query.Where("professor_id = ? AND meeting_time > ?", profId, from)
	} else {
		query = query.Where("professor_id = ? AND meeting_time > ? AND meeting_time < ?", profId, from, to)
	}

	result := query.Order("meeting_time asc").Find(&meetingsDb)
	if result.Error != nil {
		return []entities.Meeting{}, result.Error
	}

	meetings := []entities.Meeting{}
	for _, m := range meetingsDb {
		meetings = append(meetings, m.MapToEntity())
	}
	return meetings, nil
}

func (r *MeetingRepository) GetMeetingById(meetId string) (entities.Meeting, error) {
	var dbMeeting models.Meeting

	result := r.dbContext.DB.Select("*").Where("id = ?", meetId).Take(&dbMeeting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Meeting{}, usecasemodels.ErrMeetingNotFound
		}
		return entities.Meeting{}, result.Error
	}

	return dbMeeting.MapToEntity(), nil
}

func (r *MeetingRepository) GetMeetingPlannerId(meetId string) (string, error) {
	meeting := models.Meeting{}

	result := r.dbContext.DB.Select("planner_id").Where("id = ?", meetId).Take(&meeting)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", usecasemodels.ErrMeetingNotFound
		}
		return "", result.Error
	}

	return meeting.PlannerId.String, nil
}
