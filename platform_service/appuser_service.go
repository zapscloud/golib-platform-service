package platform_service

import (
	"fmt"
	"log"
	"strings"

	"github.com/rs/xid"
	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// AppUserService - Users Service structure
type AppUserService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(userId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(userId string, indata utils.Map) (utils.Map, error)
	Delete(userId string, deletePermanent bool) error
	Authenticate(auth_key string, auth_user string, auth_pwd string) (utils.Map, error)
	ChangePassword(userId string, newpwd string) (utils.Map, error)

	BusinessList(userId string, filter string, sort string, skip int64, limit int64) (utils.Map, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type appUserBaseService struct {
	db_utils.DatabaseService
	daoAppUser  platform_repository.AppUserDao
	daoBusiness platform_repository.BusinessDao
	child       AppUserService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewAppUserService(props utils.Map) (AppUserService, error) {
	p := appUserBaseService{}
	log.Printf("NewAppUserService :: Start")

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewAppUserService ", err)
		return nil, err
	}

	p.daoAppUser = platform_repository.NewAppUserDao(p.GetClient())
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.child = &p

	return &p, nil
}

func (p *appUserBaseService) EndService() {
	log.Printf("EndappUserBaseService ")
	p.CloseDatabaseService()
}

// List - List All records
func (p *appUserBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AppUserService::FindAll - Begin")

	dataresponse, err := p.daoAppUser.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("AppUserService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *appUserBaseService) Get(userId string) (utils.Map, error) {
	log.Printf("AppUserService::GetDetails::  Begin %v", userId)

	data, err := p.daoAppUser.Get(userId)

	log.Println("AppUserService::GetDetails:: End ", data, err)
	return data, err
}

func (p *appUserBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("AppUserService::GetDetails::  Begin ", filter)

	data, err := p.daoAppUser.Find(filter)

	log.Println("AppUserService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *appUserBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var appUserId string

	dataval, dataok := indata[platform_common.FLD_APP_USER_ID]
	if dataok {
		appUserId = strings.ToLower(dataval.(string))
	} else {
		guid := xid.New()
		prefix := "user_"
		log.Println("Unique Profile ID", prefix, guid.String())
		appUserId = prefix + guid.String()
	}
	log.Println("Provided Profile ID:", appUserId)

	// Update converted/generated id back to indata
	indata[platform_common.FLD_APP_USER_ID] = appUserId

	dataval, dataok = indata[platform_common.FLD_APP_USER_PASSWORD]
	if dataok {
		indata[platform_common.FLD_APP_USER_PASSWORD] = utils.SHA(dataval.(string))
	}

	dataCreated, err := p.daoAppUser.Create(indata)
	if err != nil {
		return dataCreated, err
	}
	log.Println("UserService::Create - End ")
	return dataCreated, nil
}

// Update - Update Service
func (p *appUserBaseService) Update(userId string, indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Update - Begin")

	// Delete the Key fields
	delete(indata, platform_common.FLD_APP_USER_ID)

	// Check whether the password is sent
	if dataVal, dataOk := indata[platform_common.FLD_APP_USER_PASSWORD]; dataOk {
		indata[platform_common.FLD_APP_USER_PASSWORD] = utils.SHA(dataVal.(string))
	}

	data, err := p.daoAppUser.Update(userId, indata)

	log.Println("UserService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *appUserBaseService) Delete(userId string, deletePermanent bool) error {

	log.Println("UserService::Delete - Begin", userId)

	if deletePermanent {
		result, err := p.daoAppUser.Delete(userId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}

		data, err := p.Update(userId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("UserService::Delete - End")
	return nil
}

// GetDetails - Find By Code
func (p *appUserBaseService) Authenticate(auth_key string, auth_login string, auth_pwd string) (utils.Map, error) {
	log.Println("Authenticate::  Begin ", auth_key, auth_login, auth_pwd)

	log.Println("User Password from API", auth_pwd)
	encpwd := utils.SHA(auth_pwd)
	dataUser, err := p.daoAppUser.Authenticate(auth_key, auth_login, encpwd)

	log.Println("Length of dataUser :", dataUser)

	if err != nil {
		err := &utils.AppError{ErrorCode: "S30340101", ErrorMsg: "Wrong Credentials", ErrorDetail: "Authenticate credentials is wrong !!"}
		return utils.Map{}, err
	}

	isSuspended, err := utils.GetMemberDataBool(dataUser, db_common.FLD_IS_SUSPENDED)
	if err == nil && isSuspended {
		err := &utils.AppError{ErrorCode: "S30340102", ErrorMsg: "User is in suspended mode. Contact Admin!", ErrorDetail: "User not in Active Mode. Contact Admin!"}
		return utils.Map{}, err
	}

	// isVerified, err := utils.GetMemberDataBool(dataUser, db_common.FLD_IS_VERIFIED)
	// if err == nil && !isVerified {
	// 	err := &utils.AppError{ErrorCode: "S30340103", ErrorMsg: "User not yet verified!", ErrorDetail: "User not yet verified!!"}
	// 	return utils.Map{}, err
	// }

	return dataUser, nil
}

// Update - Update Service
func (p *appUserBaseService) ChangePassword(userId string, newpwd string) (utils.Map, error) {

	log.Println("AppUserService::ChangePassword - Begin")
	indata := utils.Map{
		platform_common.FLD_APP_USER_PASSWORD: utils.SHA(newpwd),
	}
	data, err := p.daoAppUser.Update(userId, indata)

	log.Println("AppUserService::ChangePassword - End ")
	return data, err
}

func (p *appUserBaseService) BusinessList(userId string, filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	_, err := p.daoAppUser.Get(userId)
	if err != nil {
		return utils.Map{}, err
	}

	dataBusiness, err := p.daoAppUser.BusinessList(userId, filter, sort, skip, limit)
	if err != nil {
		return utils.Map{}, err
	}

	for _, value := range dataBusiness[db_common.LIST_RESULT].([]utils.Map) {
		log.Println("AppUserService: BusinessList ", value)
		dataBusiness, err := p.daoBusiness.Get(value[platform_common.FLD_BUSINESS_ID].(string))
		if err != nil {
			continue
		}
		value["app_business"] = dataBusiness
	}
	log.Println("SysBusinessService::BusinessList ", dataBusiness)

	log.Println("SysBusinessService::BusinessList:: End ", dataBusiness, err)
	return dataBusiness, err
}
