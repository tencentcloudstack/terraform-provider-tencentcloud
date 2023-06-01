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
