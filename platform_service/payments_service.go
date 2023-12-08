package platform_service

import (
	"log"
	"strings"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// PaymentsService - Paymentss Service structure
type PaymentsService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(PaymentsId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(PaymentsId string, indata utils.Map) (utils.Map, error)
	Delete(PaymentsId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// PaymentsBaseService - Paymentss Service structure
type PaymentsBaseService struct {
	db_utils.DatabaseService
	daoPayments platform_repository.PaymentsDao
	child      PaymentsService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewPaymentsService(props utils.Map) (PaymentsService, error) {

	p := PaymentsBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewIndustryDBService ", err)
		return nil, err
	}
	log.Printf("IndustryDBService ")

	// Instantiate other services
	p.daoPayments = platform_repository.NewPaymentsDao(p.GetClient())

	p.child = &p

	return &p, nil
}

func (p *PaymentsBaseService) EndService() {
	p.CloseDatabaseService()

}

// List - List All records
func (p *PaymentsBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("PaymentsService::FindAll - Begin")

	daoPayments := p.daoPayments
	response, err := daoPayments.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("PaymentsService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *PaymentsBaseService) Get(PaymentsId string) (utils.Map, error) {
	log.Printf("PaymentsService::FindByCode::  Begin %v", PaymentsId)

	data, err := p.daoPayments.Get(PaymentsId)
	log.Println("PaymentsService::FindByCode:: End ", err)
	return data, err
}

func (p *PaymentsBaseService) Find(filter string) (utils.Map, error) {
	log.Println("PaymentsService::FindByCode::  Begin ", filter)

	data, err := p.daoPayments.Find(filter)
	log.Println("PaymentsService::FindByCode:: End ", data, err)
	return data, err
}

func (p *PaymentsBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var PaymentsId string

	dataval, dataok := indata[platform_common.FLD_PAYMENT_ID]
	if dataok {
		PaymentsId = strings.ToLower(dataval.(string))
	} else {
		PaymentsId = utils.GenerateUniqueId("pay")
		log.Println("Unique Payments ID", PaymentsId)
	}
	indata[platform_common.FLD_PAYMENT_ID] = PaymentsId
	log.Println("Provided Payments ID:", PaymentsId)

	_, err := p.daoPayments.Get(PaymentsId)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Payments ID !", ErrorDetail: "Given Payments ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoPayments.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *PaymentsBaseService) Update(PaymentsId string, indata utils.Map) (utils.Map, error) {

	log.Println("PaymentsService::Update - Begin")

	data, err := p.daoPayments.Get(PaymentsId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, platform_common.FLD_PAYMENT_ID)
	delete(indata, platform_common.FLD_BUSINESS_ID)

	data, err = p.daoPayments.Update(PaymentsId, indata)
	log.Println("PaymentsService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *PaymentsBaseService) Delete(PaymentsId string, delete_permanent bool) error {

	log.Println("PaymentsService::Delete - Begin", PaymentsId)

	daoPayments := p.daoPayments
	if delete_permanent {
		result, err := daoPayments.Delete(PaymentsId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(PaymentsId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("PaymentsService::Delete - End")
	return nil
}
