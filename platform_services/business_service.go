package platform_services

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"

	"github.com/rs/xid"
)

// BusinessService - Users Service structure
type BusinessService interface {

	// Get Business Details
	Get(businessId string) (utils.Map, error)

	// Get Business List
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Find Business
	Find(filter string) (utils.Map, error)

	// Create New Business
	Create(indata utils.Map) (utils.Map, error)

	// Update Business Details
	Update(businessId string, indata utils.Map) (utils.Map, error)

	//Delete Business
	Delete(businessId string, delete_permanent bool) error

	// AddUser Business & User
	AddUser(businessId string, userId string) (utils.Map, error)

	// RemoveUser User from Business
	RemoveUser(businessId string, userId string) (string, error)

	// Get Access Details
	GetUserDetails(businessId string, userId string) (utils.Map, error)

	// Get Users in the business
	GetUsers(businessId string, filter string, sort string, skip int64, limit int64) (utils.Map, error)

	// Get Businesses for given userId form BusinessUser Table
	GetBusinessList(userId string, filter string, sort string, skip int64, limit int64) (utils.Map, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type businessBaseService struct {
	db_utils.DatabaseService
	daoBusiness  platform_repository.BusinessDao
	daoAppUser   platform_repository.AppUserDao
	daoAppRegion platform_repository.RegionDao
	child        BusinessService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewBusinessService(props utils.Map) (BusinessService, error) {
	p := businessBaseService{}

	// Open Database Service
	err := p.OpenDatabaseService(props)
	if err != nil {
		return nil, err
	}

	// Create other Instances
	p.daoBusiness = platform_repository.NewBusinessDao(p.GetClient())
	p.daoAppUser = platform_repository.NewAppUserDao(p.GetClient())
	p.daoAppRegion = platform_repository.NewRegionDao(p.GetClient())

	log.Printf("BusinessService ")
	p.child = &p

	return &p, nil
}

func (p *businessBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *businessBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("BusinessService::FindAll - Begin")

	dataresponse, err := p.daoBusiness.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("BusinessService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *businessBaseService) Get(businessId string) (utils.Map, error) {
	log.Printf("BusinessService::GetDetails::  Begin %v", businessId)

	data, err := p.daoBusiness.Get(businessId)

	log.Println("BusinessService::GetDetails:: End ", data, err)
	return data, err
}

func (p *businessBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("BusinessService::GetDetails::  Begin ", filter)

	data, err := p.daoBusiness.Find(filter)

	log.Println("BusinessService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *businessBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("BusinessService::Create - Begin")
	var businessId string

	dataval, dataok := indata[platform_common.FLD_BUSINESS_ID]
	if dataok {
		// Convert Id to Lowercase
		businessId = strings.ToLower(dataval.(string))
	} else {
		guid := xid.New()
		prefix := "biz_"

		log.Println("Unique Profile ID", prefix, guid.String())
		businessId = prefix + guid.String()
	}
	log.Println("Provided Profile ID:", businessId)

	// Assign new/case converted businessId
	indata[platform_common.FLD_BUSINESS_ID] = businessId

	_, err := p.daoBusiness.Get(businessId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S3030201", ErrorMsg: "Invalid Business id !", ErrorDetail: "Business ID given is already exist"}
		return indata, err
	}

	dataval, dataok = indata[platform_common.FLD_BUSINESS_REGION_ID]
	if !dataok {
		err := &utils.AppError{ErrorCode: "S3030202", ErrorMsg: "Business region missing!", ErrorDetail: "Business Region should be specified!"}
		return indata, err
	}
	_, err = p.daoAppRegion.Get(dataval.(string))
	if err != nil {
		err := &utils.AppError{ErrorCode: "S3030203", ErrorMsg: "Invalid Business region !", ErrorDetail: "Business Region given is invalid"}
		return indata, err
	}

	dataBusiness, err := p.daoBusiness.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("BusinessService::Create - End ")
	return dataBusiness, nil
}

// Update - Update Service
func (p *businessBaseService) Update(businessId string, indata utils.Map) (utils.Map, error) {

	log.Println("BusinessService::Update - Begin")

	// Delete the Key fields
	delete(indata, platform_common.FLD_BUSINESS_ID)
	delete(indata, platform_common.FLD_BUSINESS_REGION_ID)
	delete(indata, platform_common.FLD_BUSINESS_IS_TENANT_DB)

	data, err := p.daoBusiness.Update(businessId, indata)

	log.Println("BusinessService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *businessBaseService) Delete(businessId string, delete_permanent bool) error {

	log.Println("BusinessService::Delete - Begin", businessId)

	_, err := p.validateKeyExist(businessId)
	if err != nil {
		return err
	}

	if delete_permanent {
		result, err := p.daoBusiness.Delete(businessId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(businessId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("BusinessService::Delete")
	return nil
}

// AddUser - Grand Access for the Business
func (p *businessBaseService) AddUser(businessId string, userId string) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	_, err := p.daoBusiness.Get(businessId)
	if err != nil {
		return nil, err
	}

	_, err = p.daoAppUser.Get(userId)
	if err != nil {
		return nil, err
	}

	indata := utils.Map{
		platform_common.FLD_BUSINESS_USER_ID: utils.GetMD5Hash(businessId + "_" + userId),
		platform_common.FLD_BUSINESS_ID:      businessId,
		platform_common.FLD_APP_USER_ID:      userId,
		db_common.FLD_IS_DELETED:             true,
	}

	dataUser, err := p.daoBusiness.AddUser(indata)
	if err != nil {
		return dataUser, err
	}
	log.Println("UserService::Create - End ")
	return dataUser, nil
}

func (p *businessBaseService) RemoveUser(businessId string, userId string) (string, error) {

	log.Println("UserService::RemoveUser - Begin")

	_, err := p.daoBusiness.Get(businessId)
	if err != nil {
		return "", err
	}

	_, err = p.daoAppUser.Get(userId)
	if err != nil {
		return "", err
	}

	accessid := utils.GetMD5Hash(businessId + "_" + userId)
	data, err := p.daoBusiness.GetAccessDetails(accessid)
	if err != nil {
		return "", err
	}
	log.Println("SysBusinessService::RemoveUser ", data)

	revokeresponse, err := p.daoBusiness.RemoveUser(accessid)
	if err != nil {
		return "", err
	}

	log.Println("UserService::RemoveUser - End ", revokeresponse)
	return revokeresponse, nil
}

// GetDetails - Find By Code
func (p *businessBaseService) GetUserDetails(businessId string, userId string) (utils.Map, error) {
	_, err := p.daoBusiness.Get(businessId)
	if err != nil {
		return utils.Map{}, err
	}

	_, err = p.daoAppUser.Get(userId)
	if err != nil {
		return utils.Map{}, err
	}

	accessid := utils.GetMD5Hash(businessId + "_" + userId)
	data, err := p.daoBusiness.GetAccessDetails(accessid)
	if err != nil {
		return utils.Map{}, err
	}
	log.Println("SysBusinessService::GetDetails ", data)

	log.Println("BusinessService::GetDetails:: End ", data, err)
	return data, err
}

// GetUsers - Get all BusinessUser records for given businessId
func (p *businessBaseService) GetUsers(businessId string, filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	_, err := p.daoBusiness.Get(businessId)
	if err != nil {
		return utils.Map{}, err
	}

	data, err := p.daoBusiness.UserList(businessId, filter, sort, skip, limit)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("BusinessService::GetUsers:: End ", data, err)
	return data, err
}

// GetBusinessList - Get all BusinessUser records for given businessId
func (p *businessBaseService) GetBusinessList(userId string, filter string, sort string, skip int64, limit int64) (utils.Map, error) {
	_, err := p.daoAppUser.Get(userId)
	if err != nil {
		return utils.Map{}, err
	}

	data, err := p.daoBusiness.BusinessList(userId, filter, sort, skip, limit)
	if err != nil {
		return utils.Map{}, err
	}

	log.Println("BusinessService::GetBusiness:: End ", data, err)
	return data, err
}

func (p *businessBaseService) validateKeyExist(key string) (utils.Map, error) {
	data, err := p.daoBusiness.Get(key)
	if err != nil {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "BusinessId not found"}
		return utils.Map{}, err
	}
	return data, nil
}
