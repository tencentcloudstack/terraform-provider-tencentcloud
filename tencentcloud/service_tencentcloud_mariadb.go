package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MariadbService struct {
	client *connectivity.TencentCloudClient
}

func (me *MariadbService) InitDbInstance(ctx context.Context, instanceId string, params []*mariadb.DBParamValue) (initRet bool, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(8*readRetryTimeout, func() *resource.RetryError {
		dbInstance, errResp := me.DescribeMariadbDbInstance(ctx, instanceId)
		if errResp != nil {
			return retryError(errResp, InternalError)
		}
		if *dbInstance.Status < 0 {
			return resource.NonRetryableError(fmt.Errorf("db instance init status is %v, operate failed", *dbInstance.Status))
		}
		if *dbInstance.Status == 2 {
			return nil
		}
		if *dbInstance.Status == 3 {
			iniRequest := mariadb.NewInitDBInstancesRequest()
			iniRequest.InstanceIds = []*string{&instanceId}
			iniRequest.Params = params
			initErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := me.client.UseMariadbClient().InitDBInstances(iniRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if initErr != nil {
				return resource.NonRetryableError(fmt.Errorf("db instance init error %v, operate failed", initErr))
			}
			return resource.RetryableError(fmt.Errorf("db instance initializing, retry..."))
		}
		return resource.RetryableError(fmt.Errorf("db instance init status is %v, retry...", *dbInstance.Status))
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (me *MariadbService) DescribeMariadbDbInstancesByFilter(ctx context.Context, param map[string]interface{}) (dbInstances []*mariadb.DBInstance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_ids" {
			request.InstanceIds = v.([]*string)
		}

		if k == "project_ids" {
			request.ProjectIds = v.([]*int64)
		}

		if k == "search_name" {
			request.SearchName = v.(*string)
		}

		if k == "vpc_id" {
			request.VpcId = v.(*string)
		}

		if k == "subnet_id" {
			request.SubnetId = v.(*string)
		}

		if k == "excluster_type" {
			request.IsFilterExcluster = helper.Bool(true)
			request.ExclusterType = v.(*int64)
		}

		if k == "excluster_ids" {
			request.ExclusterIds = v.([]*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseMariadbClient().DescribeDBInstances(request)
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
		dbInstances = append(dbInstances, response.Response.Instances...)
		if len(response.Response.Instances) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *MariadbService) DescribeMariadbDbInstance(ctx context.Context, instanceId string) (dbInstance *mariadb.DBInstance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceIds = []*string{&instanceId}

	response, err := me.client.UseMariadbClient().DescribeDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || len(response.Response.Instances) < 1 {
		return
	}
	dbInstance = response.Response.Instances[0]
	return
}

func (me *MariadbService) DescribeMariadbDbInstanceDetail(ctx context.Context, instanceId string) (dbInstanceDetail *mariadb.DescribeDBInstanceDetailResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBInstanceDetailRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId

	response, err := me.client.UseMariadbClient().DescribeDBInstanceDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	dbInstanceDetail = response.Response
	return
}

func (me *MariadbService) DeleteMariadbDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDestroyDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DestroyDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DeleteMariadbHourDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDestroyHourDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DestroyHourDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DeleteMariadbDbInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewTerminateDedicatedDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().TerminateDedicatedDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DescribeMariadbAccount(ctx context.Context, instanceId, userName, host string) (account *mariadb.DBAccount, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeAccountsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	// request.UserName = &userName

	response, err := me.client.UseMariadbClient().DescribeAccounts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && len(response.Response.Users) > 0 {
		for _, v := range response.Response.Users {
			if *v.UserName == userName && *v.Host == host {
				account = v
				break
			}
		}
	}
	return
}

func (me *MariadbService) DeleteMariadbAccountById(ctx context.Context, instanceId, userName, host string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDeleteAccountRequest()

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
	response, err := me.client.UseMariadbClient().DeleteAccount(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DescribeMariadbParameters(ctx context.Context, instanceId string) (parameters *mariadb.DescribeDBParametersResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBParametersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseMariadbClient().DescribeDBParameters(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	parameters = response.Response
	return
}

func (me *MariadbService) DescribeMariadbLogFileRetentionPeriod(ctx context.Context, instanceId string) (logFileRetentionPeriod *mariadb.DescribeLogFileRetentionPeriodResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeLogFileRetentionPeriodRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseMariadbClient().DescribeLogFileRetentionPeriod(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	logFileRetentionPeriod = response.Response
	return
}

func (me *MariadbService) DescribeMariadbSecurityGroup(ctx context.Context, instanceId, securityGroupId, product string) (securityGroup *mariadb.SecurityGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBSecurityGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Product = &product

	response, err := me.client.UseMariadbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && len(response.Response.Groups) > 0 {
		for _, v := range response.Response.Groups {
			if *v.SecurityGroupId == securityGroupId {
				securityGroup = v
				break
			}
		}
	}
	return
}

func (me *MariadbService) DeleteMariadbSecurityGroupsById(ctx context.Context, instanceId, securityGroupId, product string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDisassociateSecurityGroupsRequest()

	request.InstanceIds = []*string{&instanceId}
	request.SecurityGroupId = &securityGroupId
	request.Product = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DisassociateSecurityGroups(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DescribeMariadbAccountsByFilter(ctx context.Context, param map[string]interface{}) (accounts []*mariadb.DBAccount, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeAccountsRequest()
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
	response, err := me.client.UseMariadbClient().DescribeAccounts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil || len(response.Response.Users) > 0 {
		accounts = response.Response.Users
	}

	return
}

func (me *MariadbService) DescribeMariadbSecurityGroupsByFilter(ctx context.Context, param map[string]interface{}) (securityGroups []*mariadb.SecurityGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBSecurityGroupsRequest()
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

		if k == "product" {
			request.Product = v.(*string)
		}

	}
	response, err := me.client.UseMariadbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil || len(response.Response.Groups) > 0 {
		securityGroups = response.Response.Groups
	}

	return
}

func (me *MariadbService) DescribeMariadbInstanceById(ctx context.Context, instanceId string) (instance *mariadb.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMariadbClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Instances) < 1 {
		return
	}

	instance = response.Response.Instances[0]
	return
}

func (me *MariadbService) IsolateDBInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewIsolateDBInstanceRequest()
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMariadbClient().IsolateDBInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DeleteMariadbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDestroyDBInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMariadbClient().DestroyDBInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MariadbService) DescribeMariadbDatabaseObjectsByFilter(ctx context.Context, instanceId, dbName string) (databaseObjects *mariadb.DescribeDatabaseObjectsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDatabaseObjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.DbName = &dbName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DescribeDatabaseObjects(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		databaseObjects = response.Response
	}

	return
}

func (me *MariadbService) DescribeMariadbDatabasesByFilter(ctx context.Context, instanceId string) (databases []*mariadb.Database, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDatabasesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DescribeDatabases(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Databases) < 1 {
		return
	}
	databases = response.Response.Databases

	return
}

func (me *MariadbService) DescribeMariadbDatabaseTableByFilter(ctx context.Context, param map[string]interface{}) (cols []*mariadb.TableColumn, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDatabaseTableRequest()
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
		if k == "DbName" {
			request.DbName = v.(*string)
		}
		if k == "Table" {
			request.Table = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DescribeDatabaseTable(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil || len(response.Response.Cols) > 0 {
		cols = response.Response.Cols
	}

	return
}

func (me *MariadbService) DescribeDBEncryptAttributes(ctx context.Context, instanceId string) (encryptAttributes *mariadb.DescribeDBEncryptAttributesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBEncryptAttributesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseMariadbClient().DescribeDBEncryptAttributes(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	encryptAttributes = response.Response

	return
}

func (me *MariadbService) DescribeFlowById(ctx context.Context, flowId int64) (flowParams *mariadb.DescribeFlowResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeFlowRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.FlowId = &flowId
	response, err := me.client.UseMariadbClient().DescribeFlow(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	flowParams = response.Response

	return
}

func (me *MariadbService) DescribeMariadbAccountPrivilegesById(ctx context.Context, instanceId string, user string, host string) (accountPrivileges *mariadb.DescribeAccountPrivilegesResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &instanceId
	request.UserName = &user
	request.Host = &host
	request.DbName = common.StringPtr("*")

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMariadbClient().DescribeAccountPrivileges(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accountPrivileges = response.Response
	return
}

func (me *MariadbService) DescribeMariadbBackupTimeById(ctx context.Context, instanceId string) (backupTime *mariadb.DBBackupTimeConfig, errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDescribeBackupTimeRequest()
	request.InstanceIds = common.StringPtrs([]string{instanceId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMariadbClient().DescribeBackupTime(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	backupTime = response.Response.Items[0]
	return
}

func (me *MariadbService) DescribeDBInstanceDetailById(ctx context.Context, instanceId string) (dbDetail *mariadb.DescribeDBInstanceDetailResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDescribeDBInstanceDetailRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMariadbClient().DescribeDBInstanceDetail(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	dbDetail = response.Response
	return
}
