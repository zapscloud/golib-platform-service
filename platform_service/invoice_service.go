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
)

// InvoiceService - Invoices Service structure
type InvoiceService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(invoiceId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(invoiceId string, indata utils.Map) (utils.Map, error)
	Delete(invoiceId string, deletePermanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

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
	log.Printf("NewInvoiceService :: Start")

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewInvoiceService ", err)
		return nil, err
	}

	p.daoInvoice = platform_repository.NewInvoiceDao(p.GetClient())
	p.child = &p

	return &p, nil
}

func (p *InvoiceBaseService) EndService() {
	log.Printf("EndInvoiceBaseService ")
	p.CloseDatabaseService()
}

// List - List All records
func (p *InvoiceBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("InvoiceService::FindAll - Begin")

	dataresponse, err := p.daoInvoice.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}
	log.Println("InvoiceService::FindAll - End ")
	return dataresponse, nil
}

// GetDetails - Find By Code
func (p *InvoiceBaseService) Get(invoiceId string) (utils.Map, error) {
	log.Printf("InvoiceService::GetDetails::  Begin %v", invoiceId)

	data, err := p.daoInvoice.Get(invoiceId)

	log.Println("InvoiceService::GetDetails:: End ", data, err)
	return data, err
}

func (p *InvoiceBaseService) Find(filter string) (utils.Map, error) {
	fmt.Println("InvoiceService::GetDetails::  Begin ", filter)

	data, err := p.daoInvoice.Find(filter)

	log.Println("InvoiceService::GetDetails:: End ", data, err)
	return data, err
}

func (p *InvoiceBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("InvoiceService::Create - Begin")

	var invoiceId string

	dataval, dataok := indata[platform_common.FLD_INVOICE_ID]
	if dataok {
		invoiceId = strings.ToLower(dataval.(string))
	} else {
		invoiceId = utils.GenerateUniqueId("inice_")
		log.Println("Unique Invoice ID", invoiceId)
	}
	indata[platform_common.FLD_INVOICE_ID] = invoiceId
	log.Println("Provided Invoice ID:",invoiceId)

	_, err := p.daoInvoice.Get(invoiceId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Invoice ID !", ErrorDetail: "Given Invoice ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoInvoice.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("InvoiceService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *InvoiceBaseService) Update(invoiceId string, indata utils.Map) (utils.Map, error) {

	log.Println("InvoiceService::Update - Begin")

	data, err := p.daoInvoice.Update(invoiceId, indata)

	log.Println("InvoiceService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *InvoiceBaseService) Delete(invoiceId string, deletePermanent bool) error {

	log.Println("InvoiceService::Delete - Begin", invoiceId)

	if deletePermanent {
		result, err := p.daoInvoice.Delete(invoiceId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}

		data, err := p.Update(invoiceId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("InvoiceService::Delete - End")
	return nil
}
