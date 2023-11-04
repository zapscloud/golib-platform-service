package platform_services

import (
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// InvoiceService - Invoices Service structure
type InvoiceService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(InvoiceId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(InvoiceId string, indata utils.Map) (utils.Map, error)
	Delete(InvoiceId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// InvoiceBaseService - Invoices Service structure
type InvoiceBaseService struct {
	db_utils.DatabaseService
	daoInvoice platform_repository.InvoiceDao
	child      InvoiceService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewInvoiceService(props utils.Map) (InvoiceService, error) {

	p := InvoiceBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewIndustryDBService ", err)
		return nil, err
	}
	log.Printf("IndustryDBService ")

	// Instantiate other services
	p.daoInvoice = platform_repository.NewInvoiceDao(p.GetClient())

	p.child = &p

	return &p, nil
}

func (p *InvoiceBaseService) EndService() {
	p.CloseDatabaseService()

}

// List - List All records
func (p *InvoiceBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("InvoiceService::FindAll - Begin")

	daoInvoice := p.daoInvoice
	response, err := daoInvoice.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("InvoiceService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *InvoiceBaseService) Get(InvoiceId string) (utils.Map, error) {
	log.Printf("InvoiceService::FindByCode::  Begin %v", InvoiceId)

	data, err := p.daoInvoice.Get(InvoiceId)
	log.Println("InvoiceService::FindByCode:: End ", err)
	return data, err
}

func (p *InvoiceBaseService) Find(filter string) (utils.Map, error) {
	log.Println("InvoiceService::FindByCode::  Begin ", filter)

	data, err := p.daoInvoice.Find(filter)
	log.Println("InvoiceService::FindByCode:: End ", data, err)
	return data, err
}

func (p *InvoiceBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var InvoiceId string

	dataval, dataok := indata[platform_common.FLD_INVOICE_ID]
	if dataok {
		InvoiceId = strings.ToLower(dataval.(string))
	} else {
		InvoiceId = utils.GenerateUniqueId("inice")
		log.Println("Unique Invoice ID", InvoiceId)
	}
	indata[platform_common.FLD_INVOICE_ID] = InvoiceId
	log.Println("Provided Invoice ID:", InvoiceId)

	_, err := p.daoInvoice.Get(InvoiceId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Invoice ID !", ErrorDetail: "Given Invoice ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoInvoice.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *InvoiceBaseService) Update(InvoiceId string, indata utils.Map) (utils.Map, error) {

	log.Println("InvoiceService::Update - Begin")

	data, err := p.daoInvoice.Get(InvoiceId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, platform_common.FLD_INVOICE_ID)
	delete(indata, platform_common.FLD_BUSINESS_ID)

	data, err = p.daoInvoice.Update(InvoiceId, indata)
	log.Println("InvoiceService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *InvoiceBaseService) Delete(InvoiceId string, delete_permanent bool) error {

	log.Println("InvoiceService::Delete - Begin", InvoiceId)

	daoInvoice := p.daoInvoice
	if delete_permanent {
		result, err := daoInvoice.Delete(InvoiceId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(InvoiceId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("InvoiceService::Delete - End")
	return nil
}
