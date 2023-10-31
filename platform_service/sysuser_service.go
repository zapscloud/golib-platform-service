package platform_services

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

// SysUserService - Users Service structure
type SysUserService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(userID string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(userID string, indata utils.Map) (utils.Map, error)
	Delete(userID string, delete_permanent bool) error
	Authenticate(auth_key string, auth_login string, auth_pwd string) (utils.Map, error)
	ChangePassword(userid string, newpwd string) (utils.Map, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type sysUserBaseService struct {
	db_utils.DatabaseService
	daoSysUser platform_repository.SysUserDao
	child      SysUserService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewSysUserService(props utils.Map) (SysUserService, error) {
	p := sysUserBaseService{}

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewSysUserService ", err)
		return nil, err
	}

	log.Printf("sysUserMongoService ")
	p.daoSysUser = platform_repository.NewSysUserDao(p.GetClient())

	p.child = &p

	return &p, nil
}

func (p *sysUserBaseService) EndService() {
	log.Printf("EndsysUserService ")
	p.CloseDatabaseService()
}

// List - List All records
func (p *sysUserBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("SysUserService::FindAll - Begin")

	dataresponse, err := p.daoSysUser.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("SysUserService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *sysUserBaseService) Get(userID string) (utils.Map, error) {
	log.Printf("SysUserService::GetDetails::  Begin %v", userID)

	data, err := p.daoSysUser.Get(userID)

	log.Println("SysUserService::GetDetails:: End ", data, err)
	return data, err
}

func (p *sysUserBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("SysUserService::GetDetails::  Begin ", filter)

	data, err := p.daoSysUser.Find(filter)

	log.Println("SysUserService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *sysUserBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var sysUserId string

	dataval, dataok := indata[platform_common.FLD_SYS_USER_ID]
	if dataok {
		sysUserId = strings.ToLower(dataval.(string))
	} else {
		guid := xid.New()
		prefix := "syusr_"
		log.Println("Unique Profile ID", prefix, guid.String())
		sysUserId = prefix + guid.String()
	}
	log.Println("Provided Profile ID:", sysUserId)

	// Update converted/generated id back to indata
	indata[platform_common.FLD_SYS_USER_ID] = sysUserId

	dataval, dataok = indata[platform_common.FLD_SYS_USER_PASSWORD]
	if dataok {
		indata[platform_common.FLD_SYS_USER_PASSWORD] = utils.SHA(dataval.(string))
	}

	dataCreated, err := p.daoSysUser.Create(indata)
	if err != nil {
		return dataCreated, err
	}
	log.Println("UserService::Create - End ")
	return dataCreated, nil
}

// Update - Update Service
func (p *sysUserBaseService) Update(userID string, indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Update - Begin")

	// Delete the Key fields
	delete(indata, platform_common.FLD_SYS_USER_ID)

	data, err := p.daoSysUser.Update(userID, indata)

	log.Println("UserService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *sysUserBaseService) Delete(userID string, delete_permanent bool) error {

	log.Println("UserService::Delete - Begin", userID)

	if delete_permanent {
		result, err := p.daoSysUser.Delete(userID)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(userID, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("UserService::Delete - End")
	return nil
}

// GetDetails - Find By Code
func (p *sysUserBaseService) Authenticate(auth_key string, auth_login string, auth_pwd string) (utils.Map, error) {
	log.Println("Authenticate::  Begin ", auth_key, auth_login, auth_pwd)

	log.Println("User Password from API", auth_pwd)
	encpwd := utils.SHA(auth_pwd)
	dataUser, err := p.daoSysUser.Authenticate(auth_key, auth_login, encpwd)

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
func (p *sysUserBaseService) ChangePassword(userid string, newpwd string) (utils.Map, error) {

	log.Println("SysUserService::ChangePassword - Begin")
	indata := utils.Map{
		platform_common.FLD_SYS_USER_PASSWORD: utils.SHA(newpwd),
	}
	data, err := p.daoSysUser.Update(userid, indata)

	log.Println("SysUserService::ChangePassword - End ")
	return data, err
}
