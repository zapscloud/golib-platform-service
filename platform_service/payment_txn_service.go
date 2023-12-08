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

// Payment_txnService - Payment_txns Service structure
type Payment_txnService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	Get(Payment_txnId string) (utils.Map, error)
	Find(filter string) (utils.Map, error)
	Create(indata utils.Map) (utils.Map, error)
	Update(Payment_txnId string, indata utils.Map) (utils.Map, error)
	Delete(Payment_txnId string, delete_permanent bool) error

	BeginTransaction()
	CommitTransaction()
	RollbackTransaction()

	EndService()
}

// Payment_txnBaseService - Payment_txns Service structure
type Payment_txnBaseService struct {
	db_utils.DatabaseService
	daoPayment_txn platform_repository.Payment_txnDao
	child      Payment_txnService
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func NewPayment_txnService(props utils.Map) (Payment_txnService, error) {

	p := Payment_txnBaseService{}
	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewIndustryDBService ", err)
		return nil, err
	}
	log.Printf("IndustryDBService ")

	// Instantiate other services
	p.daoPayment_txn = platform_repository.NewPayment_txnDao(p.GetClient())

	p.child = &p

	return &p, nil
}

func (p *Payment_txnBaseService) EndService() {
	p.CloseDatabaseService()

}

// List - List All records
func (p *Payment_txnBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("Payment_txnService::FindAll - Begin")

	daoPayment_txn := p.daoPayment_txn
	response, err := daoPayment_txn.List(filter, sort, skip, limit)
	if err != nil {
		return nil, err
	}

	log.Println("Payment_txnService::FindAll - End ")
	return response, nil
}

// FindByCode - Find By Code
func (p *Payment_txnBaseService) Get(Payment_txnId string) (utils.Map, error) {
	log.Printf("Payment_txnService::FindByCode::  Begin %v", Payment_txnId)

	data, err := p.daoPayment_txn.Get(Payment_txnId)
	log.Println("Payment_txnService::FindByCode:: End ", err)
	return data, err
}

func (p *Payment_txnBaseService) Find(filter string) (utils.Map, error) {
	log.Println("Payment_txnService::FindByCode::  Begin ", filter)

	data, err := p.daoPayment_txn.Find(filter)
	log.Println("Payment_txnService::FindByCode:: End ", data, err)
	return data, err
}

func (p *Payment_txnBaseService) Create(indata utils.Map) (utils.Map, error) {

	log.Println("UserService::Create - Begin")

	var Payment_txn_Id string

	dataval, dataok := indata[platform_common.FLD_PAYMENT_TXN_ID]
	if dataok {
		Payment_txn_Id = strings.ToLower(dataval.(string))
	} else {
		Payment_txn_Id = utils.GenerateUniqueId("pay_txn")
		log.Println("Unique Payment_txn ID", Payment_txn_Id)
	}
	indata[platform_common.FLD_PAYMENT_TXN_ID] = Payment_txn_Id
	log.Println("Provided Payment_txn ID:", Payment_txn_Id)

	_, err := p.daoPayment_txn.Get(Payment_txn_Id)
	if err == nil {
		err := &utils.AppError{ErrorCode: "S30102", ErrorMsg: "Existing Payment_txn ID !", ErrorDetail: "Given Payment_txn ID already exist"}
		return indata, err
	}

	insertResult, err := p.daoPayment_txn.Create(indata)
	if err != nil {
		return indata, err
	}
	log.Println("UserService::Create - End ", insertResult)
	return indata, err
}

// Update - Update Service
func (p *Payment_txnBaseService) Update(Payment_txn_Id string, indata utils.Map) (utils.Map, error) {

	log.Println("Payment_txnService::Update - Begin")

	data, err := p.daoPayment_txn.Get(Payment_txn_Id)
	if err != nil {
		return data, err
	}

	// Delete key fields
	delete(indata, platform_common.FLD_PAYMENT_TXN_ID)
	delete(indata, platform_common.FLD_BUSINESS_ID)

	data, err = p.daoPayment_txn.Update(Payment_txn_Id, indata)
	log.Println("Payment_txnService::Update - End ")
	return data, err
}

// Delete - Delete Service
func (p *Payment_txnBaseService) Delete(Payment_txn_Id string, delete_permanent bool) error {

	log.Println("Payment_txnService::Delete - Begin", Payment_txn_Id)

	daoPayment_txn := p.daoPayment_txn
	if delete_permanent {
		result, err := daoPayment_txn.Delete(Payment_txn_Id)
		if err != nil {
			return err
		}
		log.Printf("Delete %v", result)
	} else {
		indata := utils.Map{db_common.FLD_IS_DELETED: true}
		data, err := p.Update(Payment_txn_Id, indata)
		if err != nil {
			return err
		}
		log.Println("Update for Delete Flag", data)
	}

	log.Printf("Payment_txnService::Delete - End")
	return nil
}
