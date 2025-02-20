package accountrepository

import (
	"errors"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
	"strconv"

	"gorm.io/gorm"
)

type AccountRepository struct {
	dbContext database.Database
}

func InitAccountRepository(dbcxt database.Database) *AccountRepository {
	return &AccountRepository{
		dbContext: dbcxt,
	}
}

func (r *AccountRepository) GetAccountByLogin(login string) (usecasemodels.Account, error) {
	acc := models.Account{}

	result := r.dbContext.DB.Select("*").Where("login = ?", login).Take(&acc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.Account{}, usecasemodels.ErrAccountNotFound
		}
		return usecasemodels.Account{}, result.Error
	}

	return acc.MapToUseCaseModel(), nil
}

func (r *AccountRepository) DeleteAccountByLogin(login string) error {
	acc, err := r.GetAccountByLogin(login)
	if err != nil {
		return err
	}

	profId, err := strconv.Atoi(acc.Id)
	if err != nil {
		return err
	}

	result := r.dbContext.DB.Delete(&models.Professor{
		Id: uint(profId),
	})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.ErrProfessorNotFound
		}
		return result.Error
	}
	return nil
}

func (r *AccountRepository) AddProfessor(prof entities.Professor) (entities.Professor, error) {
	dbProf := models.Professor{}
	dbProf.MapEntityToThis(prof)

	result := r.dbContext.DB.Create(&dbProf)
	if result.Error != nil {
		return entities.Professor{}, result.Error
	}

	return dbProf.MapToEntity(), nil
}

func (r *AccountRepository) DeleteProfessor(profId int) error {
	dbProf := models.Professor{Id: uint(profId)}

	result := r.dbContext.DB.Delete(&dbProf)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *AccountRepository) AddAccount(account usecasemodels.Account) error {
	dbAcc := models.Account{}
	dbAcc.MapUseCaseModelToThis(account)

	result := r.dbContext.DB.Create(&dbAcc)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *AccountRepository) GetProfessorById(id string) (entities.Professor, error) {
	prof := models.Professor{}

	result := r.dbContext.DB.Select("*").Where("id = ?", id).Take(&prof)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Professor{}, usecasemodels.ErrProfessorNotFound
		}
		return entities.Professor{}, result.Error
	}

	return prof.MapToEntity(), nil
}

func (r *AccountRepository) GetAccountPlannerData(id string) (usecasemodels.PlannerIntegration, error) {
	dbPlanner := models.PlannerIntegration{}

	result := r.dbContext.DB.Select("*").Where("account_id = ?", id).Take(&dbPlanner)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.PlannerIntegration{}, usecasemodels.ErrAccountPlannerDataNotFound
		}
		return usecasemodels.PlannerIntegration{}, result.Error
	}

	return dbPlanner.MapToUseCaseModel(), nil
}

func (r *AccountRepository) GetAccountDriveData(id string) (usecasemodels.CloudDriveIntegration, error) {
	dbDrive := models.DriveIntegration{}

	result := r.dbContext.DB.Select("*").Where("account_id = ?", id).Take(&dbDrive)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.CloudDriveIntegration{}, usecasemodels.ErrAccountDriveDataNotFound
		}
		return usecasemodels.CloudDriveIntegration{}, result.Error
	}

	return dbDrive.MapToUseCaseModel(), nil
}

// can return multiple for 1 account, should consider this
func (r *AccountRepository) GetAccountRepoHubData(id string) (usecasemodels.BaseIntegration, error) {
	dbRHub := models.GitRepositoryIntegration{}

	result := r.dbContext.DB.Select("*").Where("account_id = ?", id).Take(&dbRHub)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.BaseIntegration{}, usecasemodels.ErrAccountRepoHubDataNotFound
		}
		return usecasemodels.BaseIntegration{}, result.Error
	}

	return dbRHub.MapToUseCaseModel(), nil
}

// func (r *AccountRepository) DeleteAccountPlannerData(id int) error {
// 	dbPl := models.PlannerIntegration{AccountId: uint(id)}

// 	result := r.dbContext.DB.Where("account_id = ?", id).Delete(&dbPl)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r *AccountRepository) DeleteAccountDriveData(id int) error {
// 	dbDrive := models.DriveIntegration{AccountId: uint(id)}

// 	result := r.dbContext.DB.Where("account_id = ?", id).Delete(&dbDrive)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r *AccountRepository) DeleteAccountRepoHubData(id int) error {
// 	dbRepo := models.GitRepositoryIntegration{AccountId: uint(id)}

// 	result := r.dbContext.DB.Where("account_id = ?", id).Delete(&dbRepo)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

func (r *AccountRepository) AddAccountPlannerIntegration(integr usecasemodels.PlannerIntegration) error {
	dbPlanner := models.PlannerIntegration{}
	dbPlanner.MapUseCaseModelToThis(integr)

	result := r.dbContext.DB.Create(&dbPlanner)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (r *AccountRepository) AddAccountDriveIntegration(integr usecasemodels.CloudDriveIntegration) error {
	dbDrive := models.DriveIntegration{}
	dbDrive.MapUseCaseModelToThis(integr)

	result := r.dbContext.DB.Create(&dbDrive)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (r *AccountRepository) AddAccountRepoHubIntegration(integr usecasemodels.BaseIntegration) error {
	dbRepoHub := models.GitRepositoryIntegration{}
	dbRepoHub.MapUseCaseModelToThis(integr)

	result := r.dbContext.DB.Create(&dbRepoHub)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *AccountRepository) UpdateAccountPlannerIntegration(integr usecasemodels.PlannerIntegration) error {
	plannerDb := models.PlannerIntegration{}
	plannerDb.MapUseCaseModelToThis(integr)

	result := r.dbContext.DB.Where("account_id = ?", integr.AccountId).Save(&plannerDb)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *AccountRepository) UpdateAccountDriveIntegration(integr usecasemodels.CloudDriveIntegration) error {
	result := r.dbContext.DB.Model(&models.DriveIntegration{}).Where("account_id = ?", integr.AccountId).Update("api_key", integr.ApiKey)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return usecasemodels.ErrAccountDriveDataNotFound
	}

	return nil
}

func (r *AccountRepository) UpdateAccountRepoHubIntegration(integr usecasemodels.BaseIntegration) error {
	result := r.dbContext.DB.Model(&models.GitRepositoryIntegration{}).Where("account_id = ?", integr.AccountId).Update("api_key", integr.ApiKey)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return usecasemodels.ErrAccountDriveDataNotFound
	}

	return nil
}
