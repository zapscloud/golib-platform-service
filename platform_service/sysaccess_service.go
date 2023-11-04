package platform_service

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// AccessService - Accesss Service structure
type SysAccessService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(appaccessid string) (utils.Map, error)

	GrantPermission(indata utils.Map) (utils.Map, error)
	RevokePermission(access_id string) (int64, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type sysAccessBaseService struct {
	db_utils.DatabaseService
	daoSysAccess platform_repository.SysAccessDao
	daoSysUser   platform_repository.SysUserDao
	daoAppUser   platform_repository.AppUserDao
	daoBusiness  platform_repository.BusinessDao
	child        SysAccessService
	businessID   string
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewSysAccessService(props utils.Map) (SysAccessService, error) {
	funcode := platform_common.GetServiceModuleCode() + "M" + "01"

	p := sysAccessBaseService{}

	// Open Database Service
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Fatal(err)
	}

	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, platform_common.FLD_BUSINESS_ID)
	if err != nil {
		p.CloseDatabaseService()
		return nil, err
	}

	// Assign the BusinessId
	p.businessID = businessId
	log.Printf("SysAccessMongoService ")

	p.daoSysAccess = platform_repository.NewSysAccessDao(p.GetClient(), p.businessID)
	p.daoSysUser = platform_repository.NewSysUserDao(p.GetClient())
	p.daoAppUser = platform_repository.NewAppUserDao(p.GetClient())
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())

	_, err = p.daoBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid sys_business_id", ErrorDetail: "Given sys_business_id is not exist"}
		return nil, err
	}

	p.child = &p

	return &p, err
}

// EndSysAccessService - Close all the services
func (p *sysAccessBaseService) EndService() {
	log.Printf("EndAccessService ")
	p.CloseDatabaseService()
}

func (p *sysAccessBaseService) getServiceModuleCode() string {
	return platform_common.GetServiceModuleCode() + "05"
}

// List - List All records
func (p *sysAccessBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AccessService::FindAll - Begin")

	daoAccess := p.daoSysAccess
	response, err := daoAccess.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("AccessService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *sysAccessBaseService) Get(appaccessid string) (utils.Map, error) {
	log.Printf("AccessService::FindByCode::  Begin %v", appaccessid)

	data, err := p.daoSysAccess.Get(appaccessid)
	log.Println("AccessService::FindByCode:: End ", err)
	return data, err
}

func (p *sysAccessBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("AccessService::FindByCode::  Begin ", filter)

	data, err := p.daoSysAccess.Find(filter)
	log.Println("AccessService::FindByCode:: End ", data, err)
	return data, err
}

// Update - Update Service
func (p *sysAccessBaseService) GrantPermission(indata utils.Map) (utils.Map, error) {

	funcode := p.getServiceModuleCode() + "01"

	log.Println("AccessService::Update - Begin")

	access_key := ""
	if valUserId, okUserId := indata[platform_common.FLD_SYS_USER_ID]; !okUserId {
		log.Println("GrantPermission: UserId not found  ", valUserId)
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "UserId not found ", ErrorDetail: "UserId not found "}
		return indata, err
	} else if _, err := p.daoSysUser.Get(valUserId.(string)); err != nil {
		log.Println("GrantPermission: UserId not found  ", valUserId)
		err := &utils.AppError{ErrorCode: funcode + "02", ErrorMsg: "UserId not found ", ErrorDetail: "UserId not found "}
		return indata, err
	} else {
		access_key += valUserId.(string)
	}

	if valRoleId, okRoleId := indata[platform_common.FLD_SYS_ROLE_ID]; !okRoleId {
		log.Println("GrantPermission: Missing RoleId ", valRoleId)
		err := &utils.AppError{ErrorCode: funcode + "03", ErrorMsg: "Missing RoleId ", ErrorDetail: "Missing RoleId "}
		return indata, err

	} else if _, err := p.daoSysAccess.GetRoleDetails(valRoleId.(string)); err != nil {
		log.Println("GrantPermission: RoleId not found  ", valRoleId)
		err := &utils.AppError{ErrorCode: funcode + "04", ErrorMsg: "UserId not found ", ErrorDetail: "UserId not found "}
		return indata, err
	}

	if valSiteId, okSiteId := indata["app_site_id"]; !okSiteId {
		// Ignore Site Id Field
		access_key = "-" + access_key
	} else if _, err := p.daoSysAccess.GetSiteDetails(valSiteId.(string)); err != nil {
		log.Println("GrantPermission: RoleId not found  ", valSiteId)
		err := &utils.AppError{ErrorCode: funcode + "04", ErrorMsg: "UserId not found ", ErrorDetail: "UserId not found "}
		return indata, err
	} else {
		access_key += valSiteId.(string)
	}

	if valDeptId, okDeptId := indata["app_dept_id"]; !okDeptId {
		// Ignore Department ID
		access_key = "-" + access_key
	} else if _, err := p.daoSysAccess.GetDepartmentDetails(valDeptId.(string)); err != nil {
		log.Println("GrantPermission: RoleId not found  ", valDeptId)
		err := &utils.AppError{ErrorCode: funcode + "04", ErrorMsg: "UserId not found ", ErrorDetail: "UserId not found "}
		return indata, err
	} else {
		access_key += valDeptId.(string)
	}

	access_id := utils.GenerateChecksumId("aces", access_key)
	dataAccess, err := p.daoSysAccess.Get(access_id)
	if err != nil {
		return dataAccess, err
	}

	indata[platform_common.FLD_SYS_ACCESS_ID] = access_id

	dataAccess, err = p.daoSysAccess.GrantPermission(indata)
	log.Println("AccessService::Update - End ")
	return dataAccess, err
}

// RevokePermission - RevokePermission Service
func (p *sysAccessBaseService) RevokePermission(access_id string) (int64, error) {

	log.Println("AccessService::RevokePermission - Begin", access_id)

	daoUser := p.daoSysAccess
	result, err := daoUser.RevokePermission(access_id)
	if err != nil {
		return result, err
	}

	log.Printf("UserService::Delete - End %v", result)
	return result, nil
}
