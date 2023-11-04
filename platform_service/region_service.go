package platform_service

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// RegionService - Users Service structure
type RegionService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(regionid string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(regionid string, indata utils.Map) (utils.Map, error)
	Delete(regionid string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type regionBaseService struct {
	db_utils.DatabaseService
	daoRegion platform_repository.RegionDao
	child     RegionService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewRegionService(props utils.Map) (RegionService, error) {
	p := regionBaseService{}

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewRegionMongoService Connection Error ", err)
		return nil, err
	}

	log.Printf("RegionMongoService ")
	p.daoRegion = platform_repository.NewRegionDao(p.GetClient())
	p.child = &p

	return &p, nil
}

func (p *regionBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *regionBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("RegionService::FindAll - Begin")

	dataresponse, err := p.daoRegion.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("RegionService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *regionBaseService) Get(regionid string) (utils.Map, error) {
	log.Printf("RegionService::GetDetails::  Begin %v", regionid)

	data, err := p.daoRegion.Get(regionid)

	log.Println("RegionService::GetDetails:: End ", data, err)
	return data, err
}

func (p *regionBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("RegionService::GetDetails::  Begin ", filter)

	data, err := p.daoRegion.Find(filter)

	log.Println("RegionService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *regionBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	// Conver the RegionId to Lowercase
	regionId := strings.ToLower(indata[platform_common.FLD_REGION_ID].(string))

	indata, err := p.validateCreate(indata)
	if err != nil {
		return nil, err
	}

	log.Println("Provided Profile ID:", regionId)
	// Update converted regionId back to indata
	indata[platform_common.FLD_REGION_ID] = regionId

	_, err = p.daoRegion.Get(regionId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Region ID !", ErrorDetail: "Given Region ID already exist"}
		return indata, err
	}

	_, err = p.daoRegion.Create(indata)
	if err != nil {
		return indata, err
	}

	log.Println("UserService::Create - End ")
	return p.daoRegion.Get(regionId)
}

// Update - Update Service
func (p *regionBaseService) Update(regionid string, indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Update - Begin")

	result, err := p.validateKeyExist(regionid)
	if err != nil {
		return result, err
	}

	// Remove key and default fields from indata
	delete(indata, platform_common.FLD_REGION_ID)

	data, err := p.daoRegion.Update(regionid, indata)

	log.Println("UserService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *regionBaseService) Delete(regionid string, delete_permanent bool) error {

	log.Println("UserService::Delete - Begin", regionid)

	_, err := p.validateKeyExist(regionid)
	if err != nil {
		return err
	}

	if delete_permanent {
		result, err := p.daoRegion.Delete(regionid)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(regionid, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("UserService::Delete")
	return nil
}

func (p *regionBaseService) validateKeyExist(key string) (utils.Map, error) {
	data, err := p.daoRegion.Get(key)
	if err != nil {
		err := &utils.AppError{ErrorStatus: 400, ErrorMsg: "Bad Request", ErrorDetail: "Region not found"}
		return utils.Map{}, err
	}
	return data, nil
}

// Private functions
func (p *regionBaseService) validateCreate(dataRegion utils.Map) (utils.Map, error) {

	var err error = nil

	if _, err := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_ID); err != nil {
		err = &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Missing Region ID!", ErrorDetail: "Region ID parameter is missing"}
	} else if _, err := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_NAME); err != nil {
		err = &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Missing value", ErrorDetail: "Parameter " + platform_common.FLD_REGION_NAME + " is missing"}
	} else if _, err := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_DB_TYPE); err != nil {
		err = &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Missing value", ErrorDetail: "Parameter " + platform_common.FLD_REGION_DB_TYPE + " is missing"}
	}

	return dataRegion, err
}
