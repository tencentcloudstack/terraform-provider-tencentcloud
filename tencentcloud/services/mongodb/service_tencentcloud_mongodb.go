package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMongodbService(client *connectivity.TencentCloudClient) MongodbService {
	return MongodbService{client: client}
}

type MongodbService struct {
	client *connectivity.TencentCloudClient
}

func (me *MongodbService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *mongodb.InstanceDetail, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = instanceId
	var response *mongodb.DescribeDBInstancesResponse
	err := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient(iacExtInfo).DescribeDBInstances(request)
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
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewResetDBInstancePasswordRequest()
	request.InstanceId = &instanceId
	request.UserName = &accountName
	request.Password = &password
	var response *mongodb.ResetDBInstancePasswordResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
		if err = me.DescribeAsyncRequestInfo(ctx, *response.Response.AsyncRequestId, 3*tccommon.ReadRetryTimeout); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) ModifyMongosMemory(ctx context.Context, instanceId string, mongosMemory int) (dealId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewModifyDBInstanceSpecRequest()
	request.InstanceId = &instanceId
	request.MongosMemory = helper.String(helper.IntToStr(mongosMemory))

	var response *mongodb.ModifyDBInstanceSpecResponse
	tradeError := false
	err := resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && response.Response.DealId != nil {
		dealId = *response.Response.DealId
	}
	return
}
func (me *MongodbService) UpgradeInstance(ctx context.Context, instanceId string, memory int, volume int, params map[string]interface{}) (dealId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewModifyDBInstanceSpecRequest()
	request.InstanceId = &instanceId
	request.Memory = helper.IntUint64(memory)
	request.Volume = helper.IntUint64(volume)
	if v, ok := params["node_num"]; ok {
		request.NodeNum = helper.IntUint64(v.(int))
	}
	if v, ok := params["add_node_list"]; ok {
		addNodeList := v.([]interface{})
		for _, addNode := range addNodeList {
			addNodeMap := addNode.(map[string]interface{})
			request.AddNodeList = append(request.AddNodeList, &mongodb.AddNodeList{
				Role: helper.String(addNodeMap["role"].(string)),
				Zone: helper.String(addNodeMap["zone"].(string)),
			})
		}
	}
	if v, ok := params["remove_node_list"]; ok {
		removeNodeList := v.([]interface{})
		for _, removeNode := range removeNodeList {
			removeNodeMap := removeNode.(map[string]interface{})
			request.RemoveNodeList = append(request.RemoveNodeList, &mongodb.RemoveNodeList{
				Role:     helper.String(removeNodeMap["role"].(string)),
				Zone:     helper.String(removeNodeMap["zone"].(string)),
				NodeName: helper.String(removeNodeMap["node_name"].(string)),
			})
		}
	}
	if v, ok := params["in_maintenance"]; ok {
		request.InMaintenance = helper.IntUint64(v.(int))
	}
	var response *mongodb.ModifyDBInstanceSpecResponse
	tradeError := false
	err := resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && response.Response.DealId != nil {
		dealId = *response.Response.DealId
	}

	return
}

func (me *MongodbService) ModifyProjectId(ctx context.Context, instanceId string, projectId int) (errRet error) {
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)
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

func (me *MongodbService) ModifySecurityGroups(ctx context.Context, instanceId string, securityGroups []*string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewModifyDBInstanceSecurityGroupRequest()
	request.InstanceId = &instanceId
	request.SecurityGroupIds = securityGroups
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMongodbClient().ModifyDBInstanceSecurityGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) ModifyNetworkAddress(ctx context.Context, instanceId string, vpcId string, subnetId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewModifyDBInstanceNetworkAddressRequest()
	request.InstanceId = &instanceId
	request.NewUniqVpcId = &vpcId
	request.NewUniqSubnetId = &subnetId
	request.OldIpExpiredTime = helper.Uint64(0)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMongodbClient().ModifyDBInstanceNetworkAddress(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) DescribeInstancesByFilter(ctx context.Context, instanceId string,
	clusterType int) (mongodbs []*mongodb.InstanceDetail, errRet error) {

	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewIsolateDBInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *mongodb.IsolateDBInstanceResponse
	err := resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().IsolateDBInstance(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return tccommon.RetryError(e)
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
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewTerminateDBInstancesRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *mongodb.TerminateDBInstancesResponse
	err := resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
	logId := tccommon.GetLogId(ctx)
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
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, e := me.client.UseMongodbClient().RenewDBInstances(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (me *MongodbService) DescribeAsyncRequestInfo(ctx context.Context, asyncId string, timeout time.Duration) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewDescribeAsyncRequestInfoRequest()
	request.AsyncRequestId = &asyncId
	err := resource.Retry(timeout, func() *resource.RetryError {
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
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewOfflineIsolatedDBInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var response *mongodb.OfflineIsolatedDBInstanceResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
	err = resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewDescribeSecurityGroupRequest()
	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseMongodbClient().DescribeSecurityGroup(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		groups = response.Response.Groups
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *MongodbService) DescribeDBInstanceNodeProperty(ctx context.Context, instanceId string) (replicateSets []*mongodb.ReplicateSetInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := mongodb.NewDescribeDBInstanceNodePropertyRequest()
	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseMongodbClient().DescribeDBInstanceNodeProperty(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		replicateSets = response.Response.ReplicateSets
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *MongodbService) DescribeMongodbInstanceAccountById(ctx context.Context, instanceId string, userName string) (instanceAccount *mongodb.UserInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewDescribeAccountUsersRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DescribeAccountUsers(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Users) < 1 {
		return
	}

	for _, user := range response.Response.Users {
		if *user.UserName == userName {
			instanceAccount = user
			return
		}
	}
	return
}

func (me *MongodbService) DeleteMongodbInstanceAccountById(ctx context.Context, instanceId string, userName string, mongoUserPassword string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewDeleteAccountUserRequest()
	request.InstanceId = &instanceId
	request.UserName = &userName
	request.MongoUserPassword = &mongoUserPassword

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DeleteAccountUser(request)
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		if err = me.DescribeAsyncRequestInfo(ctx, helper.Int64ToStr(*response.Response.FlowId), 3*tccommon.ReadRetryTimeout); err != nil {
			errRet = err
			return
		}
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MongodbService) DescribeMongodbInstanceBackupDownloadTaskById(ctx context.Context, instanceId string, backupName string) (instanceBackupDownloadTask []*mongodb.BackupDownloadTask, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewDescribeBackupDownloadTaskRequest()
	request.InstanceId = &instanceId
	request.BackupName = &backupName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DescribeBackupDownloadTask(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Tasks) < 1 {
		return
	}

	instanceBackupDownloadTask = response.Response.Tasks
	return
}

func (me *MongodbService) DescribeMongodbInstanceBackupsByFilter(ctx context.Context, param map[string]interface{}) (instanceBackups []*mongodb.BackupInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mongodb.NewDescribeDBBackupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "backup_method" {
			request.BackupMethod = v.(*int64)
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
		response, err := me.client.UseMongodbClient().DescribeDBBackups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.BackupList) < 1 {
			break
		}
		instanceBackups = append(instanceBackups, response.Response.BackupList...)
		if len(response.Response.BackupList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MongodbService) DescribeMongodbInstanceConnectionsByFilter(ctx context.Context, param map[string]interface{}) (instanceConnections []*mongodb.ClientConnection, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mongodb.NewDescribeClientConnectionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
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
		response, err := me.client.UseMongodbClient().DescribeClientConnections(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Clients) < 1 {
			break
		}
		instanceConnections = append(instanceConnections, response.Response.Clients...)
		if len(response.Response.Clients) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MongodbService) DescribeMongodbInstanceCurrentOpByFilter(ctx context.Context, param map[string]interface{}) (instanceCurrentOp []*mongodb.CurrentOp, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mongodb.NewDescribeCurrentOpRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "ns" {
			request.Ns = v.(*string)
		}
		if k == "millisecond_running" {
			request.MillisecondRunning = v.(*uint64)
		}
		if k == "op" {
			request.Op = v.(*string)
		}
		if k == "replica_set_name" {
			request.ReplicaSetName = v.(*string)
		}
		if k == "state" {
			request.State = v.(*string)
		}
		if k == "order_by" {
			request.OrderBy = v.(*string)
		}
		if k == "order_by_type" {
			request.OrderByType = v.(*string)
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
		response, err := me.client.UseMongodbClient().DescribeCurrentOp(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CurrentOps) < 1 {
			break
		}
		instanceCurrentOp = append(instanceCurrentOp, response.Response.CurrentOps...)
		if len(response.Response.CurrentOps) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MongodbService) DescribeMongodbInstanceParams(ctx context.Context, param map[string]interface{}) (instanceParams *mongodb.DescribeInstanceParamsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mongodb.NewDescribeInstanceParamsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceParams = response.Response

	return
}

func (me *MongodbService) DescribeMongodbInstanceSlowLogByFilter(ctx context.Context, param map[string]interface{}) (instanceSlowLog []*string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mongodb.NewDescribeSlowLogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "slow_ms" {
			request.SlowMS = v.(*uint64)
		}
		if k == "format" {
			request.Format = v.(*string)
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
		response, err := me.client.UseMongodbClient().DescribeSlowLogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SlowLogs) < 1 {
			break
		}
		instanceSlowLog = append(instanceSlowLog, response.Response.SlowLogs...)
		if len(response.Response.SlowLogs) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MongodbService) DescribeDBInstanceDeal(ctx context.Context, dealId string) (dealResponseParams *mongodb.DescribeDBInstanceDealResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewDescribeDBInstanceDealRequest()
	request.DealId = helper.String(dealId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DescribeDBInstanceDeal(request)
	if err != nil {
		errRet = err
		return
	}
	dealResponseParams = response.Response
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MongodbService) SetInstanceMaintenance(ctx context.Context, instanceId, maintenanceStart, maintenanceEnd string) error {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewSetInstanceMaintenanceRequest()
	request.InstanceId = helper.String(instanceId)
	request.MaintenanceStart = helper.String(maintenanceStart)
	request.MaintenanceEnd = helper.String(maintenanceEnd)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().SetInstanceMaintenance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *MongodbService) DescribeMongodbInstanceParamValues(ctx context.Context, instanceId string, paramNames []string) (res map[string]string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewDescribeInstanceParamsRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	res = make(map[string]string)
	for _, param := range response.Response.InstanceEnumParam {
		for _, paramName := range paramNames {
			if *param.ParamName == paramName {
				res[paramName] = *param.CurrentValue
			}
		}
	}
	for _, param := range response.Response.InstanceIntegerParam {
		for _, paramName := range paramNames {
			if *param.ParamName == paramName {
				res[paramName] = *param.CurrentValue
			}
		}
	}
	for _, param := range response.Response.InstanceMultiParam {
		for _, paramName := range paramNames {
			if *param.ParamName == paramName {
				res[paramName] = *param.CurrentValue
			}
		}
	}
	for _, param := range response.Response.InstanceTextParam {
		for _, paramName := range paramNames {
			if *param.ParamName == paramName {
				res[paramName] = *param.CurrentValue
			}
		}
	}

	return
}

func (me *MongodbService) DescribeMongodbInstanceUrls(ctx context.Context, instanceId string) (ret []*mongodb.DbURL, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mongodb.NewDescribeDBInstanceURLRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = helper.String(instanceId)

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMongodbClient().DescribeDBInstanceURL(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response.Urls
	return
}

func (me *MongodbService) DescribeMongodbInstanceSSLById(ctx context.Context, instanceId string) (sslStatus *mongodb.DescribeInstanceSSLResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewDescribeInstanceSSLRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *mongodb.DescribeInstanceSSLResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().DescribeInstanceSSL(request)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read mongodb instance ssl failed, reason: %v", logId, err)
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	sslStatus = response
	return
}

func (me *MongodbService) ModifyMongodbInstanceSSL(ctx context.Context, instanceId string, enable bool) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mongodb.NewInstanceEnableSSLRequest()
	request.InstanceId = &instanceId
	request.Enable = &enable

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMongodbClient().InstanceEnableSSL(request)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify mongodb instance ssl failed, reason: %v", logId, err)
		errRet = err
		return
	}

	return
}
