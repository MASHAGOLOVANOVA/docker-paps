package accountrepository

import (
	"mvp-2-spms/database"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/models"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dsn = "root:root@tcp(127.0.0.1:3306)/student_project_management_testing?parseTime=true"

func connectDB() *database.Database {
	gdb, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	db := database.InitDatabade(gdb)
	return db
}

func TestAccountRepo_GetAccountByLogin(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent account", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		_, err := ar.GetAccountByLogin("123")

		// assert
		assert.ErrorIs(t, err, models.ErrAccountNotFound)
	})

	t.Run("ok, get existing account", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		login := "test"
		err := addTestingAccount(login, ar)
		assert.NoError(t, err)

		// act
		_, err = ar.GetAccountByLogin(login)

		// assert
		assert.NoError(t, err)

		// cleanup
		err = deleteTestingAccount(login, ar)
		assert.NoError(t, err)
	})
}
func TestAccountRepo_AddProfessor(t *testing.T) {
	db := connectDB()

	t.Run("ok", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		prof := domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "Test",
				Surname:    "1",
				Middlename: "2",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		}

		// act
		prof, err := ar.AddProfessor(prof)

		// assert
		assert.NoError(t, err)
		foundProf, err := ar.GetProfessorById(prof.Id)
		assert.NoError(t, err)
		assert.Equal(t, prof.Id, foundProf.Id)
		assert.Equal(t, prof.Name, foundProf.Name)
		assert.Equal(t, prof.Surname, foundProf.Surname)
		assert.Equal(t, prof.Middlename, foundProf.Middlename)
		assert.Equal(t, prof.ScienceDegree, foundProf.ScienceDegree)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)
	})

	t.Run("fail, uni id doesnt exist", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		prof := domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "Test",
				Surname:    "1",
				Middlename: "2",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
			UniversityId:  "2",
		}

		// act
		_, err := ar.AddProfessor(prof)

		// assert
		assert.Error(t, err)
	})
}

func TestAccountRepo_AddAccount(t *testing.T) {
	db := connectDB()

	t.Run("ok", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		prof := domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "Test",
				Surname:    "1",
				Middlename: "2",
			},
			ScienceDegree: "sd",
		}
		prof, err := ar.AddProfessor(prof)
		assert.NoError(t, err)

		login := time.Now().Format(time.RFC3339)
		acc := models.Account{
			Login: login,
			Hash:  []byte{5, 6, 2},
			Salt:  "123232434",
			Id:    prof.Id,
		}

		// act
		err = ar.AddAccount(acc)

		// assert
		assert.NoError(t, err)

		foundAcc, err := ar.GetAccountByLogin(login)
		assert.NoError(t, err)
		assert.Equal(t, acc.Id, foundAcc.Id)
		assert.Equal(t, acc.Login, foundAcc.Login)
		assert.Equal(t, acc.Hash, foundAcc.Hash)
		assert.Equal(t, acc.Salt, foundAcc.Salt)

		// cleanup
		err = deleteTestingAccount(login, ar)
		assert.NoError(t, err)
	})

	t.Run("fail, prof id doesnt exist", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		login := time.Now().Format(time.RFC3339)
		acc := models.Account{
			Login: login,
			Hash:  []byte{5, 6, 2},
			Salt:  "123232434",
			Id:    "0",
		}

		// act
		err := ar.AddAccount(acc)

		// assert
		assert.Error(t, err)
	})
}

func TestAccountRepo_GetProfessorById(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent prof", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		_, err := ar.GetProfessorById("123")

		// assert
		assert.ErrorIs(t, err, models.ErrProfessorNotFound)
	})

	t.Run("ok, get existing prof", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "",
				Surname:    "",
				Middlename: "",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		// act
		foundProf, err := ar.GetProfessorById(prof.Id)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, prof.Id, foundProf.Id)
		assert.Equal(t, prof.Name, foundProf.Name)
		assert.Equal(t, prof.Surname, foundProf.Surname)
		assert.Equal(t, prof.Middlename, foundProf.Middlename)
		assert.Equal(t, prof.ScienceDegree, foundProf.ScienceDegree)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)
	})
}

func TestAccountRepo_GetAccountPlannerData(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent planner data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		_, err := ar.GetAccountPlannerData("123")

		// assert
		assert.ErrorIs(t, err, models.ErrAccountPlannerDataNotFound)
	})

	t.Run("ok, get existing planner data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		planner := models.PlannerIntegration{
			BaseIntegration: models.BaseIntegration{
				AccountId: prof.Id,
				ApiKey:    "api",
				Type:      int(models.GoogleCalendar),
			},
			PlannerData: models.PlannerData{
				Id: time.Now().Format(time.RFC3339),
			},
		}

		err = ar.AddAccountPlannerIntegration(planner)
		assert.NoError(t, err)

		// act
		foundPl, err := ar.GetAccountPlannerData(planner.AccountId)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, planner.AccountId, foundPl.AccountId)
		assert.Equal(t, planner.ApiKey, foundPl.ApiKey)
		assert.Equal(t, planner.Id, foundPl.Id)
		assert.Equal(t, planner.Type, foundPl.Type)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)
	})
}

func TestAccountRepo_GetAccountDriveData(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent drive data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		_, err := ar.GetAccountDriveData("123")

		// assert
		assert.ErrorIs(t, err, models.ErrAccountDriveDataNotFound)
	})

	t.Run("ok, get existing drive data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		drive := models.CloudDriveIntegration{
			BaseIntegration: models.BaseIntegration{
				AccountId: prof.Id,
				ApiKey:    "api",
				Type:      int(models.GoogleDrive),
			},
			DriveData: models.DriveData{
				BaseFolderId: time.Now().Format(time.RFC3339),
			},
		}

		err = ar.AddAccountDriveIntegration(drive)
		assert.NoError(t, err)

		// act
		foundDr, err := ar.GetAccountDriveData(drive.AccountId)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, drive.AccountId, foundDr.AccountId)
		assert.Equal(t, drive.ApiKey, foundDr.ApiKey)
		assert.Equal(t, drive.BaseFolderId, foundDr.BaseFolderId)
		assert.Equal(t, drive.Type, foundDr.Type)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)
	})
}

func TestAccountRepo_GetAccountRepoHubData(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent repo hub data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		_, err := ar.GetAccountRepoHubData("123")

		// assert
		assert.ErrorIs(t, err, models.ErrAccountRepoHubDataNotFound)
	})

	t.Run("ok, get existing repo hub data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		prof, err := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err)

		repo := models.BaseIntegration{
			AccountId: prof.Id,
			ApiKey:    time.Now().Format(time.RFC3339),
			Type:      int(models.GoogleDrive),
		}

		err = ar.AddAccountRepoHubIntegration(repo)
		assert.NoError(t, err)

		// act
		foundRepo, err := ar.GetAccountRepoHubData(repo.AccountId)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, repo.AccountId, foundRepo.AccountId)
		assert.Equal(t, repo.ApiKey, foundRepo.ApiKey)
		assert.Equal(t, repo.Type, foundRepo.Type)

		// cleanup
		profId, err := strconv.Atoi(prof.Id)
		assert.NoError(t, err)
		err = ar.DeleteProfessor(profId)
		assert.NoError(t, err)
	})
}

func addTestingAccount(name string, ar *AccountRepository) error {
	prof, err := ar.AddProfessor(domainaggregate.Professor{
		Person: domainaggregate.Person{
			Name:       "",
			Surname:    "",
			Middlename: "",
		},
		ScienceDegree: "",
	})
	if err != nil {
		return err
	}
	err = ar.AddAccount(models.Account{
		Login: name,
		Hash:  []byte{},
		Salt:  "",
		Id:    prof.Id,
	})
	if err != nil {
		return err
	}
	return nil
}

func deleteTestingAccount(name string, ar *AccountRepository) error {
	err := ar.DeleteAccountByLogin(name)
	if err != nil {
		return err
	}
	return nil
}
