package tencentcloud

import (
	"context"
	"log"

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
func (me *TdcpgService) DescribeTdcpgCluster(ctx context.Context, clusterId string) (cluster *tdcpg.DescribeClustersResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdcpg.NewDescribeClustersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.Filters = []*tdcpg.Filter{
		{
			Name:   helper.String(TDCPG_CLUSTER_FILTER_ID),
			Values: []*string{&clusterId},
		},
	}

	response, err := me.client.UseTdcpgClient().DescribeClusters(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	cluster = response.Response
	return
}

func (me *TdcpgService) DeleteTdcpgClusterById(ctx context.Context, clusterId string) (errRet error) {
	logId := getLogId(ctx)

	request := tdcpg.NewDeleteClusterRequest()

	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdcpgClient().DeleteCluster(request)
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
		logId   = getLogId(ctx)
		request = tdcpg.NewDescribeClusterInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = clusterId
	request.Filters = []*tdcpg.Filter{
		{
			Name:   helper.String(TDCPG_INSTANCE_FILTER_ID),
			Values: []*string{instanceId},
		},
	}

	response, err := me.client.UseTdcpgClient().DescribeClusterInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	instance = response.Response
	return
}

func (me *TdcpgService) DescribeTdcpgResourceByDealName(ctx context.Context, dealName []*string) (resourceInfo []*tdcpg.ResourceIdInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tdcpg.NewDescribeResourcesByDealNameRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.DealName = dealName

	response, err := me.client.UseTdcpgClient().DescribeResourcesByDealName(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	resourceInfo = response.Response.ResourceIdInfoSet
	return
}

func (me *TdcpgService) DeleteTdcpgInstanceById(ctx context.Context, clusterId, instanceId *string) (errRet error) {
	logId := getLogId(ctx)

	request := tdcpg.NewDeleteClusterInstancesRequest()

	request.ClusterId = clusterId
	request.InstanceIdSet[0] = instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTdcpgClient().DeleteClusterInstances(request)
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
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

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
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
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
		if len(response.Response.ClusterSet) < int(pageSize) {
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
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

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

		response, err := me.client.UseTdcpgClient().DescribeClusterInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
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
		if len(response.Response.InstanceSet) < int(pageSize) {
			break
		}
		currNumber++
	}
	return
}
