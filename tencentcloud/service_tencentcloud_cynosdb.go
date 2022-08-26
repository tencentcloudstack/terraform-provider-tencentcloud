package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CynosdbService struct {
	client *connectivity.TencentCloudClient
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

func (me *CynosdbService) DescribeClusterById(ctx context.Context, clusterId string) (renewFlag int64, clusterInfo *cynosdb.CynosdbClusterDetail, has bool, errRet error) {
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
		if *clusters[0].Status == CYNOSDB_STATUS_ISOLATED || *clusters[0].Status == CYNOSDB_STATUS_OFFLINE || *clusters[0].Status == CYNOSDB_STATUS_DELETED {
			notExist = true
			return nil
		} else if *clusters[0].Status == CYNOSDB_STATUS_RUNNING {
			renewFlag = *clusters[0].RenewFlag
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cynosdb cluster %s is still in processing", clusterId))
		}
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
	items = response.Response.Items

	return
}
