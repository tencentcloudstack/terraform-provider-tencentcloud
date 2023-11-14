package tencentcloud

import (
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"log"
)

type CdbService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdbService) DescribeCdbBackupDatabasesByFilter(ctx context.Context, param map[string]interface{}) (backupDatabases []*cdb.DatabaseName, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeBackupDatabasesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "SearchDatabase" {
			request.SearchDatabase = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*int64)
		}
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*int64)
		}
		if k == "Items" {
			request.Items = v.([]*cdb.DatabaseName)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeBackupDatabases(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		backupDatabases = append(backupDatabases, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbBackupOverviewByFilter(ctx context.Context, param map[string]interface{}) (backupOverview []*cdb.DescribeBackupOverviewResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeBackupOverviewRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeBackupOverview(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.BackupCount) < 1 {
			break
		}
		backupOverview = append(backupOverview, response.Response.BackupCount...)
		if len(response.Response.BackupCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbBackupSummariesByFilter(ctx context.Context, param map[string]interface{}) (backupSummaries []*cdb.DescribeBackupSummariesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeBackupSummariesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderDirection" {
			request.OrderDirection = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeBackupSummaries(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		backupSummaries = append(backupSummaries, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbBackupTablesByFilter(ctx context.Context, param map[string]interface{}) (backupTables []*cdb.TableName, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeBackupTablesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "DatabaseName" {
			request.DatabaseName = v.(*string)
		}
		if k == "SearchTable" {
			request.SearchTable = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*int64)
		}
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*int64)
		}
		if k == "Items" {
			request.Items = v.([]*cdb.TableName)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeBackupTables(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		backupTables = append(backupTables, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbBinLogByFilter(ctx context.Context, param map[string]interface{}) (binLog []*cdb.DescribeBinlogsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeBinlogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeBinlogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		binLog = append(binLog, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbBinlogBackupOverviewByFilter(ctx context.Context, param map[string]interface{}) (binlogBackupOverview []*cdb.DescribeBinlogBackupOverviewResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeBinlogBackupOverviewRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeBinlogBackupOverview(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.BinlogBackupVolume) < 1 {
			break
		}
		binlogBackupOverview = append(binlogBackupOverview, response.Response.BinlogBackupVolume...)
		if len(response.Response.BinlogBackupVolume) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbCloneListByFilter(ctx context.Context, param map[string]interface{}) (cloneList []*cdb.DescribeCloneListResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeCloneListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeCloneList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		cloneList = append(cloneList, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbDataBackupOverviewByFilter(ctx context.Context, param map[string]interface{}) (dataBackupOverview []*cdb.DescribeDataBackupOverviewResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDataBackupOverviewRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDataBackupOverview(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DataBackupVolume) < 1 {
			break
		}
		dataBackupOverview = append(dataBackupOverview, response.Response.DataBackupVolume...)
		if len(response.Response.DataBackupVolume) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbDatabasesByFilter(ctx context.Context, param map[string]interface{}) (databases []*cdb.DescribeDatabasesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDatabasesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*int64)
		}
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "DatabaseRegexp" {
			request.DatabaseRegexp = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*int64)
		}
		if k == "Items" {
			request.Items = v.([]*string)
		}
		if k == "DatabaseList" {
			request.DatabaseList = v.([]*cdb.DatabasesWithCharacterLists)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDatabases(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		databases = append(databases, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbDbFeaturesByFilter(ctx context.Context, param map[string]interface{}) (dbFeatures []*cdb.DescribeDBFeaturesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDBFeaturesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDBFeatures(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.IsSupportAudit) < 1 {
			break
		}
		dbFeatures = append(dbFeatures, response.Response.IsSupportAudit...)
		if len(response.Response.IsSupportAudit) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbErrorLogByFilter(ctx context.Context, param map[string]interface{}) (errorLog []*cdb.DescribeErrorLogDataResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeErrorLogDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*uint64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*uint64)
		}
		if k == "KeyWords" {
			request.KeyWords = v.([]*string)
		}
		if k == "InstType" {
			request.InstType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeErrorLogData(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		errorLog = append(errorLog, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbInstTablesByFilter(ctx context.Context, param map[string]interface{}) (instTables []*cdb.DescribeTablesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeTablesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Database" {
			request.Database = v.(*string)
		}
		if k == "TableRegexp" {
			request.TableRegexp = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeTables(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		instTables = append(instTables, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbInstanceCharsetByFilter(ctx context.Context, param map[string]interface{}) (instanceCharset []*cdb.DescribeDBInstanceCharsetResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDBInstanceCharsetRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDBInstanceCharset(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Charset) < 1 {
			break
		}
		instanceCharset = append(instanceCharset, response.Response.Charset...)
		if len(response.Response.Charset) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbInstanceInfoByFilter(ctx context.Context, param map[string]interface{}) (instanceInfo []*cdb.DescribeDBInstanceInfoResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDBInstanceInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDBInstanceInfo(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceId) < 1 {
			break
		}
		instanceInfo = append(instanceInfo, response.Response.InstanceId...)
		if len(response.Response.InstanceId) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbInstanceParamRecordByFilter(ctx context.Context, param map[string]interface{}) (instanceParamRecord []*cdb.DescribeInstanceParamRecordsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeInstanceParamRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeInstanceParamRecords(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		instanceParamRecord = append(instanceParamRecord, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbInstanceRebootTimeByFilter(ctx context.Context, param map[string]interface{}) (instanceRebootTime []*cdb.DescribeDBInstanceRebootTimeResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDBInstanceRebootTimeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDBInstanceRebootTime(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		instanceRebootTime = append(instanceRebootTime, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbParamTemplatesByFilter(ctx context.Context, param map[string]interface{}) (paramTemplates []*cdb.ParamTemplateInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeParamTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "EngineVersions" {
			request.EngineVersions = v.([]*string)
		}
		if k == "EngineTypes" {
			request.EngineTypes = v.([]*string)
		}
		if k == "TemplateNames" {
			request.TemplateNames = v.([]*string)
		}
		if k == "TemplateIds" {
			request.TemplateIds = v.([]*int64)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*int64)
		}
		if k == "Items" {
			request.Items = v.([]*cdb.ParamTemplateInfo)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeParamTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		paramTemplates = append(paramTemplates, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbProjectSecurityGroupByFilter(ctx context.Context, param map[string]interface{}) (projectSecurityGroup []*cdb.DescribeProjectSecurityGroupsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeProjectSecurityGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProjectId" {
			request.ProjectId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeProjectSecurityGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		projectSecurityGroup = append(projectSecurityGroup, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbRoMinScaleByFilter(ctx context.Context, param map[string]interface{}) (roMinScale []*cdb.DescribeRoMinScaleResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeRoMinScaleRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "RoInstanceId" {
			request.RoInstanceId = v.(*string)
		}
		if k == "MasterInstanceId" {
			request.MasterInstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeRoMinScale(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Memory) < 1 {
			break
		}
		roMinScale = append(roMinScale, response.Response.Memory...)
		if len(response.Response.Memory) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbRollbackRangeTimeByFilter(ctx context.Context, param map[string]interface{}) (rollbackRangeTime []*cdb.DescribeRollbackRangeTimeResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeRollbackRangeTimeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
		if k == "IsRemoteZone" {
			request.IsRemoteZone = v.(*string)
		}
		if k == "BackupRegion" {
			request.BackupRegion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeRollbackRangeTime(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		rollbackRangeTime = append(rollbackRangeTime, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbSlowLogByFilter(ctx context.Context, param map[string]interface{}) (slowLog []*cdb.DescribeSlowLogsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeSlowLogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeSlowLogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		slowLog = append(slowLog, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbSlowLogDataByFilter(ctx context.Context, param map[string]interface{}) (slowLogData []*cdb.DescribeSlowLogDataResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeSlowLogDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*uint64)
		}
		if k == "EndTime" {
			request.EndTime = v.(*uint64)
		}
		if k == "UserHosts" {
			request.UserHosts = v.([]*string)
		}
		if k == "UserNames" {
			request.UserNames = v.([]*string)
		}
		if k == "DataBases" {
			request.DataBases = v.([]*string)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "InstType" {
			request.InstType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeSlowLogData(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		slowLogData = append(slowLogData, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbSupportedPrivilegesByFilter(ctx context.Context, param map[string]interface{}) (supportedPrivileges []*cdb.DescribeSupportedPrivilegesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeSupportedPrivilegesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeSupportedPrivileges(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.GlobalSupportedPrivileges) < 1 {
			break
		}
		supportedPrivileges = append(supportedPrivileges, response.Response.GlobalSupportedPrivileges...)
		if len(response.Response.GlobalSupportedPrivileges) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbSwitchRecordByFilter(ctx context.Context, param map[string]interface{}) (switchRecord []*cdb.DescribeDBSwitchRecordsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeDBSwitchRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDBSwitchRecords(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		switchRecord = append(switchRecord, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbTablesByFilter(ctx context.Context, param map[string]interface{}) (tables []*cdb.DescribeTablesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeTablesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Database" {
			request.Database = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*int64)
		}
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "TableRegexp" {
			request.TableRegexp = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*int64)
		}
		if k == "Items" {
			request.Items = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeTables(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		tables = append(tables, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbUploadedFilesByFilter(ctx context.Context, param map[string]interface{}) (uploadedFiles []*cdb.DescribeUploadedFilesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cdb.NewDescribeUploadedFilesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Path" {
			request.Path = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeUploadedFiles(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TotalCount) < 1 {
			break
		}
		uploadedFiles = append(uploadedFiles, response.Response.TotalCount...)
		if len(response.Response.TotalCount) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CdbService) DescribeCdbAccountsById(ctx context.Context, instanceId string) (accounts *cdb.AccountInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeAccountsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AccountInfo) < 1 {
		return
	}

	accounts = response.Response.AccountInfo[0]
	return
}

func (me *CdbService) DeleteCdbAccountsById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDeleteAccountsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DeleteAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdbService) CdbAccountsStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAsyncRequestInfo(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbAuditLogFileById(ctx context.Context, instanceId string, fileName string) (auditLogFile *cdb.AuditLogFile, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeAuditLogFilesRequest()
	request.InstanceId = &instanceId
	request.FileName = &fileName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeAuditLogFiles(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AuditLogFile) < 1 {
		return
	}

	auditLogFile = response.Response.AuditLogFile[0]
	return
}

func (me *CdbService) DeleteCdbAuditLogFileById(ctx context.Context, instanceId string, fileName string) (errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDeleteAuditLogFileRequest()
	request.InstanceId = &instanceId
	request.FileName = &fileName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DeleteAuditLogFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdbService) CdbAuditLogFileStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAuditLogFiles(ctx, instanceId, fileName)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbBackupDownloadRestrictionById(ctx context.Context, idsHash string) (backupDownloadRestriction *cdb.DescribeBackupDownloadRestrictionResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeBackupDownloadRestrictionRequest()
	request.IdsHash = &idsHash

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeBackupDownloadRestriction(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	backupDownloadRestriction = response.Response
	return
}

func (me *CdbService) DescribeCdbBackupEncryptionStatusById(ctx context.Context, instanceId string) (backupEncryptionStatus *cdb.DescribeBackupEncryptionStatusResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeBackupEncryptionStatusRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeBackupEncryptionStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	backupEncryptionStatus = response.Response
	return
}

func (me *CdbService) DescribeCdbDbImportById(ctx context.Context, instanceId string) (dbImport *cdb.DescribeDBImportRecordsResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDBImportRecordsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeDBImportRecords(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dbImport = response.Response
	return
}

func (me *CdbService) DeleteCdbDbImportById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewStopDBImportJobRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().StopDBImportJob(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdbService) CdbDbImportStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAsyncRequestInfo(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbDeployGroupById(ctx context.Context, deployGroupId string) (deployGroup *cdb.DeployGroupInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDeployGroupListRequest()
	request.DeployGroupId = &deployGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	instances := make([]*cdb.DeployGroupInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCdbClient().DescribeDeployGroupList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DeployGroupInfo) < 1 {
			break
		}
		instances = append(instances, response.Response.DeployGroupInfo...)
		if len(response.Response.DeployGroupInfo) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	deployGroup = instances[0]
	return
}

func (me *CdbService) DeleteCdbDeployGroupById(ctx context.Context, deployGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDeleteDeployGroupsRequest()
	request.DeployGroupId = &deployGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DeleteDeployGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdbService) DescribeCdbInstanceEncryptionById(ctx context.Context, instanceId string) (instanceEncryption *cdb.DescribeDBInstanceInfoResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDBInstanceInfoRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeDBInstanceInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceEncryption = response.Response
	return
}

func (me *CdbService) CdbInstanceEncryptionStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeDBInstanceInfo(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Encryption), nil
	}
}

func (me *CdbService) DescribeCdbInstanceParamById(ctx context.Context, instanceId string) (instanceParam *cdb.ParameterDetail, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeInstanceParamsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ParameterDetail) < 1 {
		return
	}

	instanceParam = response.Response.ParameterDetail[0]
	return
}

func (me *CdbService) DescribeCdbInstanceTypeById(ctx context.Context, instanceId string) (instanceType *cdb.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceInfo) < 1 {
		return
	}

	instanceType = response.Response.InstanceInfo[0]
	return
}

func (me *CdbService) CdbInstanceTypeStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAsyncRequestInfo(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbLocalBinlogConfigById(ctx context.Context, instanceId string) (localBinlogConfig *cdb.LocalBinlogConfig, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeLocalBinlogConfigRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeLocalBinlogConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.LocalBinlogConfig) < 1 {
		return
	}

	localBinlogConfig = response.Response.LocalBinlogConfig[0]
	return
}

func (me *CdbService) DescribeCdbParamTemplateById(ctx context.Context, templateId string) (paramTemplate *cdb.DescribeParamTemplateInfoResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeParamTemplateInfoRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeParamTemplateInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	paramTemplate = response.Response
	return
}

func (me *CdbService) DeleteCdbParamTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDeleteParamTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DeleteParamTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdbService) DescribeCdbPasswordComplexityById(ctx context.Context, instanceId string) (passwordComplexity *cdb.ParameterDetail, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeInstanceParamsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ParameterDetail) < 1 {
		return
	}

	passwordComplexity = response.Response.ParameterDetail[0]
	return
}

func (me *CdbService) CdbPasswordComplexityStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAsyncRequestInfo(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbRemoteBackupConfigById(ctx context.Context, instanceId string) (remoteBackupConfig *cdb.DescribeRemoteBackupConfigResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeRemoteBackupConfigRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeRemoteBackupConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	remoteBackupConfig = response.Response
	return
}

func (me *CdbService) DescribeCdbRestartDBInstancesById(ctx context.Context, idsHash string) (restartDBInstances *cdb.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDBInstancesRequest()
	request.IdsHash = &idsHash

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceInfo) < 1 {
		return
	}

	restartDBInstances = response.Response.InstanceInfo[0]
	return
}

func (me *CdbService) CdbRestartDBInstancesStateRefreshFunc(idsHash string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeCdbProxyInfo(ctx, idsHash)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.TaskStatus), nil
	}
}

func (me *CdbService) DescribeCdbRoGroupById(ctx context.Context, instanceId string, roGroupId string) (roGroup *cdb.RoGroup, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeRoGroupsRequest()
	request.InstanceId = &instanceId
	request.RoGroupId = &roGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeRoGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RoGroup) < 1 {
		return
	}

	roGroup = response.Response.RoGroup[0]
	return
}

func (me *CdbService) CdbRoGroupStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAsyncRequestInfo(ctx, instanceId, roGroupId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbRollbackById(ctx context.Context, instanceId string) (rollback *cdb.DescribeRollbackTaskDetailResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeRollbackTaskDetailRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeRollbackTaskDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	rollback = response.Response
	return
}

func (me *CdbService) CdbRollbackStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeAsyncRequestInfo(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.Status), nil
	}
}

func (me *CdbService) DescribeCdbSecurityGroupsAttachmentById(ctx context.Context, securityGroupId string, instanceId string) (securityGroupsAttachment *cdb.SecurityGroup, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDBSecurityGroupsRequest()
	request.SecurityGroupId = &securityGroupId
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SecurityGroup) < 1 {
		return
	}

	securityGroupsAttachment = response.Response.SecurityGroup[0]
	return
}

func (me *CdbService) DeleteCdbSecurityGroupsAttachmentById(ctx context.Context, securityGroupId string, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDisassociateSecurityGroupsRequest()
	request.SecurityGroupId = &securityGroupId
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DisassociateSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CdbService) DescribeCdbSwitchForUpgradeById(ctx context.Context, instanceId string) (switchForUpgrade *cdb.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceInfo) < 1 {
		return
	}

	switchForUpgrade = response.Response.InstanceInfo[0]
	return
}

func (me *CdbService) CdbSwitchForUpgradeStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeDBInstances(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.TaskStatus), nil
	}
}

func (me *CdbService) DescribeCdbTimeWindowById(ctx context.Context, instanceId string) (timeWindow *cdb.DescribeTimeWindowResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeTimeWindowRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdbClient().DescribeTimeWindow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	timeWindow = response.Response
	return
}
