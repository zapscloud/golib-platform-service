package platform_service

import (
	"fmt"
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// SysSettingService - Users Service structure
type SysSettingService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(clientid string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (string, error)
	Update(clientid string, indata utils.Map) (utils.Map, error)
	Delete(clientid string) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type appSettingBaseService struct {
	db_utils.DatabaseService
	daoSysSetting platform_repository.SysSettingDao
	child         SysSettingService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewSysSettingService(props utils.Map) (SysSettingService, error) {
	p := appSettingBaseService{}

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewSysSettingService: Connection Error ", err)
		return nil, err
	}
	log.Printf("NewSysSettingService ")

	p.daoSysSetting = platform_repository.NewSysSettingDao(p.GetClient())
	p.child = &p

	return &p, nil
}

func (p *appSettingBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *appSettingBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("SysSettingService::FindAll - Begin")

	dataresponse, err := p.daoSysSetting.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("SysSettingService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *appSettingBaseService) Get(clientid string) (utils.Map, error) {
	log.Printf("SysSettingService::GetDetails::  Begin %v", clientid)

	data, err := p.daoSysSetting.Get(clientid)

	log.Println("SysSettingService::GetDetails:: End ", data, err)
	return data, err
}

func (p *appSettingBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("SysSettingService::GetDetails::  Begin ", filter)

	data, err := p.daoSysSetting.Find(filter)

	log.Println("SysSettingService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *appSettingBaseService) Create(indata utils.Map) (string, error) {

	log.Println("ClientService::Create - Begin")
	var settingsId string

	dataval, dataok := indata[platform_common.FLD_SETTING_ID]
	if dataok {
		settingsId = strings.ToLower(dataval.(string))
	} else {
		err := &utils.AppError{ErrorCode: "S3040101", ErrorMsg: "Missing app_setting_id", ErrorDetail: "Missing required field app_setting_id !!"}
		return "", err
	}
	log.Println("Provided Settings ID:", settingsId)

	_, err := p.daoSysSetting.Get(settingsId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S3040102", ErrorMsg: "Existing app_setting_id", ErrorDetail: "Given app_setting_id is already exist"}
		return settingsId, err
	}

	createdId, err := p.daoSysSetting.Create(indata)
	if err != nil {
		return "", err
	}
	log.Println("ClientService::Create - End ")
	return createdId, nil
}

// Update - Update Service
func (p *appSettingBaseService) Update(clientid string, indata utils.Map) (utils.Map, error) {

	log.Println("ClientService::Update - Begin")

	// Delete the Key fields
	delete(indata, platform_common.FLD_SETTING_ID)

	data, err := p.daoSysSetting.Update(clientid, indata)

	log.Println("ClientService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *appSettingBaseService) Delete(clientid string) error {

	log.Println("ClientService::Delete - Begin", clientid)

	result, err := p.daoSysSetting.Delete(clientid)
	if err != nil {
		return err
	}

	log.Printf("ClientService::Delete - End %v", result)
	return nil
}
