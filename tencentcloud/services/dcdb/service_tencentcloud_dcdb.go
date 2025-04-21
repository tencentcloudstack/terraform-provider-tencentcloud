package dcdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewDcdbService(client *connectivity.TencentCloudClient) DcdbService {
	return DcdbService{client: client}
}

type DcdbService struct {
	client *connectivity.TencentCloudClient
}

// dc_account
func (me *DcdbService) DescribeDcdbAccount(ctx context.Context, instanceId, userName string) (account *dcdb.DescribeAccountsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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

// dc_db_instance
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
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(15*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		dbInstances, errResp := me.DescribeDcdbDbInstance(ctx, instanceId)
		if errResp != nil {
			return tccommon.RetryError(errResp, tccommon.InternalError)
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
			initErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := me.client.UseDcdbClient().InitDCDBInstances(iniRequest)
				if e != nil {
					return tccommon.RetryError(e)
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

// dc_hourdb_instance
func (me *DcdbService) DescribeDcdbHourdbInstance(ctx context.Context, instanceId string) (hourdbInstance *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	return me.DescribeDcdbDbInstance(ctx, instanceId)
}

func (me *DcdbService) DeleteDcdbHourdbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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

// dc_sg
func (me *DcdbService) DescribeDcdbSecurityGroup(ctx context.Context, instanceId string) (securityGroup *dcdb.DescribeDBSecurityGroupsResponseParams, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = dcdb.NewDescribeDBSecurityGroupsRequest()
		response = dcdb.NewDescribeDBSecurityGroupsResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.Product = helper.String("dcdb") // api only use this fixed value
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseDcdbClient().DescribeDBSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	securityGroup = response.Response

	return
}

func (me *DcdbService) DeleteDcdbSecurityGroupAttachmentById(ctx context.Context, instanceId, securityGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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
		logId    = tccommon.GetLogId(ctx)
		request  = dcdb.NewDescribeDCDBInstancesRequest()
		response = dcdb.NewDescribeDCDBInstancesResponse()
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
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseDcdbClient().DescribeDCDBInstances(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}

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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

	request := dcdb.NewDescribeDcnDetailRequest()
	response := dcdb.NewDescribeDcnDetailResponse()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseDcdbClient().DescribeDcnDetail(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

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
	logId := tccommon.GetLogId(ctx)

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
		ctx := tccommon.ContextNil

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
	logId := tccommon.GetLogId(ctx)

	request := dcdb.NewDescribeAccountPrivilegesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	idSplit := strings.Split(ids, tccommon.FILED_SP)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		logId      = tccommon.GetLogId(ctx)
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
		logId      = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

	request := dcdb.NewDescribeDCDBInstanceDetailRequest()
	response := dcdb.NewDescribeDCDBInstanceDetailResponse()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseDcdbClient().DescribeDCDBInstanceDetail(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	dedicatedClusterDbInstance = response.Response
	return
}

func (me *DcdbService) DeleteDcdbDedicatedClusterDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		ctx := tccommon.ContextNil

		object, err := me.DescribeDcdbDbSyncModeConfigById(ctx, instanceId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.Int64ToStr(*object.IsModifying), nil
	}
}

func (me *DcdbService) DcdbDcnStateRefreshFunc(instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil
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
		ctx := tccommon.ContextNil

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
	logId := tccommon.GetLogId(ctx)
	var flowId *int64

	if enable {
		request := dcdb.NewOpenDBExtranetAccessRequest()
		request.InstanceId = &instanceId
		request.Ipv6Flag = helper.IntInt64(ipv6Flag)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseDcdbClient().OpenDBExtranetAccess(request)
			if e != nil {
				return tccommon.RetryError(e)
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

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := me.client.UseDcdbClient().CloseDBExtranetAccess(request)
			if e != nil {
				return tccommon.RetryError(e)
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
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 2*tccommon.ReadRetryTimeout, time.Second, me.DcdbDbInstanceStateRefreshFunc(flowId, []string{"1"}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}
	return
}

func (me *DcdbService) SetRealServerAccessStrategy(ctx context.Context, instanceId string, rsAccessStrategy int) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := dcdb.NewModifyRealServerAccessStrategyRequest()
	request.InstanceId = &instanceId
	request.RsAccessStrategy = helper.IntInt64(rsAccessStrategy)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseDcdbClient().ModifyRealServerAccessStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	logId := tccommon.GetLogId(ctx)
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseDcdbClient().ModifyInstanceNetwork(request)
		if e != nil {
			return tccommon.RetryError(e)
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
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 2*tccommon.ReadRetryTimeout, time.Second, me.DcdbDbInstanceStateRefreshFunc(flowId, []string{"1"}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	return
}

func (me *DcdbService) DescribeDcdbFileDownloadUrlByFilter(ctx context.Context, param map[string]interface{}) (fileDownloadUrl *string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeFileDownloadUrlRequest()
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
		if k == "ShardId" {
			request.ShardId = v.(*string)
		}
		if k == "FilePath" {
			request.FilePath = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeFileDownloadUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	fileDownloadUrl = response.Response.PreSignedUrl

	return
}

func (me *DcdbService) DescribeDcdbLogFilesByFilter(ctx context.Context, param map[string]interface{}) (ret *dcdb.DescribeDBLogFilesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDBLogFilesRequest()
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
		if k == "ShardId" {
			request.ShardId = v.(*string)
		}
		if k == "Type" {
			request.Type = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDBLogFiles(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	ret = response.Response

	return
}

func (me *DcdbService) DescribeDcdbInstanceNodeInfoByFilter(ctx context.Context, param map[string]interface{}) (instanceNodeInfo []*dcdb.BriefNodeInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDCDBInstanceNodeInfoRequest()
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
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDcdbClient().DescribeDCDBInstanceNodeInfo(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NodesInfo) < 1 {
			break
		}
		instanceNodeInfo = append(instanceNodeInfo, response.Response.NodesInfo...)
		if len(response.Response.NodesInfo) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DcdbService) DescribeDcdbOrdersByFilter(ctx context.Context, param map[string]interface{}) (orders []*dcdb.Deal, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeOrdersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DealNames" {
			request.DealNames = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeOrders(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Deals) < 1 {
		return
	}

	orders = response.Response.Deals

	return
}

func (me *DcdbService) DescribeDcdbPriceByFilter(ctx context.Context, param map[string]interface{}) (ret *dcdb.DescribeDCDBPriceResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDCDBPriceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceCount" {
			request.Count = v.(*int64)
		}
		if k == "Zone" {
			request.Zone = v.(*string)
		}
		if k == "Period" {
			request.Period = v.(*int64)
		}
		if k == "ShardNodeCount" {
			request.ShardNodeCount = v.(*int64)
		}
		if k == "ShardMemory" {
			request.ShardMemory = v.(*int64)
		}
		if k == "ShardStorage" {
			request.ShardStorage = v.(*int64)
		}
		if k == "ShardCount" {
			request.ShardCount = v.(*int64)
		}
		if k == "Paymode" {
			request.Paymode = v.(*string)
		}
		if k == "AmountUnit" {
			request.AmountUnit = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDCDBPrice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	ret = response.Response

	return
}

func (me *DcdbService) DescribeDcdbProjectsByFilter(ctx context.Context) (projects []*dcdb.Project, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeProjects(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Projects) < 1 {
		return
	}
	projects = response.Response.Projects

	return
}

func (me *DcdbService) DescribeDcdbProjectSecurityGroupsByFilter(ctx context.Context, param map[string]interface{}) (projectSecurityGroups []*dcdb.SecurityGroup, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeProjectSecurityGroupsRequest()
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
		if k == "ProjectId" {
			request.ProjectId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeProjectSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Groups) < 1 {
		return
	}
	projectSecurityGroups = response.Response.Groups

	return
}

func (me *DcdbService) DescribeDcdbRenewalPriceByFilter(ctx context.Context, param map[string]interface{}) (ret *dcdb.DescribeDCDBRenewalPriceResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDCDBRenewalPriceRequest()
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
		if k == "Period" {
			request.Period = v.(*int64)
		}
		if k == "AmountUnit" {
			request.AmountUnit = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDCDBRenewalPrice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	ret = response.Response

	return
}

func (me *DcdbService) DescribeDcdbSaleInfoByFilter(ctx context.Context, param map[string]interface{}) (regionInfo []*dcdb.RegionInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDCDBSaleInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDCDBSaleInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.RegionList) < 1 {
		return
	}

	regionInfo = response.Response.RegionList

	return
}

func (me *DcdbService) DescribeDcdbShardSpecByFilter(ctx context.Context, param map[string]interface{}) (specConfigs []*dcdb.SpecConfig, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeShardSpecRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseDcdbClient().DescribeShardSpec(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.SpecConfig) < 1 {
		return
	}

	specConfigs = response.Response.SpecConfig

	return
}

func (me *DcdbService) DescribeDcdbSlowLogsByFilter(ctx context.Context, param map[string]interface{}) (slowLogs []*dcdb.SlowLogData, ret *dcdb.DescribeDBSlowLogsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDBSlowLogsRequest()
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
		if k == "ShardId" {
			request.ShardId = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "Db" {
			request.Db = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*string)
		}
		if k == "Slave" {
			request.Slave = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseDcdbClient().DescribeDBSlowLogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		slowLogs = append(slowLogs, response.Response.Data...)
		ret = response.Response
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DcdbService) DescribeDcdbUpgradePriceByFilter(ctx context.Context, param map[string]interface{}) (ret *dcdb.DescribeDCDBUpgradePriceResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dcdb.NewDescribeDCDBUpgradePriceRequest()
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
		if k == "UpgradeType" {
			request.UpgradeType = v.(*string)
		}
		if k == "AddShardConfig" {
			request.AddShardConfig = v.(*dcdb.AddShardConfig)
		}
		if k == "ExpandShardConfig" {
			request.ExpandShardConfig = v.(*dcdb.ExpandShardConfig)
		}
		if k == "SplitShardConfig" {
			request.SplitShardConfig = v.(*dcdb.SplitShardConfig)
		}
		if k == "AmountUnit" {
			request.AmountUnit = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDcdbClient().DescribeDCDBUpgradePrice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	ret = response.Response

	return
}
