package platform_services

import (
	"fmt"
	"log"
	"strings"

	"github.com/rs/xid"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// AppRoleService - Users Service structure
type AppRoleService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(role_id string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(role_id string, indata utils.Map) (utils.Map, error)
	Delete(role_id string) error

	AddCredentials(role_id string, indata utils.Map) (utils.Map, error)
	GetCredentials(role_id string) (utils.Map, error)
	FindCredential(filter string) (utils.Map, error)

	AddUsers(role_id string, indata utils.Map) (utils.Map, error)
	FindUser(filter string) (utils.Map, error)
	GetUsers(rold_id string) (utils.Map, error)

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

type appRoleBaseService struct {
	db_utils.DatabaseService
	daoAppRole platform_repository.AppRoleDao
	child      AppRoleService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewAppRoleService(props utils.Map) (AppRoleService, error) {
	p := appRoleBaseService{}

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewAppRoleService ", err)
		return nil, err
	}

	log.Printf("appRoleService ")
	p.daoAppRole = platform_repository.NewAppRoleDao(p.GetClient())
	p.child = &p

	return &p, nil
}

func (p *appRoleBaseService) EndService() {
	p.CloseDatabaseService()
}

// List - List All records
func (p *appRoleBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("AppRoleService::FindAll - Begin")

	dataresponse, err := p.daoAppRole.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("AppRoleService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *appRoleBaseService) Get(role_id string) (utils.Map, error) {
	log.Printf("AppRoleService::GetDetails::  Begin %v", role_id)

	data, err := p.daoAppRole.Get(role_id)

	log.Println("AppRoleService::GetDetails:: End ", data, err)
	return data, err
}

func (p *appRoleBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("AppRoleService::GetDetails::  Begin ", filter)

	data, err := p.daoAppRole.Find(filter)

	log.Println("AppRoleService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *appRoleBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")
	var roleId string

	dataval, dataok := indata[platform_common.FLD_APP_ROLE_ID]
	if dataok {
		roleId = strings.ToLower((dataval.(string)))
	} else {
		guid := xid.New()
		prefix := "role_"
		log.Println("Unique Role ID", guid.String())
		roleId = prefix + guid.String()
	}
	log.Println("Provided Role ID:", roleId)

	// Update the new Id
	indata[platform_common.FLD_APP_ROLE_ID] = roleId

	dataCreated, err := p.daoAppRole.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ")
	return dataCreated, nil
}

// Update - Update Service
func (p *appRoleBaseService) Update(role_id string, indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Update - Begin")

	// Delete the Key fields
	delete(indata, platform_common.FLD_APP_ROLE_ID)

	data, err := p.daoAppRole.Update(role_id, indata)

	log.Println("UserService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *appRoleBaseService) Delete(role_id string) error {

	log.Println("UserService::Delete - Begin", role_id)

	result, err := p.daoAppRole.Delete(role_id)
	if err != nil {
		return err
	}

	log.Printf("UserService::Delete - End %v", result)
	return nil
}

// Create - Create Service
func (p *appRoleBaseService) AddCredentials(role_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AddCredentials::Add - Begin")

	log.Println("Provided Role ID:", role_id, indata)

	dataRes, err := p.daoAppRole.AddCredentials(role_id, indata)
	if err != nil {
		return indata, err
	}
	log.Println("AddCredentials::Add - End ")
	return dataRes, nil
}

// Check Credentails - Get Credentail by given role and credentail
func (p *appRoleBaseService) FindCredential(filter string) (utils.Map, error) {

	log.Println("AddCredentials::Add - Begin")

	log.Println("Provided Role ID:", filter)

	dataRes, err := p.daoAppRole.FindCredential(filter)
	if err != nil {
		return nil, err
	}
	log.Println("AddCredentials::Add - End ")
	return dataRes, nil
}

// Check Credentails - Get Credentail by given role and credentail
func (p *appRoleBaseService) GetCredentials(rold_id string) (utils.Map, error) {

	log.Println("GetCredentials::Get - Begin")

	log.Println("Provided Role ID:", rold_id)

	dataCreds, err := p.daoAppRole.GetCredentials(rold_id)
	if err != nil {
		return nil, err
	}
	log.Println("GetCredentials::Get - End ")
	return dataCreds, nil
}

// Create - Create Service
func (p *appRoleBaseService) AddUsers(role_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AddUsers::Add - Begin")

	log.Println("Provided Role ID:", role_id, indata)

	dataRes, err := p.daoAppRole.AddUsers(role_id, indata)
	if err != nil {
		return indata, err
	}
	log.Println("AddUsers::Add - End ")
	return dataRes, nil
}

// Check Credentails - Get Credentail by given role and credentail
func (p *appRoleBaseService) FindUser(filter string) (utils.Map, error) {

	log.Println("FindUser::Add - Begin")

	log.Println("Provided Role ID:", filter)

	dataRes, err := p.daoAppRole.FindUser(filter)
	if err != nil {
		return nil, err
	}
	log.Println("FindUser::Add - End ")
	return dataRes, nil
}

// GetUsers - Get Users by given role and credentail
func (p *appRoleBaseService) GetUsers(rold_id string) (utils.Map, error) {

	log.Println("GetUsers::Get - Begin")

	log.Println("Provided Role ID:", rold_id)

	dataRes, err := p.daoAppRole.GetUsers(rold_id)
	if err != nil {
		return nil, err
	}
	log.Println("GetUsers::Get - End ")
	return dataRes, nil
}
