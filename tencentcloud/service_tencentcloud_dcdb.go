package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DcdbService struct {
	client *connectivity.TencentCloudClient
}

//dc_account
func (me *DcdbService) DescribeDcdbAccount(ctx context.Context, instanceId, userName string) (account *dcdb.DescribeAccountsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeAccountsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseDcdbClient().DescribeAccounts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if userName != "" {
		// filter
		for _, user := range response.Response.Users {
			if *user.UserName == userName {
				account = &dcdb.DescribeAccountsResponseParams{
					InstanceId: response.Response.InstanceId,
					RequestId:  response.Response.RequestId,
					Users:      []*dcdb.DBAccount{user},
				}
				return
			}
		}
		return
	}

	account = response.Response
	return
}

func (me *DcdbService) DeleteDcdbAccountById(ctx context.Context, instanceId, userName, host string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDeleteAccountRequest()

	request.InstanceId = &instanceId
	request.UserName = &userName
	request.Host = &host

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DeleteAccount(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

//dc_db_instance
func (me *DcdbService) DescribeDcdbDbInstance(ctx context.Context, instanceId string) (instances *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	params := make(map[string]interface{})
	params["instance_ids"] = []*string{&instanceId}

	ret, err := me.DescribeDcdbInstancesByFilter(ctx, params)
	if err != nil {
		return nil, err
	}

	result := &dcdb.DescribeDCDBInstancesResponseParams{
		Instances:  ret,
		TotalCount: helper.IntInt64(len(ret)),
	}

	return result, nil
}

func (me *DcdbService) InitDcdbDbInstance(ctx context.Context, instanceId string, params []*dcdb.DBParamValue) (initRet bool, flowId *uint64, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(15*readRetryTimeout, func() *resource.RetryError {
		dbInstances, errResp := me.DescribeDcdbDbInstance(ctx, instanceId)
		if errResp != nil {
			return retryError(errResp, InternalError)
		}
		if dbInstances.Instances[0] == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeDcdbDbInstance return result(dcdb instance) is nil!"))
		}

		dbInstance := dbInstances.Instances[0]
		if *dbInstance.Status < 0 {
			return resource.NonRetryableError(fmt.Errorf("dcdb instance init status is %v, operate failed", *dbInstance.Status))
		}
		if *dbInstance.Status == 2 {
			return nil
		}
		if *dbInstance.Status == 3 {
			iniRequest := dcdb.NewInitDCDBInstancesRequest()
			iniRequest.InstanceIds = []*string{&instanceId}
			iniRequest.Params = params
			initErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := me.client.UseDcdbClient().InitDCDBInstances(iniRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				flowId = result.Response.FlowIds[0]
				return nil
			})
			if initErr != nil {
				return resource.NonRetryableError(fmt.Errorf("dcdb instance init error %v, operate failed", initErr))
			}
			time.Sleep(10 * time.Second)
			return resource.RetryableError(fmt.Errorf("dcdb instance initializing, retry..."))
		}
		return resource.RetryableError(fmt.Errorf("dcdb instance init status is %v, retry...", *dbInstance.Status))
	})
	if err != nil {
		return false, nil, err
	}
	return true, flowId, nil
}

//dc_hourdb_instance
func (me *DcdbService) DescribeDcdbHourdbInstance(ctx context.Context, instanceId string) (hourdbInstance *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	return me.DescribeDcdbDbInstance(ctx, instanceId)
}

func (me *DcdbService) DeleteDcdbHourdbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDestroyHourDCDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DestroyHourDCDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

//dc_sg
func (me *DcdbService) DescribeDcdbSecurityGroup(ctx context.Context, instanceId string) (securityGroup *dcdb.DescribeDBSecurityGroupsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDBSecurityGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.Product = helper.String("dcdb") // api only use this fixed value

	response, err := me.client.UseDcdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	securityGroup = response.Response

	return
}

func (me *DcdbService) DeleteDcdbSecurityGroupAttachmentById(ctx context.Context, instanceId, securityGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDisassociateSecurityGroupsRequest()

	request.Product = helper.String("dcdb") // api only use this fixed value
	request.InstanceIds = []*string{&instanceId}
	request.SecurityGroupId = &securityGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DisassociateSecurityGroups(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// for data_source
// tencentcloud_dcdb_instances
func (me *DcdbService) DescribeDcdbInstancesByFilter(ctx context.Context, params map[string]interface{}) (instances []*dcdb.DCDBInstanceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range params {
		if k == "instance_ids" {
			var ids []*string
			ids = append(ids, v.([]*string)...)
			request.InstanceIds = ids
		}

		if k == "search_name" {
			request.SearchName = v.(*string)
		}

		if k == "search_key" {
			request.SearchKey = v.(*string)
		}

		if k == "project_ids" {
			var ids []*int64
			ids = append(ids, v.([]*int64)...)
			request.ProjectIds = ids
		}

		if k == "is_filter_excluster" {
			request.IsFilterExcluster = v.(*bool)
		}

		if k == "excluster_type" {
			request.ExclusterType = v.(*int64)
		}

		if k == "is_filter_vpc" {
			request.IsFilterVpc = v.(*bool)
		}

		if k == "vpc_id" {
			request.VpcId = v.(*string)
		}

		if k == "subnet_id" {
			request.SubnetId = v.(*string)
		}

	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize

		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDcdbClient().DescribeDCDBInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Instances) < 1 {
			break
		}
		instances = append(instances, response.Response.Instances...)
		if len(response.Response.Instances) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

// tencentcloud_dcdb_accounts
func (me *DcdbService) DescribeDcdbAccountsByFilter(ctx context.Context, param map[string]interface{}) (accounts []*dcdb.DBAccount, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeAccountsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDcdbClient().DescribeAccounts(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Users) < 1 {
			break
		}
		accounts = append(accounts, response.Response.Users...)
		if len(response.Response.Users) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

// tencentcloud_dcdb_databases
func (me *DcdbService) DescribeDcdbDatabasesByFilter(ctx context.Context, param map[string]interface{}) (databases []*dcdb.Database, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDatabasesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDcdbClient().DescribeDatabases(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Databases) < 1 {
			break
		}
		databases = append(databases, response.Response.Databases...)
		if len(response.Response.Databases) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

// tencentcloud_dcdb_parameters
func (me *DcdbService) DescribeDcdbParametersByFilter(ctx context.Context, param map[string]interface{}) (parameters []*dcdb.ParamDesc, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDBParametersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DescribeDBParameters(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	parameters = append(parameters, response.Response.Params...)

	return
}

// tencentcloud_dcdb_shards
func (me *DcdbService) DescribeDcdbShardsByFilter(ctx context.Context, param map[string]interface{}) (shards []*dcdb.DCDBShardInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBShardsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}

		if k == "shard_instance_ids" {
			request.ShardInstanceIds = v.([]*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDcdbClient().DescribeDCDBShards(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Shards) < 1 {
			break
		}
		shards = append(shards, response.Response.Shards...)
		if len(response.Response.Shards) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

// tencentcloud_dcdb_security_groups
func (me *DcdbService) DescribeDcdbSecurityGroupsByFilter(ctx context.Context, param map[string]interface{}) (securityGroups []*dcdb.SecurityGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDBSecurityGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
	}
	request.Product = helper.String("dcdb") // api only use this fixed value

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	securityGroups = append(securityGroups, response.Response.Groups...)

	return
}

// tencentcloud_dcdb_db_instance
func (me *DcdbService) DeleteDcdbDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDestroyDCDBInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DestroyDCDBInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DcdbService) DescribeDcnDetailById(ctx context.Context, instanceId string) (dcnDetails []*dcdb.DcnDetailItem, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDcnDetailRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDcnDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DcnDetails) < 1 {
		return
	}

	// we need this relationship about master and dcn, so no need to filter the results
	// if instanceId != "" {
	// 	for _, detail := range response.Response.DcnDetails {
	// 		if *detail.InstanceId == instanceId {
	// 			dbInstance = detail
	// 			return
	// 		}
	// 	}
	// 	return
	// }

	dcnDetails = response.Response.DcnDetails
	return
}

func (me *DcdbService) DescribeDcdbFlowById(ctx context.Context, flowId *int64) (dbInstance *dcdb.DescribeFlowResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeFlowRequest()
	if flowId != nil {
		request.FlowId = flowId
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeFlow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dbInstance = response.Response
	return
}

func (me *DcdbService) DcdbDbInstanceStateRefreshFunc(flowId *int64, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		if *flowId == 0 {
			return &dcdb.DescribeFlowResponseParams{}, "0", nil
		}

		object, err := me.DescribeDcdbFlowById(ctx, flowId)

		if err != nil {
			return nil, "", err
		}

		for _, str := range failStates {
			if strings.Contains(str, helper.Int64ToStr(*object.Status)) {
				return &dcdb.DescribeFlowResponseParams{}, "1", nil
			}
		}

		return object, helper.Int64ToStr(*object.Status), nil
	}
}

// tencentcloud_dcdb_account_privileges
func (me *DcdbService) DescribeDcdbAccountPrivilegesById(ctx context.Context, ids string, dbName, aType, object, colName *string) (accountPrivileges *dcdb.DescribeAccountPrivilegesResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeAccountPrivilegesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	idSplit := strings.Split(ids, FILED_SP)
	if len(idSplit) != 7 {
		return nil, fmt.Errorf("[service_tc_dbdb]id is broken,%s", ids)
	}

	request.InstanceId = helper.String(idSplit[0])
	request.UserName = helper.String(idSplit[1])
	request.Host = helper.String(idSplit[2])

	all := helper.String("*")

	if dbName != nil {
		request.DbName = dbName
	} else {
		request.DbName = all
	}

	if aType != nil {
		request.Type = aType
	} else {
		request.Type = all
	}

	if object != nil {
		request.Object = object
	} else {
		request.Object = all
	}

	if colName != nil {
		request.ColName = colName
	} else {
		request.ColName = all
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeAccountPrivileges(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accountPrivileges = response.Response
	return
}

func (me *DcdbService) DescribeDcdbDatabases(ctx context.Context, instanceId string) (rets *dcdb.DescribeDatabasesResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDatabasesRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDatabases(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	rets = response.Response
	return
}

func (me *DcdbService) DescribeDcdbDBTables(ctx context.Context, instanceId string, dbName string, tableName string) (rets *dcdb.DescribeDatabaseTableResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDatabaseTableRequest()
	request.InstanceId = &instanceId
	request.DbName = &dbName
	request.Table = &tableName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDatabaseTable(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	rets = response.Response
	return
}

func (me *DcdbService) DescribeDcdbDBObjects(ctx context.Context, instanceId string, dbName string) (rets *dcdb.DescribeDatabaseObjectsResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDatabaseObjectsRequest()
	request.InstanceId = &instanceId
	request.DbName = &dbName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDatabaseObjects(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	rets = response.Response
	return
}

// tencentcloud_dcdb_db_parameters
func (me *DcdbService) DescribeDcdbDbParametersById(ctx context.Context, instanceId string) (dbParameters *dcdb.DescribeDBParametersResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDBParametersRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDBParameters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dbParameters = response.Response
	return
}

// tencentcloud_dcdb_database_objects
func (me *DcdbService) DescribeDcdbDBObjectsByFilter(ctx context.Context, param map[string]interface{}) (response *dcdb.DescribeDatabaseObjectsResponseParams, errRet error) {
	var (
		logId      = getLogId(ctx)
		instanceId *string
		dbName     *string
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s filter api[%s] fail, reason[%s]\n",
				logId, "query db objects", errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			instanceId = v.(*string)
		}
		if k == "db_name" {
			dbName = v.(*string)
		}
	}

	response, errRet = me.DescribeDcdbDBObjects(ctx, *instanceId, *dbName)
	if errRet != nil {
		return
	}

	return
}

// tencentcloud_dcdb_database_tables
func (me *DcdbService) DescribeDcdbDBTablesByFilter(ctx context.Context, param map[string]interface{}) (response *dcdb.DescribeDatabaseTableResponseParams, errRet error) {
	var (
		logId      = getLogId(ctx)
		instanceId *string
		dbName     *string
		tableName  *string
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s filter api[%s] fail, reason[%s]\n",
				logId, "query db tables", errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			instanceId = v.(*string)
		}
		if k == "db_name" {
			dbName = v.(*string)
		}
		if k == "table" {
			tableName = v.(*string)
		}
	}

	response, errRet = me.DescribeDcdbDBTables(ctx, *instanceId, *dbName, *tableName)
	if errRet != nil {
		return
	}

	return
}

func (me *DcdbService) DescribeDcdbDbInstanceDetailById(ctx context.Context, instanceId string) (dedicatedClusterDbInstance *dcdb.DescribeDCDBInstanceDetailResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDCDBInstanceDetailRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDCDBInstanceDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dedicatedClusterDbInstance = response.Response
	return
}

func (me *DcdbService) DeleteDcdbDedicatedClusterDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewTerminateDedicatedDBInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().TerminateDedicatedDBInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DcdbService) DescribeDcdbEncryptAttributesConfigById(ctx context.Context, instanceId string) (encryptAttributesConfig *dcdb.DescribeDBEncryptAttributesResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDBEncryptAttributesRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDBEncryptAttributes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	encryptAttributesConfig = response.Response
	return
}

func (me *DcdbService) DescribeDcdbDbSyncModeConfigById(ctx context.Context, instanceId string) (dbSyncModeConfig *dcdb.DescribeDBSyncModeResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDescribeDBSyncModeRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDBSyncMode(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dbSyncModeConfig = response.Response
	return
}

func (me *DcdbService) DcdbDbSyncModeConfigStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeDcdbDbSyncModeConfigById(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.Int64ToStr(*object.IsModifying), nil
	}
}

func (me *DcdbService) DcdbDcnStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil
		rets, err := me.DescribeDcnDetailById(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		for _, object := range rets {
			if *object.InstanceId == instanceId {
				return object, helper.Int64ToStr(*object.DcnStatus), nil
			}
		}
		return &dcdb.DcnDetailItem{}, "0", nil
	}
}

func (me *DcdbService) DcdbAccountRefreshFunc(instanceId string, userName string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeDcdbAccount(ctx, instanceId, userName)

		if err != nil {
			return nil, "", err
		}

		if object == nil || len(object.Users) < 1 {
			return &dcdb.DBAccount{}, "deleted", nil
		}

		user := object.Users[0]
		return user, *user.UserName, nil
	}
}

func (me *DcdbService) SetDcdbExtranetAccess(ctx context.Context, instanceId string, ipv6Flag int, enable bool) (errRet error) {
	logId := getLogId(ctx)
	var flowId *int64

	if enable {
		request := dcdb.NewOpenDBExtranetAccessRequest()
		request.InstanceId = &instanceId
		request.Ipv6Flag = helper.IntInt64(ipv6Flag)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseDcdbClient().OpenDBExtranetAccess(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			flowId = result.Response.FlowId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate dcdb openDBExtranetAccessOperation failed, reason:%+v", logId, err)
			errRet = err
			return
		}

	} else {
		request := dcdb.NewCloseDBExtranetAccessRequest()
		request.InstanceId = &instanceId
		request.Ipv6Flag = helper.IntInt64(ipv6Flag)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseDcdbClient().CloseDBExtranetAccess(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			flowId = result.Response.FlowId
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate dcdb closeDBExtranetAccessOperation failed, reason:%+v", logId, err)
			errRet = err
			return
		}
	}

	if flowId != nil {
		// need to wait operation complete
		// 0:success; 1:failed, 2:running
		conf := BuildStateChangeConf([]string{}, []string{"0"}, 2*readRetryTimeout, time.Second, me.DcdbDbInstanceStateRefreshFunc(flowId, []string{"1"}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}
	return
}

func (me *DcdbService) SetRealServerAccessStrategy(ctx context.Context, instanceId string, rsAccessStrategy int) (errRet error) {
	logId := getLogId(ctx)
	request := dcdb.NewModifyRealServerAccessStrategyRequest()
	request.InstanceId = &instanceId
	request.RsAccessStrategy = helper.IntInt64(rsAccessStrategy)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseDcdbClient().ModifyRealServerAccessStrategy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb modifyRealServerAccessStrategyOperation failed, reason:%+v", logId, err)
		return err
	}

	return
}

func (me *DcdbService) SetNetworkVip(ctx context.Context, instanceId, vpcId, subnetId, vip, vipv6 string) (errRet error) {
	logId := getLogId(ctx)
	var flowId *int64

	request := dcdb.NewModifyInstanceNetworkRequest()

	request.InstanceId = &instanceId
	request.VpcId = &vpcId
	request.SubnetId = &subnetId
	if vip != "" {
		request.Vip = &vip
	}
	if vipv6 != "" {
		request.Vipv6 = &vipv6
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseDcdbClient().ModifyInstanceNetwork(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		flowId = result.Response.FlowId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb modifyInstanceNetworkOperation failed, reason:%+v", logId, err)
		return err
	}

	if flowId != nil {
		// need to wait operation complete
		// 0:success; 1:failed, 2:running
		conf := BuildStateChangeConf([]string{}, []string{"0"}, 2*readRetryTimeout, time.Second, me.DcdbDbInstanceStateRefreshFunc(flowId, []string{"1"}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	return
}
