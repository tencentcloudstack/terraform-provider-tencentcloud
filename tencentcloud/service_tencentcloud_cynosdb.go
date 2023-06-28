package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CynosdbService struct {
	client *connectivity.TencentCloudClient
}

func (me *CynosdbService) DescribeRedisZoneConfig(ctx context.Context) (instanceSpecSet []*cynosdb.InstanceSpec, err error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeInstanceSpecsRequest()

	request.DbType = helper.String("MYSQL")
	request.IncludeZoneStocks = helper.Bool(true)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseCynosdbClient().DescribeInstanceSpecs(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		}
		instanceSpecSet = response.Response.InstanceSpecSet
		return nil
	})

	return
}

func (me *CynosdbService) DescribeClusters(ctx context.Context, filters map[string]string) (clusters []*cynosdb.CynosdbCluster, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeClustersRequest()

	queryFilters := make([]*cynosdb.QueryFilter, 0, len(filters))

	for k, v := range filters {
		filter := cynosdb.QueryFilter{
			Names:  []*string{helper.String(k)},
			Values: []*string{helper.String(v)},
		}
		queryFilters = append(queryFilters, &filter)
	}
	request.Filters = queryFilters

	offset := CYNOSDB_DEFAULT_OFFSET
	pageSize := CYNOSDB_MAX_LIMIT
	clusters = make([]*cynosdb.CynosdbCluster, 0)
	for {
		request.Offset = helper.IntInt64(offset)
		request.Limit = helper.IntInt64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCynosdbClient().DescribeClusters(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}

		clusters = append(clusters, response.Response.ClusterSet...)

		if len(response.Response.ClusterSet) < pageSize {
			break
		}
		offset += pageSize
	}

	return
}

/**
Return values:
	clusterItem: ResponseBody of DescribeClusters, include `renew_flag` and `db_mode`
    clusterInfo: ResponseBody of DescribeClusterDetailResponse, primary args setter.
*/
func (me *CynosdbService) DescribeClusterById(ctx context.Context, clusterId string) (clusterItem *cynosdb.CynosdbCluster, clusterInfo *cynosdb.CynosdbClusterDetail, has bool, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeClusterDetailRequest()
	request.ClusterId = &clusterId

	// get cluster status
	var notExist bool
	var clusters []*cynosdb.CynosdbCluster
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		clusters, errRet = me.DescribeClusters(ctx, map[string]string{"ClusterId": clusterId})
		if errRet != nil {
			return retryError(errRet)
		}
		if len(clusters) == 0 {
			notExist = true
			return nil
		}
		if len(clusters) != 1 {
			return resource.NonRetryableError(fmt.Errorf("[CRITAL] mutiple cluster found by cluster id %s", clusterId))
		}
		clusterItem = clusters[0]
		clusterStatus := clusterItem.Status
		if clusterStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("cluster %s status is nil", clusterId))
		}
		if *clusterStatus == CYNOSDB_STATUS_RUNNING {
			return nil
		}
		if *clusterStatus == CYNOSDB_STATUS_ISOLATED || *clusterStatus == CYNOSDB_STATUS_OFFLINE || *clusterStatus == CYNOSDB_STATUS_DELETED {
			notExist = true
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cynosdb cluster %s is still in processing", clusterId))
	})
	if errRet != nil || notExist {
		return
	}
	has = true

	var response *cynosdb.DescribeClusterDetailResponse
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseCynosdbClient().DescribeClusterDetail(request)
		if errRet != nil {
			return retryError(errRet)
		}

		return nil
	})
	if errRet != nil {
		log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
		return
	}

	clusterInfo = response.Response.Detail
	return
}

func (me *CynosdbService) UpgradeInstance(ctx context.Context, instanceId string, cpu, mem int64) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewUpgradeInstanceRequest()

	request.InstanceId = &instanceId
	request.Cpu = &cpu
	request.Memory = &mem
	request.UpgradeType = helper.String(CYNOSDB_UPGRADE_IMMEDIATE)

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().UpgradeInstance(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) DescribeMaintainPeriod(ctx context.Context, instanceId string) (response *cynosdb.DescribeMaintainPeriodResponse, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeMaintainPeriodRequest()

	request.InstanceId = &instanceId

	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseCynosdbClient().DescribeMaintainPeriod(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) ModifyMaintainPeriodConfig(ctx context.Context, instanceId string, startTime, duration int64, weekdays []*string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewModifyMaintainPeriodConfigRequest()

	request.InstanceId = &instanceId
	request.MaintainStartTime = &startTime
	request.MaintainDuration = &duration
	request.MaintainWeekDays = weekdays

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().ModifyMaintainPeriodConfig(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) IsolateCluster(ctx context.Context, clusterId string) (flowId int64, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewIsolateClusterRequest()

	request.ClusterId = &clusterId

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseCynosdbClient().IsolateCluster(request)
		if errRet != nil {
			// isolate after creation immediately will encounter this error
			if ee, ok := errRet.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.Contains(ee.Message, "return not found valid deal") {
					return resource.RetryableError(ee)
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		if response.Response.FlowId != nil {
			flowId = *response.Response.FlowId
		}
		return nil
	})

	return
}

func (me *CynosdbService) OfflineCluster(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewOfflineClusterRequest()

	request.ClusterId = &clusterId

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().OfflineCluster(request)
		if errRet != nil {
			if ee, ok := errRet.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.Contains(ee.Message, "IsolateInstanceFlow failed") {
					return resource.RetryableError(ee)
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})

	return
}

func (me *CynosdbService) DescribeInstances(ctx context.Context, filters map[string]string) (instances []*cynosdb.CynosdbInstance, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeInstancesRequest()

	queryFilters := make([]*cynosdb.QueryFilter, 0, len(filters))

	for k, v := range filters {
		filter := cynosdb.QueryFilter{
			Names:  []*string{helper.String(k)},
			Values: []*string{helper.String(v)},
		}
		queryFilters = append(queryFilters, &filter)
	}
	request.Filters = queryFilters

	offset := CYNOSDB_DEFAULT_OFFSET
	pageSize := CYNOSDB_MAX_LIMIT
	instances = make([]*cynosdb.CynosdbInstance, 0)
	for {
		request.Offset = helper.IntInt64(offset)
		request.Limit = helper.IntInt64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCynosdbClient().DescribeInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}

		instances = append(instances, response.Response.InstanceSet...)

		if len(response.Response.InstanceSet) < pageSize {
			break
		}
		offset += pageSize
	}

	return
}

func (me *CynosdbService) DescribeInstanceById(ctx context.Context, instanceId string) (clusterId string, instanceInfo *cynosdb.CynosdbInstanceDetail, has bool, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeInstanceDetailRequest()
	request.InstanceId = &instanceId

	var notExist bool
	var instances []*cynosdb.CynosdbInstance
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		instances, errRet = me.DescribeInstances(ctx, map[string]string{"InstanceId": instanceId})
		if errRet != nil {
			return retryError(errRet)
		}
		if len(instances) == 0 {
			notExist = true
			return nil
		}
		if len(instances) != 1 {
			return resource.NonRetryableError(fmt.Errorf("[CRITAL] mutiple instance found by cluster id %s", instanceId))
		}
		if *instances[0].Status == CYNOSDB_STATUS_ISOLATED || *instances[0].Status == CYNOSDB_STATUS_OFFLINE || *instances[0].Status == CYNOSDB_STATUS_DELETED {
			notExist = true
			return nil
		} else if *instances[0].Status == CYNOSDB_STATUS_RUNNING {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cynosdb instance %s is still in processing", instanceId))
		}
	})
	if errRet != nil || notExist {
		return
	}
	has = true
	clusterId = *instances[0].ClusterId

	var response *cynosdb.DescribeInstanceDetailResponse
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseCynosdbClient().DescribeInstanceDetail(request)
		if errRet != nil {
			return retryError(errRet)
		}

		return nil
	})
	if errRet != nil {
		log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
		return
	}

	instanceInfo = response.Response.Detail
	return
}

func (me *CynosdbService) DescribeClusterInstanceGrps(ctx context.Context, clusterId string) (response *cynosdb.DescribeClusterInstanceGrpsResponse, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeClusterInstanceGrpsRequest()
	request.ClusterId = &clusterId

	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseCynosdbClient().DescribeClusterInstanceGrps(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) DescribeInsGrpSecurityGroups(ctx context.Context, instanceGrpId string) (response *cynosdb.DescribeDBSecurityGroupsResponse, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &instanceGrpId

	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseCynosdbClient().DescribeDBSecurityGroups(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) ModifyInsGrpSecurityGroups(ctx context.Context, insGrp, az string, sg []*string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewModifyDBInstanceSecurityGroupsRequest()

	request.InstanceId = &insGrp
	request.Zone = &az
	request.SecurityGroupIds = sg

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().ModifyDBInstanceSecurityGroups(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) ModifyClusterParam(ctx context.Context, request *cynosdb.ModifyClusterParamRequest) (asyncReqId string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCynosdbClient().ModifyClusterParam(request)

	if err != nil {
		errRet = err
		return
	}

	asyncReqId = *response.Response.AsyncRequestId

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) IsolateInstance(ctx context.Context, clusterId, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewIsolateInstanceRequest()

	request.ClusterId = &clusterId
	request.InstanceIdList = []*string{helper.String(instanceId)}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().IsolateInstance(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})

	return
}

func (me *CynosdbService) OfflineInstance(ctx context.Context, clusterId, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewOfflineInstanceRequest()

	request.ClusterId = &clusterId
	request.InstanceIdList = []*string{helper.String(instanceId)}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().OfflineInstance(request)
		if errRet != nil {
			if ee, ok := errRet.(*sdkErrors.TencentCloudSDKError); ok {
				if strings.Contains(ee.Message, "OfflineInstanceFlow failed") {
					return resource.RetryableError(ee)
				}
			}

			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})

	return
}

func (me *CynosdbService) DescribeClusterParams(ctx context.Context, clusterId string) (items []*cynosdb.ParamInfo, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeClusterParamsRequest()
	request.ClusterId = &clusterId

	var response *cynosdb.DescribeClusterParamsResponse
	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseCynosdbClient().DescribeClusterParams(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}
	items = response.Response.Items

	return
}

func (me *CynosdbService) ResumeServerless(ctx context.Context, request *cynosdb.ResumeServerlessRequest) (errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCynosdbClient().ResumeServerless(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) PauseServerless(ctx context.Context, request *cynosdb.PauseServerlessRequest) (errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCynosdbClient().PauseServerless(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) SwitchServerlessCluster(ctx context.Context, clusterId string, resume bool) error {
	pause := !resume
	_, detail, _, err := me.DescribeClusterById(ctx, clusterId)
	if err != nil {
		return err
	}
	st := detail.ServerlessStatus
	if st == nil {
		return nil
	}
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if resume && *st == "resuming" || pause && *st == "pausing" {
			return resource.RetryableError(fmt.Errorf("waiting for status %s finish", *st))
		}
		if resume && *st == "resume" || pause && *st == "pause" {
			return nil
		}
		if resume && *st == "pause" {
			request := cynosdb.NewResumeServerlessRequest()
			request.ClusterId = &clusterId
			err := me.ResumeServerless(ctx, request)
			if err != nil {
				return retryError(err, cynosdb.OPERATIONDENIED_SERVERLESSCLUSTERSTATUSDENIED)
			}
			return nil
		}
		if pause && *st == "resume" {
			request := cynosdb.NewPauseServerlessRequest()
			request.ClusterId = &clusterId
			err := me.PauseServerless(ctx, request)
			if err != nil {
				return retryError(err, cynosdb.OPERATIONDENIED_SERVERLESSCLUSTERSTATUSDENIED)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	statusChangeRetry := 5
	return resource.Retry(readRetryTimeout*5, func() *resource.RetryError {
		_, detail, _, err = me.DescribeClusterById(ctx, clusterId)
		if err != nil {
			return retryError(err)
		}
		st := detail.ServerlessStatus
		if st == nil {
			return resource.NonRetryableError(fmt.Errorf("cannot read serverless cluster status"))
		}
		if resume && *st == "pause" || pause && *st == "resume" {
			if statusChangeRetry > 0 {
				statusChangeRetry -= 1
				return resource.RetryableError(fmt.Errorf("waiting for status change, retry %d", statusChangeRetry))
			}
			return resource.NonRetryableError(fmt.Errorf("api action invoked but status still %s", *st))
		}
		if resume && *st == "resuming" || pause && *st == "pausing" {
			statusChangeRetry = 0
			return resource.RetryableError(fmt.Errorf("waiting for status %s finished", *st))
		}
		return nil
	})
}

func (me *CynosdbService) DescribeCynosdbAuditLogFileById(ctx context.Context, instanceId string, fileName string) (auditLogFile *cynosdb.AuditLogFile, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeAuditLogFilesRequest()
	request.InstanceId = &instanceId
	request.FileName = &fileName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeAuditLogFiles(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) < 1 {
		return
	}

	auditLogFile = response.Response.Items[0]
	return
}

func (me *CynosdbService) DeleteCynosdbAuditLogFileById(ctx context.Context, instanceId string, fileName string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDeleteAuditLogFileRequest()
	request.InstanceId = &instanceId
	request.FileName = &fileName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DeleteAuditLogFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) DeleteCynosdbSecurityGroupById(ctx context.Context, instanceId string, securityGroupIds []*string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDisassociateSecurityGroupsRequest()
	request.InstanceIds = []*string{&instanceId}
	request.SecurityGroupIds = securityGroupIds

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DisassociateSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) DescribeCynosdbSecurityGroups(ctx context.Context, instanceId string) (securityGroups []*cynosdb.SecurityGroup, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	securityGroups = response.Response.Groups
	return
}

func (me *CynosdbService) DescribeCynosdbBackup(ctx context.Context, clusterId string, params map[string]interface{}) (backups []*cynosdb.BackupFileInfo, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeBackupListRequest()
	request.ClusterId = helper.String(clusterId)
	if v, ok := params["backup_type"]; ok {
		request.BackupType = helper.String(v.(string))
	}
	if v, ok := params["backup_name"]; ok {
		request.BackupNames = helper.Strings([]string{v.(string)})
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeBackupList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	backups = response.Response.BackupList
	return
}

func (me *CynosdbService) DescribeCynosdbBackupById(ctx context.Context, clusterId, backupId string) (backup *cynosdb.BackupFileInfo, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeBackupListRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId
	backupIdInt64, err := strconv.ParseInt(backupId, 10, 64)
	if err != nil {
		errRet = err
		return
	}
	request.BackupIds = []*int64{&backupIdInt64}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeBackupList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.BackupList) < 1 {
		return
	}

	backup = response.Response.BackupList[0]
	return
}

func (me *CynosdbService) DeleteCynosdbBackupById(ctx context.Context, clusterId, backupId string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDeleteBackupRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId
	backupIdInt64, err := strconv.ParseInt(backupId, 10, 64)
	if err != nil {
		errRet = err
		return
	}
	request.BackupIds = []*int64{&backupIdInt64}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DeleteBackup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) ModifyBackupConfig(ctx context.Context, clusterId string, params map[string]interface{}) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewModifyBackupConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	if v, ok := params["backup_time_beg"]; ok {
		request.BackupTimeBeg = helper.IntUint64(v.(int))
	}
	if v, ok := params["backup_time_end"]; ok {
		request.BackupTimeEnd = helper.IntUint64(v.(int))
	}
	if v, ok := params["reserve_duration"]; ok {
		request.ReserveDuration = helper.IntUint64(v.(int))
	}
	if v, ok := params["backup_freq"]; ok {
		backupFreqs := make([]*string, 0)
		for _, item := range v.([]interface{}) {
			backupFreq := item.(string)
			backupFreqs = append(backupFreqs, helper.String(backupFreq))
		}
		request.BackupFreq = backupFreqs
	}
	if v, ok := params["backup_type"]; ok {
		request.BackupType = helper.String(v.(string))
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseCynosdbClient().ModifyBackupConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb backupConfig failed, reason:%+v", logId, err)
		return err
	}
	return nil
}

func (me *CynosdbService) DescribeCynosdbAccountsByFilter(ctx context.Context, clusterId string, paramMap map[string]interface{}) (result *cynosdb.DescribeAccountsResponseParams, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeAccountsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = helper.String(clusterId)
	if v, ok := paramMap["account_names"]; ok {
		request.AccountNames = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := paramMap["hosts"]; ok {
		request.Hosts = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := me.client.UseCynosdbClient().DescribeAccounts(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}
		result = response.Response
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *CynosdbService) DescribeClusterParamsByFilter(ctx context.Context, clusterId string, paramMap map[string]interface{}) (items []*cynosdb.ParamInfo, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeClusterParamsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = helper.String(clusterId)
	if v, ok := paramMap["param_name"]; ok {
		request.ParamName = helper.String(v.(string))
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := me.client.UseCynosdbClient().DescribeClusterParams(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}
		items = response.Response.Items
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *CynosdbService) DescribeCynosdbParamTemplatesByFilter(ctx context.Context, paramMap map[string]interface{}) (items []*cynosdb.ParamTemplateListInfo, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeParamTemplatesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if v, ok := paramMap["engine_versions"]; ok {
		request.EngineVersions = v.([]*string)
	}
	if v, ok := paramMap["template_names"]; ok {
		request.TemplateNames = v.([]*string)
	}
	if v, ok := paramMap["template_ids"]; ok {
		request.TemplateIds = v.([]*int64)
	}
	if v, ok := paramMap["db_modes"]; ok {
		request.DbModes = v.([]*string)
	}
	if v, ok := paramMap["products"]; ok {
		request.Products = v.([]*string)
	}
	if v, ok := paramMap["template_types"]; ok {
		request.TemplateTypes = v.([]*string)
	}
	if v, ok := paramMap["engine_types"]; ok {
		request.EngineTypes = v.([]*string)
	}
	if v, ok := paramMap["offset"]; ok {
		request.Offset = v.(*int64)
	}
	if v, ok := paramMap["limit"]; ok {
		request.Limit = v.(*int64)
	}
	if v, ok := paramMap["engine_types"]; ok {
		request.EngineTypes = v.([]*string)
	}
	if v, ok := paramMap["order_by"]; ok {
		request.OrderBy = v.(*string)
	}
	if v, ok := paramMap["order_direction"]; ok {
		request.OrderDirection = v.(*string)
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		response, e := me.client.UseCynosdbClient().DescribeParamTemplates(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}
		items = response.Response.Items
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *CynosdbService) DescribeCynosdbAuditLogsByFilter(ctx context.Context, param map[string]interface{}) (auditLogs []*cynosdb.AuditLog, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeAuditLogsRequest()
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
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.(*cynosdb.AuditLogFilter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCynosdbClient().DescribeAuditLogs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}

		auditLogs = append(auditLogs, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbBackupDownloadUrlById(ctx context.Context, param map[string]interface{}) (backupDownloadUrl *cynosdb.DescribeBackupDownloadUrlResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeBackupDownloadUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "BackupId" {
			request.BackupId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeBackupDownloadUrl(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	backupDownloadUrl = response.Response

	return
}

func (me *CynosdbService) DescribeCynosdbBinlogDownloadUrlByFilter(ctx context.Context, param map[string]interface{}) (binlogDownloadUrl *cynosdb.DescribeBinlogDownloadUrlResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeBinlogDownloadUrlRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "BinlogId" {
			request.BinlogId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeBinlogDownloadUrl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	binlogDownloadUrl = response.Response

	return
}

func (me *CynosdbService) DescribeCynosdbClusterDetailDatabasesByFilter(ctx context.Context, param map[string]interface{}) (clusterDetailDatabases []*cynosdb.DbInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeClusterDetailDatabasesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "DbName" {
			request.DbName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCynosdbClient().DescribeClusterDetailDatabases(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DbInfos) < 1 {
			break
		}

		clusterDetailDatabases = append(clusterDetailDatabases, response.Response.DbInfos...)
		if len(response.Response.DbInfos) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbClusterParamLogsByFilter(ctx context.Context, param map[string]interface{}) (clusterParamLogs []*cynosdb.ClusterParamModifyLog, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeClusterParamLogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseCynosdbClient().DescribeClusterParamLogs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterParamLogs) < 1 {
			break
		}

		clusterParamLogs = append(clusterParamLogs, response.Response.ClusterParamLogs...)
		if len(response.Response.ClusterParamLogs) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbClusterByFilter(ctx context.Context, param map[string]interface{}) (cluster []*cynosdb.DatabaseTables, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewSearchClusterTablesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "Database" {
			request.Database = v.(*string)
		}
		if k == "Table" {
			request.Table = v.(*string)
		}
		if k == "TableType" {
			request.TableType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().SearchClusterTables(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Tables) < 1 {
		return
	}

	cluster = response.Response.Tables

	return
}

func (me *CynosdbService) DescribeCynosdbDescribeInstanceSlowQueriesByFilter(ctx context.Context, param map[string]interface{}) (describeInstanceSlowQueries []*cynosdb.BinlogItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeBinlogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
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
		response, err := me.client.UseCynosdbClient().DescribeBinlogs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Binlogs) < 1 {
			break
		}

		describeInstanceSlowQueries = append(describeInstanceSlowQueries, response.Response.Binlogs...)
		if len(response.Response.Binlogs) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbDescribeInstanceErrorLogsByFilter(ctx context.Context, param map[string]interface{}) (describeInstanceErrorLogs []*cynosdb.CynosdbErrorLogItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeInstanceErrorLogsRequest()
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
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*string)
		}
		if k == "LogLevels" {
			request.LogLevels = v.([]*string)
		}
		if k == "KeyWords" {
			request.KeyWords = v.([]*string)
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
		response, err := me.client.UseCynosdbClient().DescribeInstanceErrorLogs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ErrorLogs) < 1 {
			break
		}

		describeInstanceErrorLogs = append(describeInstanceErrorLogs, response.Response.ErrorLogs...)
		if len(response.Response.ErrorLogs) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbAccountAllGrantPrivilegesByFilter(ctx context.Context, param map[string]interface{}) (accountAllGrantPrivileges *cynosdb.DescribeAccountAllGrantPrivilegesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeAccountAllGrantPrivilegesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}

		if k == "Account" {
			request.Account = v.(*cynosdb.InputAccount)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeAccountAllGrantPrivileges(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	accountAllGrantPrivileges = response.Response

	return
}

func (me *CynosdbService) DescribeCynosdbProjectSecurityGroupsByFilter(ctx context.Context, param map[string]interface{}) (projectSecurityGroups []*cynosdb.SecurityGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeProjectSecurityGroupsRequest()
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
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
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
		response, err := me.client.UseCynosdbClient().DescribeProjectSecurityGroups(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Groups) < 1 {
			break
		}

		projectSecurityGroups = append(projectSecurityGroups, response.Response.Groups...)
		if len(response.Response.Groups) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbResourcePackageSaleSpecsByFilter(ctx context.Context, param map[string]interface{}) (resourcePackageSaleSpecs []*cynosdb.SalePackageSpec, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeResourcePackageSaleSpecRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceType" {
			request.InstanceType = v.(*string)
		}
		if k == "PackageRegion" {
			request.PackageRegion = v.(*string)
		}
		if k == "PackageType" {
			request.PackageType = v.(*string)
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
		response, err := me.client.UseCynosdbClient().DescribeResourcePackageSaleSpec(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Detail) < 1 {
			break
		}

		resourcePackageSaleSpecs = append(resourcePackageSaleSpecs, response.Response.Detail...)
		if len(response.Response.Detail) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbRollbackTimeRangeByFilter(ctx context.Context, param map[string]interface{}) (rollbackTimeRange *cynosdb.DescribeRollbackTimeRangeResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeRollbackTimeRangeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeRollbackTimeRange(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	rollbackTimeRange = response.Response

	return
}

func (me *CynosdbService) DescribeCynosdbRollbackTimeValidityByFilter(ctx context.Context, param map[string]interface{}) (rollbackTimeValidity *cynosdb.DescribeRollbackTimeValidityResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeRollbackTimeValidityRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "ExpectTime" {
			request.ExpectTime = v.(*string)
		}
		if k == "ExpectTimeThresh" {
			request.ExpectTimeThresh = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeRollbackTimeValidity(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	rollbackTimeValidity = response.Response

	return
}

func (me *CynosdbService) DescribeCynosdbResourcePackageListByFilter(ctx context.Context, param map[string]interface{}) (resourcePackageList []*cynosdb.Package, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeResourcePackageListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "PackageId" {
			request.PackageId = v.([]*string)
		}
		if k == "PackageName" {
			request.PackageName = v.([]*string)
		}
		if k == "PackageType" {
			request.PackageType = v.([]*string)
		}
		if k == "PackageRegion" {
			request.PackageRegion = v.([]*string)
		}
		if k == "Status" {
			request.Status = v.([]*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.([]*string)
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
		response, err := me.client.UseCynosdbClient().DescribeResourcePackageList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Detail) < 1 {
			break
		}

		resourcePackageList = append(resourcePackageList, response.Response.Detail...)
		if len(response.Response.Detail) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbClusterResourcePackagesAttachmentById(ctx context.Context, clusterId string) (clusterResourcePackagesAttachment *cynosdb.CynosdbClusterDetail, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeClusterDetailRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeClusterDetail(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Detail == nil {
		return
	}

	clusterResourcePackagesAttachment = response.Response.Detail
	return
}

func (me *CynosdbService) DeleteCynosdbClusterResourcePackagesAttachmentById(ctx context.Context, clusterId string, packageIdsSet []*string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewUnbindClusterResourcePackagesRequest()
	request.ClusterId = &clusterId
	request.PackageIds = packageIdsSet

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().UnbindClusterResourcePackages(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) DeleteCynosdbProxyById(ctx context.Context, clusterId string) (flowId *int64, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewCloseProxyRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().CloseProxy(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return response.Response.FlowId, nil
}

func (me *CynosdbService) DescribeCynosdbProxyById(ctx context.Context, clusterId, proxyGroupId string) (proxy *cynosdb.DescribeProxiesResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeProxiesRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeProxies(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	proxy = response.Response

	if proxyGroupId != "" {
		for _, proxyGroupRwInfo := range proxy.ProxyGroupInfos {
			proxyGroup := proxyGroupRwInfo.ProxyGroup
			if proxyGroupId == *proxyGroup.ProxyGroupId {
				proxy.ProxyGroupInfos = []*cynosdb.ProxyGroupInfo{proxyGroupRwInfo}
			}
		}
	}

	return
}

func (me *CynosdbService) DescribeCynosdbProxyNodeByFilter(ctx context.Context, param map[string]interface{}) (proxyNode []*cynosdb.ProxyNodeInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeProxyNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "OrderByType" {
			request.OrderByType = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*cynosdb.QueryFilter)
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
		response, err := me.client.UseCynosdbClient().DescribeProxyNodes(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ProxyNodeInfos) < 1 {
			break
		}

		proxyNode = append(proxyNode, response.Response.ProxyNodeInfos...)
		if len(response.Response.ProxyNodeInfos) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CynosdbService) DescribeCynosdbProxyVersionByFilter(ctx context.Context, param map[string]interface{}) (proxyVersion *cynosdb.DescribeSupportProxyVersionResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeSupportProxyVersionRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "ProxyGroupId" {
			request.ProxyGroupId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeSupportProxyVersion(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	proxyVersion = response.Response

	return
}

func (me *CynosdbService) CynosdbInstanceIsolateStateRefreshFunc(clusterId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, detail, _, err := me.DescribeClusterById(ctx, clusterId)

		if err != nil {
			return nil, "", err
		}

		if object == nil || object.Status == nil {
			return &cynosdb.CynosdbCluster{}, "isolated", nil
		}

		if object != nil || object.Status != nil {
			return object, helper.PString(object.Status), nil
		}

		return detail, helper.PString(detail.Status), nil
	}
}

func (me *CynosdbService) CynosdbInstanceOfflineStateRefreshFunc(clusterId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, detail, _, err := me.DescribeClusterById(ctx, clusterId)

		if err != nil {
			return nil, "", err
		}

		if object == nil || object.Status == nil {
			return &cynosdb.CynosdbCluster{}, "offlined", nil
		}

		if object != nil || object.Status != nil {
			return object, helper.PString(object.Status), nil
		}

		return detail, helper.PString(detail.Status), nil
	}
}

func (me *CynosdbService) DescribeCynosdbZoneByFilter(ctx context.Context, param map[string]interface{}) (zone []*cynosdb.SaleRegion, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = cynosdb.NewDescribeZonesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "IncludeVirtualZones" {
			request.IncludeVirtualZones = v.(*bool)
		}
		if k == "ShowPermission" {
			request.ShowPermission = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCynosdbClient().DescribeZones(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.RegionSet) < 1 {
		return
	}

	zone = response.Response.RegionSet

	return
}

func (me *CynosdbService) DescribeCynosdbAccountById(ctx context.Context, clusterId string, accountName string, host string) (account *cynosdb.Account, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeAccountsRequest()
	request.ClusterId = &clusterId
	request.AccountNames = []*string{&accountName}
	request.Hosts = []*string{&host}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AccountSet) < 1 {
		return
	}

	account = response.Response.AccountSet[0]
	return
}

func (me *CynosdbService) DeleteCynosdbAccountById(ctx context.Context, clusterId string, accountNames string, host string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDeleteAccountsRequest()
	request.ClusterId = &clusterId
	request.Accounts = []*cynosdb.InputAccount{
		{
			AccountName: &accountNames,
			Host:        &host,
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DeleteAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) DescribeCynosdbAccountPrivilegesById(ctx context.Context, clusterId string, accountName string, host string) (accountPrivileges *cynosdb.DescribeAccountAllGrantPrivilegesResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeAccountAllGrantPrivilegesRequest()
	request.ClusterId = &clusterId
	request.Account = &cynosdb.InputAccount{
		AccountName: &accountName,
		Host:        &host,
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeAccountAllGrantPrivileges(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accountPrivileges = response.Response
	return
}

func (me *CynosdbService) DeleteCynosdbAccountPrivilegesById(ctx context.Context, clusterId string, accountName string, host string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewRevokeAccountPrivilegesRequest()
	request.ClusterId = &clusterId
	request.Account = &cynosdb.InputAccount{
		AccountName: &accountName,
		Host:        &host,
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().RevokeAccountPrivileges(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) DescribeCynosdbBinlogSaveDaysById(ctx context.Context, clusterId string) (binlogSaveDays *int64, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeBinlogSaveDaysRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeBinlogSaveDays(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	binlogSaveDays = response.Response.BinlogSaveDays
	return
}

func (me *CynosdbService) DescribeCynosdbClusterDatabasesById(ctx context.Context, clusterId string, dbName string) (clusterDatabases *cynosdb.DbInfo, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeClusterDetailDatabasesRequest()
	request.ClusterId = &clusterId
	request.DbName = &dbName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeClusterDetailDatabases(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || len(response.Response.DbInfos) < 1 {
		return
	}

	clusterDatabases = response.Response.DbInfos[0]

	return
}

func (me *CynosdbService) DeleteCynosdbClusterDatabasesById(ctx context.Context, clusterId, dbName string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDeleteClusterDatabaseRequest()
	request.ClusterId = &clusterId
	request.DbNames = []*string{&dbName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DeleteClusterDatabase(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) DescribeCynosdbClusterPasswordComplexityById(ctx context.Context, clusterId string) (clusterPasswordComplexity *cynosdb.DescribeClusterPasswordComplexityResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeClusterPasswordComplexityRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeClusterPasswordComplexity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	clusterPasswordComplexity = response.Response
	return
}

func (me *CynosdbService) DeleteCynosdbClusterPasswordComplexityById(ctx context.Context, clusterId string) (flowId int64, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewCloseClusterPasswordComplexityRequest()
	request.ClusterIds = []*string{&clusterId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().CloseClusterPasswordComplexity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	flowId = *response.Response.FlowId
	return
}

func (me *CynosdbService) DescribeFlow(ctx context.Context, flowId int64) (ok bool, errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewDescribeFlowRequest()
	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCynosdbClient().DescribeFlow(request)
	if err != nil {
		errRet = err
		return
	}
	if *response.Response.Status == 2 {
		return
	}
	if *response.Response.Status == 0 {
		ok = true
		return
	}
	errRet = fmt.Errorf("redis task exe fail, task status is %v", *response.Response.Status)
	return
}

func (me *CynosdbService) CopyClusterPasswordComplexity(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewCopyClusterPasswordComplexityRequest()
	request.ClusterIds = []*string{&clusterId}
	request.SourceClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().CopyClusterPasswordComplexity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	flowId := *response.Response.FlowId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := me.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("create cynosdb clusterPasswordComplexity is processing"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb clusterPasswordComplexity fail, reason:%s\n", logId, err.Error())
		errRet = err
		return
	}

	return
}

func (me *CynosdbService) DescribeCynosdbInstanceParamById(ctx context.Context, clusterId string, instanceId string) (instanceParam *cynosdb.InstanceParamItem, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeInstanceParamsRequest()
	request.ClusterId = &clusterId
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) < 1 {
		return
	}

	instanceParam = response.Response.Items[0]
	return
}

func (me *CynosdbService) DeleteCynosdbWanById(ctx context.Context, instanceGrpId string) (flowId int64, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewCloseWanRequest()
	request.InstanceGrpId = &instanceGrpId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().CloseWan(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	flowId = *response.Response.FlowId
	return
}

func (me *CynosdbService) DescribeCynosdbParamTemplateById(ctx context.Context, templateId int64) (paramTemplate *cynosdb.DescribeParamTemplateDetailResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeParamTemplateDetailRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeParamTemplateDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	paramTemplate = response.Response
	return
}

func (me *CynosdbService) DeleteCynosdbParamTemplateById(ctx context.Context, templateId int64) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDeleteParamTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DeleteParamTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CynosdbService) SetRenewFlag(ctx context.Context, instanceId string, autoRenewFlag int64) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewSetRenewFlagRequest()
	request.ResourceIds = []*string{&instanceId}
	request.AutoRenewFlag = &autoRenewFlag

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().SetRenewFlag(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) ModifyClusterName(ctx context.Context, clusterId string, clusterName string) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewModifyClusterNameRequest()
	request.ClusterId = &clusterId
	request.ClusterName = &clusterName

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().ModifyClusterName(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) ModifyClusterStorage(ctx context.Context, clusterId string, newStorageLimit int64, oldStorageLimit int64) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewModifyClusterStorageRequest()
	request.ClusterId = &clusterId
	request.NewStorageLimit = &newStorageLimit
	request.OldStorageLimit = &oldStorageLimit
	request.DealMode = helper.IntInt64(0)

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseCynosdbClient().ModifyClusterStorage(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *CynosdbService) SwitchClusterVpc(ctx context.Context, clusterId string, vpcId string, subnetId string, oldIpReserveHours int64) (errRet error) {
	logId := getLogId(ctx)
	request := cynosdb.NewSwitchClusterVpcRequest()
	request.ClusterId = &clusterId
	request.UniqVpcId = &vpcId
	request.UniqSubnetId = &subnetId
	request.OldIpReserveHours = &oldIpReserveHours

	var flowId int64
	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet := me.client.UseCynosdbClient().SwitchClusterVpc(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		flowId = *response.Response.FlowId
		return nil
	})
	if errRet != nil {
		return
	}

	err := resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := me.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("update cynosdb SwitchClusterVpc is processing"))
		}
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb SwitchClusterVpc fail, reason:%s\n", logId, err.Error())
		errRet = err
		return
	}

	return
}

func (me *CynosdbService) DescribeCynosdbResourcePackageById(ctx context.Context, packageId string) (resourcePackage *cynosdb.PackageDetail, errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewDescribeResourcePackageDetailRequest()
	request.PackageId = &packageId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().DescribeResourcePackageDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Detail) < 1 {
		return
	}

	resourcePackage = response.Response.Detail[0]
	return
}

func (me *CynosdbService) DeleteCynosdbResourcePackageById(ctx context.Context, packageId string) (errRet error) {
	logId := getLogId(ctx)

	request := cynosdb.NewRefundResourcePackageRequest()
	request.PackageId = &packageId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCynosdbClient().RefundResourcePackage(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
