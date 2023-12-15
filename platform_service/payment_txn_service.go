package platform_service

import (
	"log"
	"strings"
	"time"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

// PaymentTxnService - PaymentTxns Service structure
type PaymentTxnService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(PaymentTxnId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(PaymentTxnId string, indata utils.Map) (utils.Map, error)
	Delete(PaymentTxnId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// PaymentTxnBaseService - PaymentTxns Service structure
type PaymentTxnBaseService struct {
	db_utils.DatabaseService
	daoPaymentTxn platform_repository.PaymentTxnDao
	child         PaymentTxnService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewPaymentTxnService(props utils.Map) (PaymentTxnService, error) {

	p := PaymentTxnBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewIndustryDBService ", err)
		return nil, err
	}
	log.Printf("IndustryDBService ")

	// Instantiate other services
	p.daoPaymentTxn = platform_repository.NewPaymentTxnDao(p.GetClient())

	p.child = &p

	return &p, nil
}

func (p *PaymentTxnBaseService) EndService() {
	p.CloseDatabaseService()

}

// List - List All records
func (p *PaymentTxnBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("PaymentTxnService::FindAll - Begin")

	daoPaymentTxn := p.daoPaymentTxn
	response, err := daoPaymentTxn.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("PaymentTxnService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *PaymentTxnBaseService) Get(PaymentTxnId string) (utils.Map, error) {
	log.Printf("PaymentTxnService::FindByCode::  Begin %v", PaymentTxnId)

	data, err := p.daoPaymentTxn.Get(PaymentTxnId)
	log.Println("PaymentTxnService::FindByCode:: End ", err)
	return data, err
}

func (p *PaymentTxnBaseService) Find(filter string) (utils.Map, error) {
	log.Println("PaymentTxnService::FindByCode::  Begin ", filter)

	data, err := p.daoPaymentTxn.Find(filter)
	log.Println("PaymentTxnService::FindByCode:: End ", data, err)
	return data, err
}

func (p *PaymentTxnBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var PaymentTxn_Id string

	dataval, dataok := indata[platform_common.FLD_PAYMENT_TXN_ID]
	if dataok {
		PaymentTxn_Id = strings.ToLower(dataval.(string))
	} else {
		PaymentTxn_Id = utils.GenerateUniqueId("pay_txn")
		log.Println("Unique PaymentTxn ID", PaymentTxn_Id)
	}
	dateTime := time.Now().Format(time.DateTime)
	indata[platform_common.FLD_DATE_TIME] = dateTime
	indata[platform_common.FLD_PAYMENT_TXN_ID] = PaymentTxn_Id
	log.Println("Provided PaymentTxn ID:", PaymentTxn_Id)

	_, err := p.daoPaymentTxn.Get(PaymentTxn_Id)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing PaymentTxn ID !", ErrorDetail: "Given PaymentTxn ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoPaymentTxn.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *PaymentTxnBaseService) Update(PaymentTxnId string, indata utils.Map) (utils.Map, error) {

	log.Println("PaymentTxnService::Update - Begin")

	data, err := p.daoPaymentTxn.Get(PaymentTxnId)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, platform_common.FLD_PAYMENT_TXN_ID)
	delete(indata, platform_common.FLD_BUSINESS_ID)

	data, err = p.daoPaymentTxn.Update(PaymentTxnId, indata)
	log.Println("PaymentTxnService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *PaymentTxnBaseService) Delete(PaymentTxnId string, delete_permanent bool) error {

	log.Println("PaymentTxnService::Delete - Begin", PaymentTxnId)

	daoPaymentTxn := p.daoPaymentTxn
	if delete_permanent {
		result, err := daoPaymentTxn.Delete(PaymentTxnId)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(PaymentTxnId, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("PaymentTxnService::Delete - End")
	return nil
}
