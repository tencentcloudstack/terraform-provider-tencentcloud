package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	TDCPG_CLUSTER_FILTER_ID         = "ClusterId"
	TDCPG_CLUSTER_FILTER_NAME       = "ClusterName"
	TDCPG_CLUSTER_FILTER_PROJECT_ID = "ProjectId"
	TDCPG_CLUSTER_FILTER_STATUS     = "Status"
	TDCPG_CLUSTER_FILTER_PAY_MODE   = "PayMode"

	TDCPG_INSTANCE_FILTER_ID          = "InstanceId"
	TDCPG_INSTANCE_FILTER_NAME        = "InstanceName"
	TDCPG_INSTANCE_FILTER_ENDPOINT_ID = "EndpointId"
	TDCPG_INSTANCE_FILTER_STATUS      = "Status"
	TDCPG_INSTANCE_FILTER_TYPE        = "InstanceType"
)

type TdcpgService struct {
	client *connectivity.TencentCloudClient
}

// tdcpg resource
func (me *TdcpgService) DescribeTdcpgCluster(ctx context.Context, clusterId *string) (cluster *tdcpg.DescribeClustersResponseParams, errRet error) {
	var (
		logId = getLogId(ctx)
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, reason[%s]\n", logId, "DescribeTdcpgCluster", errRet.Error())
		}
	}()

	paramMap := make(map[string]interface{})
	paramMap["cluster_id"] = clusterId

	result, err := me.DescribeTdcpgClustersByFilter(ctx, paramMap)
	if err != nil {
		errRet = err
		return
	}
	cluster = &tdcpg.DescribeClustersResponseParams{
		ClusterSet: result,
	}
	return
}

func (me *TdcpgService) IsolateTdcpgInstanceById(ctx context.Context, clusterId, instanceId *string) (errRet error) {
	logId := getLogId(ctx)

	request := tdcpg.NewIsolateClusterInstancesRequest()
	request.ClusterId = clusterId
	request.InstanceIdSet = []*string{instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "isolate tdcpg instance object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdcpgClient().IsolateClusterInstances(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdcpgService) DeleteTdcpgClusterById(ctx context.Context, clusterId *string) (errRet error) {
	logId := getLogId(ctx)
	var status string
	if err := me.IsolateTdcpgClusterById(ctx, clusterId); err != nil {
		return err
	}

	// polling the cluster's status to isolated
	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		result, err := me.DescribeTdcpgCluster(ctx, clusterId)
		if err != nil {
			return retryError(err)
		}
		if result != nil {
			status = *result.ClusterSet[0].Status
			if status == "isolated" || status == "deleted" {
				return nil
			}
			if status == "isolating" {
				return resource.RetryableError(fmt.Errorf("cluster status still on isolating, retry..."))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s query tdcpg cluster failed, reason:%+v", logId, err)
		return err
	}
	if status == "deleted" {
		// do not need delete
		return nil
	}

	request := tdcpg.NewDeleteClusterRequest()
	request.ClusterId = clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete tdcpg cluster object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdcpgClient().DeleteCluster(request)
	if err != nil {
		errRet = err
		return err
	}

	// wait the cluster to be deleted
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		result, err := me.DescribeTdcpgCluster(ctx, clusterId)
		if err != nil {
			return retryError(err)
		}
		if result != nil {
			status = *result.ClusterSet[0].Status
			if status == "deleted" {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		errRet = err
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdcpgService) DescribeTdcpgInstance(ctx context.Context, clusterId, instanceId *string) (instance *tdcpg.DescribeClusterInstancesResponseParams, errRet error) {
	var (
		logId = getLogId(ctx)
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, reason[%s]\n", logId, "DescribeTdcpgInstance", errRet.Error())
		}
	}()

	paramMap := make(map[string]interface{})
	paramMap["instance_id"] = instanceId

	result, err := me.DescribeTdcpgInstancesByFilter(ctx, clusterId, paramMap)
	if err != nil {
		errRet = err
		return
	}
	instance = &tdcpg.DescribeClusterInstancesResponseParams{
		InstanceSet: result,
	}
	return
}

func (me *TdcpgService) DescribeTdcpgResourceByDealName(ctx context.Context, dealNames []*string) (resourceInfo []*tdcpg.ResourceIdInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdcpg.NewDescribeResourcesByDealNameRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for _, dealname := range dealNames {
		request.DealName = dealname

		log.Printf("[DEBUG]%s dealName:[%v]\n", logId, *dealname)

		response, err := me.client.UseTdcpgClient().DescribeResourcesByDealName(request)
		if err != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		resourceInfo = response.Response.ResourceIdInfoSet
	}
	return
}

func (me *TdcpgService) IsolateTdcpgClusterById(ctx context.Context, clusterId *string) (errRet error) {
	logId := getLogId(ctx)

	request := tdcpg.NewIsolateClusterRequest()
	request.ClusterId = clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "isolate tdcpg cluster object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdcpgClient().IsolateCluster(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TdcpgService) DeleteTdcpgInstanceById(ctx context.Context, clusterId, instanceId *string) (errRet error) {
	logId := getLogId(ctx)
	var status string
	if err := me.IsolateTdcpgInstanceById(ctx, clusterId, instanceId); err != nil {
		return err
	}

	// polling the instance's status to isolated
	err := resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		result, err := me.DescribeTdcpgInstance(ctx, clusterId, instanceId)
		if err != nil {
			return retryError(err)
		}
		if result != nil {
			status = *result.InstanceSet[0].Status
			if status == "isolated" || status == "deleted" {
				return nil
			}

			if status == "isolating" {
				return resource.RetryableError(fmt.Errorf("instance status still on isolating, retry..."))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s query tdcpg instance failed, reason:%+v", logId, err)
		return err
	}
	if status == "deleted" {
		// do not need delete
		return nil
	}

	request := tdcpg.NewDeleteClusterInstancesRequest()

	request.ClusterId = clusterId
	request.InstanceIdSet = []*string{instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete tdcpg instance object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdcpgClient().DeleteClusterInstances(request)
	if err != nil {
		errRet = err
		return err
	}

	// wait the instance to be deleted
	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		result, err := me.DescribeTdcpgInstance(ctx, clusterId, instanceId)
		if err != nil {
			return retryError(err)
		}
		if result != nil {
			status = *result.InstanceSet[0].Status
			if status == "deleted" {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		errRet = err
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// tdcpg data-source
func (me *TdcpgService) DescribeTdcpgClustersByFilter(ctx context.Context, param map[string]interface{}) (clusters []*tdcpg.Cluster, errRet error) {
	var (
		logId      = getLogId(ctx)
		request    = tdcpg.NewDescribeClustersRequest()
		indx       = 0
		currNumber = 1
		pageSize   = 20
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*tdcpg.Filter, len(param))
	for k, v := range param {
		if k == "cluster_id" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_CLUSTER_FILTER_ID),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "cluster_name" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_CLUSTER_FILTER_NAME),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "status" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_CLUSTER_FILTER_STATUS),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "pay_mode" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_CLUSTER_FILTER_PAY_MODE),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "project_id" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_CLUSTER_FILTER_PROJECT_ID),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}
	}
	ratelimit.Check(request.GetAction())

	for {
		request.PageNumber = helper.IntUint64(currNumber)
		request.PageSize = helper.IntUint64(pageSize)

		response, err := me.client.UseTdcpgClient().DescribeClusters(request)
		if err != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ClusterSet) < 1 {
			break
		}

		clusters = append(clusters, response.Response.ClusterSet...)
		if len(response.Response.ClusterSet) < pageSize {
			break
		}
		currNumber++
	}
	return
}

func (me *TdcpgService) DescribeTdcpgInstancesByFilter(ctx context.Context, clusterId *string, param map[string]interface{}) (instances []*tdcpg.Instance, errRet error) {
	var (
		logId      = getLogId(ctx)
		request    = tdcpg.NewDescribeClusterInstancesRequest()
		indx       = 0
		currNumber = 1
		pageSize   = 20
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = make([]*tdcpg.Filter, len(param))
	for k, v := range param {
		if k == "instance_id" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_INSTANCE_FILTER_ID),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "instance_name" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_INSTANCE_FILTER_NAME),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "status" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_INSTANCE_FILTER_STATUS),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}

		if k == "instance_type" {
			request.Filters[indx] = &tdcpg.Filter{
				Name:   helper.String(TDCPG_INSTANCE_FILTER_TYPE),
				Values: []*string{v.(*string)},
			}
			indx++
			continue
		}
	}
	ratelimit.Check(request.GetAction())

	for {
		request.PageNumber = helper.IntUint64(currNumber)
		request.PageSize = helper.IntUint64(pageSize)
		request.ClusterId = clusterId

		response, err := me.client.UseTdcpgClient().DescribeClusterInstances(request)
		if err != nil {
			log.Printf("[CRITICAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < pageSize {
			break
		}
		currNumber++
	}
	return
}
