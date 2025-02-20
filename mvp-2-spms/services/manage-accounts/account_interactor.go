package managestudents

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/manage-accounts/outputdata"
	"mvp-2-spms/services/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/oauth2"
)

const pbkdf2Iterations int = 4096
const pbkdf2HashSize int = 32

type AccountInteractor struct {
	accountRepo interfaces.IAccountRepository
	uniRepo     interfaces.IUniversityRepository
}

func InitAccountInteractor(accRepo interfaces.IAccountRepository, uniRepo interfaces.IUniversityRepository) *AccountInteractor {
	return &AccountInteractor{
		accountRepo: accRepo,
		uniRepo:     uniRepo,
	}
}

func (a *AccountInteractor) GetAccountProfessorId(login string) (string, error) {
	res, err := a.accountRepo.GetAccountByLogin(login)
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func (a *AccountInteractor) GetProfessorInfo(input inputdata.GetProfessorInfo) (outputdata.GetProfessorInfo, error) {
	profInfo, err := a.accountRepo.GetProfessorById(fmt.Sprint(input.AccountId))
	if err != nil {
		return outputdata.GetProfessorInfo{}, err
	}

	uni, err := a.uniRepo.GetUniversityById(profInfo.UniversityId)
	if err != nil {
		if !errors.Is(err, models.ErrUniNoFound) {
			return outputdata.GetProfessorInfo{}, err
		}
		uni = entities.University{
			Id:   "-1",
			Name: "Университет неизвестен",
			City: "-",
		}
	}

	// add get account login
	output := outputdata.MapToGetAccountInfo(profInfo, uni)
	return output, nil
}

func (a *AccountInteractor) GetPlannerIntegration(input inputdata.GetPlannerIntegration) (outputdata.GetPlannerIntegration, error) {
	planner, err := a.accountRepo.GetAccountPlannerData(fmt.Sprint(input.AccountId))
	if err != nil {
		return outputdata.GetPlannerIntegration{}, err
	}

	output := outputdata.MapToGetPlannerIntegration(planner)
	return output, nil
}

func (a *AccountInteractor) GetDriveIntegration(input inputdata.GetDriveIntegration) (outputdata.GetDriveIntegration, error) {
	drive, err := a.accountRepo.GetAccountDriveData(fmt.Sprint(input.AccountId))
	if err != nil {
		return outputdata.GetDriveIntegration{}, err
	}

	output := outputdata.MapToGetDriveIntegration(drive)
	return output, nil
}

func (a *AccountInteractor) SetProfessorPlanner(plannerId, profId string) error {
	plannerInfo, err := a.accountRepo.GetAccountPlannerData(profId)
	if err != nil {
		return err
	}

	plannerInfo.PlannerData.Id = plannerId

	err = a.accountRepo.UpdateAccountPlannerIntegration(plannerInfo)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountInteractor) GetProfessorIntegrPlanners(profId string, planner interfaces.IPlannerService) (outputdata.GetProfessorIntegrPlanners, error) {
	plannerInfo, err := a.accountRepo.GetAccountPlannerData(profId)
	if err != nil {
		return outputdata.GetProfessorIntegrPlanners{}, err
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	// check for access token first????????????????????????????????????????????
	token := &oauth2.Token{
		RefreshToken: plannerInfo.ApiKey,
	}
	err = planner.Authentificate(token)
	if err != nil {
		return outputdata.GetProfessorIntegrPlanners{}, err
	}

	planners, err := planner.GetAllPlanners()
	if err != nil {
		return outputdata.GetProfessorIntegrPlanners{}, err
	}

	return outputdata.MapToGetProfessorIntegrPlanners(planners), nil
}

func (a *AccountInteractor) GetDriveBaseFolderName(folderId, profId string, cloudDrive interfaces.ICloudDrive) (string, error) {
	driveInfo, err := a.accountRepo.GetAccountDriveData(fmt.Sprint(profId))
	if err != nil {
		return "", err
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	// check for access token first????????????????????????????????????????????
	token := &oauth2.Token{
		RefreshToken: driveInfo.ApiKey,
	}
	err = cloudDrive.Authentificate(token)
	if err != nil {
		return "", err
	}

	folderName, err := cloudDrive.GetFolderNameById(folderId)
	if err != nil {
		return "", err
	}

	return folderName, nil
}

func (a *AccountInteractor) GetRepoHubIntegration(input inputdata.GetRepoHubIntegration) (outputdata.GetRepoHubIntegration, error) {
	repoHub, err := a.accountRepo.GetAccountRepoHubData(fmt.Sprint(input.AccountId))
	if err != nil {
		return outputdata.GetRepoHubIntegration{}, err
	}

	output := outputdata.MapToGetRepoHubIntegration(repoHub)
	return output, nil
}

func (a *AccountInteractor) SetPlannerIntegration(input inputdata.SetPlannerIntegration, planner interfaces.IPlannerService) (outputdata.SetPlannerIntegration, error) {
	token, err := planner.GetToken(input.AuthCode)
	if err != nil {
		return outputdata.SetPlannerIntegration{}, err
	}

	refreshTok := token.RefreshToken
	accessTok := token.AccessToken
	expires := token.Expiry

	integr := models.PlannerIntegration{
		BaseIntegration: models.BaseIntegration{
			AccountId: fmt.Sprint(input.AccountId),
			ApiKey:    refreshTok,
			Type:      input.Type,
		},
		PlannerData: models.PlannerData{},
	}

	err = a.accountRepo.AddAccountPlannerIntegration(integr)
	if err != nil {
		return outputdata.SetPlannerIntegration{}, err
	}

	return outputdata.SetPlannerIntegration{
		AccessToken: accessTok,
		Expiry:      expires,
	}, nil
}

func (a *AccountInteractor) SetDriveIntegration(input inputdata.SetDriveIntegration, drive interfaces.ICloudDrive) (outputdata.SetDriveIntegration, error) {
	token, _ := drive.GetToken(input.AuthCode)
	refreshTok := token.RefreshToken
	accessTok := token.AccessToken
	expires := token.Expiry

	err := drive.Authentificate(token)
	if err != nil {
		return outputdata.SetDriveIntegration{}, err
	}

	baseFolder, err := drive.AddProfessorBaseFolder()
	if err != nil {
		return outputdata.SetDriveIntegration{}, err
	}

	integr := models.CloudDriveIntegration{
		BaseIntegration: models.BaseIntegration{
			AccountId: fmt.Sprint(input.AccountId),
			ApiKey:    refreshTok,
			Type:      input.Type,
		},
		DriveData: baseFolder,
	}

	err = a.accountRepo.AddAccountDriveIntegration(integr)
	if err != nil {
		return outputdata.SetDriveIntegration{}, err
	}

	return outputdata.SetDriveIntegration{
		AccessToken: accessTok,
		Expiry:      expires,
	}, nil
}

func (a *AccountInteractor) SetRepoHubIntegration(input inputdata.SetRepoHubIntegration, planner interfaces.IGitRepositoryHub) (outputdata.SetRepoHubIntegration, error) {
	token, err := planner.GetToken(input.AuthCode)
	if err != nil {
		return outputdata.SetRepoHubIntegration{}, err
	}

	refreshTok := token.RefreshToken
	accessTok := token.AccessToken
	expires := token.Expiry
	integr := models.BaseIntegration{
		AccountId: fmt.Sprint(input.AccountId),
		ApiKey:    refreshTok,
		Type:      input.Type,
	}

	err = a.accountRepo.AddAccountRepoHubIntegration(integr)
	if err != nil {
		return outputdata.SetRepoHubIntegration{}, err
	}

	return outputdata.SetRepoHubIntegration{
		AccessToken: accessTok,
		Expiry:      expires,
	}, nil
}

func (a *AccountInteractor) GetAccountIntegrations(input inputdata.GetAccountIntegrations) (outputdata.GetAccountIntegrations, error) {
	var (
		outputDrive   *outputdata.GetAccountIntegrationsDrive   = nil
		outputPlanner *outputdata.GetAccountIntegrationsPlanner = nil
		outputRepos   []outputdata.GetAccountIntegrationsIntegr = []outputdata.GetAccountIntegrationsIntegr{}
	)

	found := true
	drive, err := a.accountRepo.GetAccountDriveData(fmt.Sprint(input.AccountId))
	if err != nil {
		if !errors.Is(err, models.ErrAccountDriveDataNotFound) {
			return outputdata.GetAccountIntegrations{}, err
		}
		found = false
	}

	if found {
		outputDrive = &outputdata.GetAccountIntegrationsDrive{
			Type: outputdata.GetAccountIntegrationsIntegr{
				Id:   drive.Type,
				Name: drive.GetTypeAsString(),
			},
			BaseFolderId: drive.BaseFolderId,
		}
	}

	found = true
	planner, err := a.accountRepo.GetAccountPlannerData(fmt.Sprint(input.AccountId))
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			return outputdata.GetAccountIntegrations{}, err
		}
		found = false
	}

	if found {
		outputPlanner = &outputdata.GetAccountIntegrationsPlanner{
			Type: outputdata.GetAccountIntegrationsIntegr{
				Id:   planner.Type,
				Name: planner.GetTypeAsString(),
			},
			PlannerName: planner.PlannerData.Id, ///////////////////////////////////////change
		}
	}

	found = true
	repohub, err := a.accountRepo.GetAccountRepoHubData(fmt.Sprint(input.AccountId))
	if err != nil {
		if !errors.Is(err, models.ErrAccountRepoHubDataNotFound) {
			return outputdata.GetAccountIntegrations{}, err
		}
		found = false
	}

	if found {
		outputRepos = append(outputRepos, outputdata.GetAccountIntegrationsIntegr{
			Id:   repohub.Type,
			Name: repohub.GetRepoHubTypeAsString(),
		})
	}

	return outputdata.MapToGetAccountIntegrations(outputDrive, outputPlanner, outputRepos), nil
}

func (a *AccountInteractor) CheckCredsValidity(input inputdata.CheckCredsValidity) (bool, error) {
	account, err := a.accountRepo.GetAccountByLogin(input.Login)
	if err != nil {
		return false, err
	}

	key := pbkdf2.Key([]byte(input.Password), []byte(account.Salt), pbkdf2Iterations, pbkdf2HashSize, sha512.New)

	return bytes.Equal(key, account.Hash), nil
}

func (a *AccountInteractor) CheckUsernameExists(input inputdata.CheckUsernameExists) (bool, error) {
	_, err := a.accountRepo.GetAccountByLogin(input.Login)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *AccountInteractor) SignUp(input inputdata.SignUp) (outputdata.SignUp, error) {
	salt := uuid.NewString()
	passHash := pbkdf2.Key([]byte(input.Password), []byte(salt), pbkdf2Iterations, pbkdf2HashSize, sha512.New)

	prof := entities.Professor{
		Person: entities.Person{
			Name:       input.Name,
			Surname:    input.Surname,
			Middlename: input.Middlename,
		},
		ScienceDegree: input.ScienceDegree,
		UniversityId:  fmt.Sprint(input.UniId),
	}

	prof, err := a.accountRepo.AddProfessor(prof)
	if err != nil {
		return outputdata.SignUp{}, err
	}

	account := models.Account{
		Login: input.Login,
		Hash:  passHash,
		Salt:  salt,
		Id:    prof.Id,
	}

	err = a.accountRepo.AddAccount(account)
	if err != nil {
		return outputdata.SignUp{}, err
	}

	return outputdata.SignUp{
		Id:    account.Id,
		Login: account.Login,
	}, nil
}
