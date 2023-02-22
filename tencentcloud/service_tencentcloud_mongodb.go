package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MongodbService struct {
	client *connectivity.TencentCloudClient
}

func (me *MongodbService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *mongodb.InstanceDetail, has bool, errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	var response *mongodb.DescribeDBInstancesResponse
	err := resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().DescribeDBInstances(request)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if result != nil && result.Response != nil {
			if len(result.Response.InstanceDetails) != 0 && (*result.Response.InstanceDetails[0].Status == MONGODB_INSTANCE_STATUS_PROCESSING ||
				*result.Response.InstanceDetails[0].Status == MONGODB_INSTANCE_STATUS_INITIAL) {
				return resource.RetryableError(fmt.Errorf("mongodb instance status is processing"))
			}
			response = result
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("response is null"))
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceDetails) == 0 || *response.Response.InstanceDetails[0].Status == MONGODB_INSTANCE_STATUS_EXPIRED {
		has = false
		return
	}

	has = true
	instance = response.Response.InstanceDetails[0]
	return
}

func (me *MongodbService) ModifyInstanceName(ctx context.Context, instanceId, instanceName string) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewRenameInstanceRequest()
	request.InstanceId = &instanceId
	request.NewName = &instanceName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMongodbClient().RenameInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) ResetInstancePassword(ctx context.Context, instanceId, accountName, password string) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewResetDBInstancePasswordRequest()
	request.InstanceId = &instanceId
	request.UserName = &accountName
	request.Password = &password
	var response *mongodb.ResetDBInstancePasswordResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().ResetDBInstancePassword(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.RetryableError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response != nil && response.Response != nil {
		if err = me.DescribeAsyncRequestInfo(ctx, *response.Response.AsyncRequestId); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) UpgradeInstance(ctx context.Context, instanceId string, memory int, volume int) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewModifyDBInstanceSpecRequest()
	request.InstanceId = &instanceId
	request.Memory = helper.IntUint64(memory)
	request.Volume = helper.IntUint64(volume)
	var response *mongodb.ModifyDBInstanceSpecResponse
	tradeError := false
	err := resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().ModifyDBInstanceSpec(request)
		if e != nil {
			// request might be accepted between "InvalidParameterValue.InvalidTradeOperation" and "InvalidParameterValue.StatusAbnormal" error
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "InvalidParameterValue.InvalidTradeOperation" {
					tradeError = true
					return resource.RetryableError(e)
				} else if ee.Code == "InvalidParameterValue.StatusAbnormal" && tradeError {
					response = result
					return nil
				} else {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) ModifyProjectId(ctx context.Context, instanceId string, projectId int) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewAssignProjectRequest()
	request.InstanceIds = []*string{&instanceId}
	request.ProjectId = helper.IntUint64(projectId)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMongodbClient().AssignProject(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) DescribeSpecInfo(ctx context.Context, zone string) (infos []*mongodb.SpecificationInfo, errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewDescribeSpecInfoRequest()
	if zone != "" {
		request.Zone = &zone
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMongodbClient().DescribeSpecInfo(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	infos = response.Response.SpecInfoList
	return
}

func (me *MongodbService) DescribeInstancesByFilter(ctx context.Context, instanceId string,
	clusterType int) (mongodbs []*mongodb.InstanceDetail, errRet error) {

	logId := getLogId(ctx)
	request := mongodb.NewDescribeDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	if instanceId != "" {
		request.InstanceIds = []*string{&instanceId}
	}
	if clusterType > 0 {
		temp := int64(clusterType)
		request.ClusterType = &temp
	}

	offset := MONGODB_DEFAULT_OFFSET
	pageSize := MONGODB_MAX_LIMIT
	mongodbs = make([]*mongodb.InstanceDetail, 0)
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseMongodbClient().DescribeDBInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceDetails) < 1 {
			break
		}

		mongodbs = append(mongodbs, response.Response.InstanceDetails...)

		if len(response.Response.InstanceDetails) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}

func (me *MongodbService) IsolateInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewIsolateDBInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *mongodb.IsolateDBInstanceResponse
	err := resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().IsolateDBInstance(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MongodbService) TerminateDBInstances(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewTerminateDBInstancesRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *mongodb.TerminateDBInstancesResponse
	err := resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().TerminateDBInstances(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.RetryableError(fmt.Errorf("Terminate instance %s error", instanceId))
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MongodbService) ModifyAutoRenewFlag(ctx context.Context, instanceId string, period int, renewFlag int) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewRenewDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	request.InstanceChargePrepaid = &mongodb.InstanceChargePrepaid{}
	request.InstanceChargePrepaid.Period = helper.IntInt64(period)
	request.InstanceChargePrepaid.RenewFlag = helper.String(MONGODB_AUTO_RENEW_FLAG[renewFlag])
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, e := me.client.UseMongodbClient().RenewDBInstances(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (me *MongodbService) DescribeAsyncRequestInfo(ctx context.Context, asyncId string) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewDescribeAsyncRequestInfoRequest()
	request.AsyncRequestId = &asyncId
	err := resource.Retry(readRetryTimeout*3, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().DescribeAsyncRequestInfo(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.RetryableError(e)
		}
		if *result.Response.Status == MONGODB_TASK_FAILED {
			return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s api[%s] task failed", logId, request.GetAction()))
		}
		if *result.Response.Status != MONGODB_TASK_SUCCESS {
			return resource.RetryableError(fmt.Errorf("[CRITAL]%s api[%s] task is %s, retrying", logId, request.GetAction(), *result.Response.Status))
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *MongodbService) OfflineIsolatedDBInstance(ctx context.Context, instanceId string, timeOutTolerant bool) (errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewOfflineIsolatedDBInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *mongodb.OfflineIsolatedDBInstanceResponse
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseMongodbClient().OfflineIsolatedDBInstance(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil {
		return nil
	}

	checkRequest := mongodb.NewDescribeDBInstancesRequest()
	checkRequest.InstanceIds = []*string{&instanceId}
	err = resource.Retry(20*readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(checkRequest.GetAction())
		result, e := me.client.UseMongodbClient().DescribeDBInstances(checkRequest)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if result != nil && result.Response != nil {
			if len(result.Response.InstanceDetails) != 0 {
				return resource.RetryableError(fmt.Errorf("Offline mongodb instance is processing"))
			}
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("response is null"))
	})
	if err != nil {
		return err
	}
	return nil
}

func (me *MongodbService) DescribeSecurityGroup(ctx context.Context, instanceId string) (groups []*mongodb.SecurityGroup, errRet error) {
	logId := getLogId(ctx)
	request := mongodb.NewDescribeSecurityGroupRequest()
	request.InstanceId = helper.String(instanceId)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMongodbClient().DescribeSecurityGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	groups = response.Response.Groups
	return
}
