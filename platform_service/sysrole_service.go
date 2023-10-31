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

// SysRoleService - Users Service structure
type SysRoleService interface {
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

type sysRoleBaseService struct {
	db_utils.DatabaseService
	daoSysRole platform_repository.SysRoleDao
	child      SysRoleService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewSysRoleService(props utils.Map) (SysRoleService, error) {
	p := sysRoleBaseService{}

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewSysRoleMongoService ", err)
		return nil, err
	}

	log.Printf("sysRoleMongoService ")
	p.daoSysRole = platform_repository.NewSysRoleDao(p.GetClient())
	p.child = &p

	return &p, nil
}

func (p *sysRoleBaseService) EndService() {
	log.Printf("EndsysRoleService ")
	p.CloseDatabaseService()
}

// List - List All records
func (p *sysRoleBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("SysRoleService::FindAll - Begin")

	dataresponse, err := p.daoSysRole.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("SysRoleService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *sysRoleBaseService) Get(role_id string) (utils.Map, error) {
	log.Printf("SysRoleService::GetDetails::  Begin %v", role_id)

	data, err := p.daoSysRole.Get(role_id)

	log.Println("SysRoleService::GetDetails:: End ", data, err)
	return data, err
}

func (p *sysRoleBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("SysRoleService::GetDetails::  Begin ", filter)

	data, err := p.daoSysRole.Find(filter)

	log.Println("SysRoleService::GetDetails:: End ", data, err)
	return data, err
}

// Create - Create Service
func (p *sysRoleBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var sysRoleId string

	dataval, dataok := indata[platform_common.FLD_SYS_ROLE_ID]
	if dataok {
		sysRoleId = strings.ToLower(dataval.(string))
	} else {
		guid := xid.New()
		prefix := "syrol_"
		log.Println("Unique Role ID", prefix, guid.String())
		sysRoleId = prefix + guid.String()
	}
	log.Println("Provided Role ID:", dataval)

	// Update converted/generated id back to indata
	indata[platform_common.FLD_SYS_ROLE_ID] = sysRoleId

	dataCreated, err := p.daoSysRole.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ")
	return dataCreated, nil
}

// Update - Update Service
func (p *sysRoleBaseService) Update(role_id string, indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Update - Begin")

	// Delete the Key fields
	delete(indata, platform_common.FLD_SYS_ROLE_ID)

	data, err := p.daoSysRole.Update(role_id, indata)

	log.Println("UserService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *sysRoleBaseService) Delete(role_id string) error {

	log.Println("UserService::Delete - Begin", role_id)

	result, err := p.daoSysRole.Delete(role_id)
	if err != nil {
		return err
	}

	log.Printf("UserService::Delete - End %v", result)
	return nil
}

// Create - Create Service
func (p *sysRoleBaseService) AddCredentials(role_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AddCredentials::Add - Begin")

	log.Println("Provided Role ID:", role_id, indata)

	dataRes, err := p.daoSysRole.AddCredentials(role_id, indata)
	if err != nil {
		return indata, err
	}
	log.Println("AddCredentials::Add - End ")
	return dataRes, nil
}

// Check Credentails - Get Credentail by given role and credentail
func (p *sysRoleBaseService) FindCredential(filter string) (utils.Map, error) {

	log.Println("AddCredentials::Add - Begin")

	log.Println("Provided Role ID:", filter)

	dataRes, err := p.daoSysRole.FindCredential(filter)
	if err != nil {
		return nil, err
	}
	log.Println("AddCredentials::Add - End ")
	return dataRes, nil
}

// Check Credentails - Get Credentail by given role and credentail
func (p *sysRoleBaseService) GetCredentials(rold_id string) (utils.Map, error) {

	log.Println("GetCredentials::Get - Begin")

	log.Println("Provided Role ID:", rold_id)

	dataRes, err := p.daoSysRole.GetCredentials(rold_id)
	if err != nil {
		return nil, err
	}
	log.Println("GetCredentials::Get - End ")
	return dataRes, nil
}

// Create - Create Service
func (p *sysRoleBaseService) AddUsers(role_id string, indata utils.Map) (utils.Map, error) {

	log.Println("AddUsers::Add - Begin")

	log.Println("Provided Role ID:", role_id, indata)

	dataRes, err := p.daoSysRole.AddUsers(role_id, indata)
	if err != nil {
		return indata, err
	}
	log.Println("AddUsers::Add - End ")
	return dataRes, nil
}

// Check Credentails - Get Credentail by given role and credentail
func (p *sysRoleBaseService) FindUser(filter string) (utils.Map, error) {

	log.Println("FindUser::Add - Begin")

	log.Println("Provided Role ID:", filter)

	dataRes, err := p.daoSysRole.FindUser(filter)
	if err != nil {
		return nil, err
	}
	log.Println("FindUser::Add - End ")
	return dataRes, nil
}

// GetUsers - Get Users by given role and credentail
func (p *sysRoleBaseService) GetUsers(rold_id string) (utils.Map, error) {

	log.Println("GetUsers::Get - Begin")

	log.Println("Provided Role ID:", rold_id)

	dataRes, err := p.daoSysRole.GetUsers(rold_id)
	if err != nil {
		return nil, err
	}
	log.Println("GetUsers::Get - End ")
	return dataRes, nil
}
