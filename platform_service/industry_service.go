package platform_service

import (
	"fmt"
	"log"

	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

type IndustryService interface {
	List(filter string, sort string, skip int64, limit int64) (utils.Map, error)
	GetIndustryById(industryid string) (utils.Map, error)
	EndService()
}

type industryBaseService struct {
	db_utils.DatabaseService
	daoIndustry platform_repository.IndustryDao
	child       IndustryService
}

func NewIndustryService(props utils.Map) (IndustryService, error) {
	p := industryBaseService{}

	err := p.OpenDatabaseService(props)
	if err != nil {
		log.Println("NewIndustryDBService ", err)
		return nil, err
	}
	log.Printf("IndustryDBService ")
	p.child = &p

	p.daoIndustry = platform_repository.NewIndustryDao(p.GetClient())
	return &p, nil
}

// EndIndustryService - Close all the services
func (p *industryBaseService) EndService() {
	log.Printf("EndIndustryDBService ")
	p.CloseDatabaseService()
}

// List - List All records
func (p *industryBaseService) List(filter string, sort string, skip int64, limit int64) (utils.Map, error) {

	log.Println("CustomerService::FindAll - Begin")

	listdata, err := p.daoIndustry.List(filter, sort, skip, limit)
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}

	log.Println("CustomerService::FindAll - End ")
	return listdata, nil
}

func (p *industryBaseService) GetIndustryById(industryid string) (utils.Map, error) {
	log.Println("CustomerService::FindAll - Begin")

	industryData, err := p.daoIndustry.GetIndustryById(industryid)
	if err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}

	log.Println("CustomerService::FindAll - End ")
	return industryData, nil

}
