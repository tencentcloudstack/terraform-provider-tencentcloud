package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func (me *DcdbService) InitDcdbDbInstance(ctx context.Context, instanceId string, params []*dcdb.DBParamValue) (initRet bool, errRet error) {
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
		return false, err
	}
	return true, nil
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
