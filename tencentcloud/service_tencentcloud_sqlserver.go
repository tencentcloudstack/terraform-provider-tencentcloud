package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/pkg/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SqlserverService struct {
	client *connectivity.TencentCloudClient
}

func (me *SqlserverService) DescribeZones(ctx context.Context) (zoneInfoList []*sqlserver.ZoneInfo, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeZonesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *sqlserver.DescribeZonesResponse
	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().DescribeZones(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		zoneInfoList = response.Response.ZoneSet
	}
	return
}

func (me *SqlserverService) DescribeProductConfig(ctx context.Context, zone string) (specInfoList []*sqlserver.SpecInfo, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeProductConfigRequest()
	request.Zone = &zone

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *sqlserver.DescribeProductConfigResponse
	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().DescribeProductConfig(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		specInfoList = response.Response.SpecInfoList
	}
	return
}

func (me *SqlserverService) DescribeSqlserverInstances(ctx context.Context, instanceId string, projectId int, vpcId string, subnetId string, netType int) (instanceList []*sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if instanceId != "" {
		request.InstanceIdSet = []*string{&instanceId}
	}
	if projectId != -1 {
		request.ProjectId = helper.IntUint64(projectId)
	}
	if subnetId != "" && netType != 0 {
		request.SubnetId = &subnetId
	}
	if vpcId != "" && netType != 0 {
		request.VpcId = &vpcId
	}
	if netType == 0 {
		//basic network
		request.VpcId = helper.String("")
		request.SubnetId = helper.String("")
	}
	var offset, limit int64 = 0, 20

	request.Offset = &offset
	request.Limit = &limit

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		instanceList = append(instanceList, response.Response.DBInstances...)
		if len(response.Response.DBInstances) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeSqlserverInstanceById(ctx context.Context, instanceId string) (instance *sqlserver.DBInstance, has bool, errRet error) {
	instanceList, err := me.DescribeSqlserverInstances(ctx, instanceId, -1, "", "", 1)
	if err != nil {
		errRet = err
		return
	}
	if len(instanceList) == 0 {
		return
	} else if len(instanceList) > 1 {
		errRet = fmt.Errorf("[DescribeDBInstances]SDK returns more than one instance with instanceId %s", instanceId)
	}

	instance = instanceList[0]
	if instance != nil && *instance.Status != 8 && *instance.Status != 4 && *instance.Status != 6 {
		has = true
	}
	return
}

func (me *SqlserverService) WaitForTaskFinish(ctx context.Context, flowId int64) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeFlowStatusRequest()
	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	errRet = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		taskResponse, err := me.client.UseSqlserverClient().DescribeFlowStatus(request)
		if err != nil {
			return resource.NonRetryableError(errors.WithStack(err))
		}
		if *taskResponse.Response.Status == int64(SQLSERVER_TASK_RUNNING) {
			return resource.RetryableError(errors.WithStack(fmt.Errorf("SQLSERVER task status is %d(expanding), requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId)))
		} else if *taskResponse.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(errors.WithStack(fmt.Errorf("SQLSERVER task status is %d(failed), requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId)))
		}
		return nil
	})
	return
}

func (me *SqlserverService) CreateSqlserverDB(ctx context.Context, instanceID string, dbname string, charset string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewCreateDBRequest()

	// set instance id
	request.InstanceId = &instanceID
	// set DBs
	var dbCreateInfo = sqlserver.DBCreateInfo{
		DBName:  &dbname,
		Charset: &charset,
		Remark:  &remark,
	}
	var dbInfoList []*sqlserver.DBCreateInfo
	dbInfoList = append(dbInfoList, &dbCreateInfo)
	request.DBs = dbInfoList

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *sqlserver.CreateDBResponse
	err := resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().CreateDB(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response != nil && response.Response != nil {
		return me.WaitForTaskFinish(ctx, *response.Response.FlowId)
	}
	return
}

func (me *SqlserverService) DescribeDBsOfInstance(ctx context.Context, instanceId string) (instanceDBList []*sqlserver.DBDetail, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if instanceId != "" {
		request.InstanceIdSet = []*string{&instanceId}
	}
	var offset, limit uint64 = SQLSERVER_DEFAULT_OFFSET, SQLSERVER_DEFAULT_LIMIT

	request.Offset = &offset
	request.Limit = &limit

	for {
		var response *sqlserver.DescribeDBsResponse
		err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseSqlserverClient().DescribeDBs(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response for api[%s]", request.GetAction())
			return
		}
		if len(response.Response.DBInstances) == 0 {
			return
		} else if len(response.Response.DBInstances) > 1 {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] returned multiple DB lists for one instance", logId, request.GetAction())
			return
		}
		instanceDBList = append(instanceDBList, response.Response.DBInstances[0].DBDetails...)
		if len(response.Response.DBInstances) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeDBDetailsById(ctx context.Context, dbId string) (dbInfo *sqlserver.DBDetail, has bool, errRet error) {
	idItem := strings.Split(dbId, FILED_SP)
	if len(idItem) < 2 {
		errRet = fmt.Errorf("broken ID of SQLServer DB")
		return
	}
	instanceId := idItem[0]
	dbName := idItem[1]

	instanceDBList, err := me.DescribeDBsOfInstance(ctx, instanceId)
	if err != nil {
		errRet = err
		return
	}
	if len(instanceDBList) == 0 {
		return
	}

	for _, dbDetail := range instanceDBList {
		if *dbDetail.Name == dbName {
			dbInfo = dbDetail
			if *dbDetail.Status != SQLSERVER_DB_DELETING {
				has = true
			}
			break
		}
	}
	return
}

func (me *SqlserverService) ModifySqlserverDBRemark(ctx context.Context, instanceId string, dbName string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyDBRemarkRequest()
	request.InstanceId = &instanceId
	request.DBRemarks = []*sqlserver.DBRemark{{Name: &dbName, Remark: &remark}}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	err := resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, e := me.client.UseSqlserverClient().ModifyDBRemark(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (me *SqlserverService) DeleteSqlserverDB(ctx context.Context, instanceId string, name string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteDBRequest()
	request.InstanceId = &instanceId
	request.Names = []*string{&name}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var response *sqlserver.DeleteDBResponse
	err := resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().DeleteDB(request)
		if e != nil {
			log.Printf("[CRITAL]%s %s fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response != nil && response.Response != nil {
		return me.WaitForTaskFinish(ctx, *response.Response.FlowId)
	}
	return
}
