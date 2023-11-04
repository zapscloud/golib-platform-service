package platform_service

import (
	"log"

	"github.com/zapscloud/golib-dbutils/db_common"
	"github.com/zapscloud/golib-dbutils/db_utils"
	"github.com/zapscloud/golib-platform-repository/platform_common"
	"github.com/zapscloud/golib-platform-repository/platform_repository"
	"github.com/zapscloud/golib-utils/utils"
)

func OpenRegionDatabaseService(props utils.Map) (db_utils.DatabaseService, error) {
	var dbRegion db_utils.DatabaseService

	// Get Region and Tenant Database Information
	propsRegion, err := getRegionAndTenantDBInfo(props)
	if err == nil {
		err = dbRegion.OpenDatabaseService(propsRegion)
	}
	return dbRegion, err
}

func getRegionAndTenantDBInfo(props utils.Map) (utils.Map, error) {
	var dbServices db_utils.DatabaseService
	var daoPlatformBusiness platform_repository.BusinessDao
	funcode := platform_common.GetServiceModuleCode() + "M" + "01"

	// Verify whether the business id data passed
	businessId, err := utils.GetMemberDataStr(props, platform_common.FLD_BUSINESS_ID)
	if err != nil {
		log.Println("No BusinessId found, may be opening platform database")
		// No BusinessId avaible, so it might be opening platform Database
		return props, nil
	}

	// Open Platform Database first
	err = dbServices.OpenDatabaseService(props)
	if err != nil {
		log.Println("GetTenantDBInfo:: Error wile Open Database ", err)
		return utils.Map{}, err
	}
	defer dbServices.CloseDatabaseService()

	// Open Instance of PlatformBusiness
	daoPlatformBusiness = platform_repository.NewBusinessDao(dbServices.GetClient())

	// Get the BusinessDetails
	dataBusiness, err := daoPlatformBusiness.Get(businessId)
	if err != nil {
		err := &utils.AppError{ErrorCode: funcode + "01", ErrorMsg: "Invalid business_id", ErrorDetail: "Given business_id is not exist"}
		return nil, err
	}

	// Get RegionId from PlatformBusiness
	regionId, err := utils.GetMemberDataStr(dataBusiness, platform_common.FLD_BUSINESS_REGION_ID)
	if err != nil {
		log.Println("GetTenantDBInfo:: RegionId not found in Platform_business", err)
		return nil, err
	}

	// Open Instance of Region
	daoRegion := platform_repository.NewRegionDao(dbServices.GetClient())
	dataRegion, err := daoRegion.Get(regionId)
	if err != nil {
		log.Println("GetTenantDBInfo:: No such region found", regionId, err)
		return nil, err
	}

	// Get all the Database information from the Region
	dbType, _ := utils.GetMemberDataInt(dataRegion, platform_common.FLD_REGION_DB_TYPE, true)
	dbServer, _ := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_MONGODB_SERVER)
	dbUser, _ := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_MONGODB_USER)
	dbSecret, _ := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_MONGODB_SECRET)
	dbName, _ := utils.GetMemberDataStr(dataRegion, platform_common.FLD_REGION_MONGODB_NAME)

	// Check whether the business has the tenant database enabled
	isTenantDB, _ := utils.GetMemberDataBool(dataBusiness, platform_common.FLD_BUSINESS_IS_TENANT_DB)
	if isTenantDB {
		dbName = dbName + "-" + businessId
	}

	// Create Region's Database Props
	regionDBProps := utils.Map{
		db_common.DB_TYPE:               db_common.DatabaseType(dbType),
		db_common.DB_SERVER:             dbServer,
		db_common.DB_USER:               dbUser,
		db_common.DB_SECRET:             dbSecret,
		db_common.DB_NAME:               dbName,
		platform_common.FLD_BUSINESS_ID: businessId,
	}

	return regionDBProps, nil

}
