package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TseService struct {
	client *connectivity.TencentCloudClient
}

func (me *TseService) DescribeTseInstanceById(ctx context.Context, instanceId string) (instance *tse.SREInstance, errRet error) {
	logId := getLogId(ctx)

	request := tse.NewDescribeSREInstancesRequest()
	filter := &tse.Filter{
		Name:   helper.String("InstanceId"),
		Values: []*string{&instanceId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeSREInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Content) < 1 {
		return
	}

	instance = response.Response.Content[0]
	return
}

func (me *TseService) CheckTseInstanceStatusById(ctx context.Context, instanceId, operate string) (errRet error) {
	logId := getLogId(ctx)

	err := resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		instance, e := me.DescribeTseInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if operate == "create" {
			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("instance %s instance not exists", instanceId))
			}

			if *instance.Status == "creating" || *instance.Status == "restarting" {
				return resource.RetryableError(fmt.Errorf("create instance status is %v,start retrying ...", *instance.Status))
			}
			if *instance.Status == "running" {
				return nil
			}
		}

		if operate == "update" {
			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("instance %s instance not exists", instanceId))
			}

			if *instance.Status == "updating" || *instance.Status == "restarting" {
				return resource.RetryableError(fmt.Errorf("update instance status is %v,start retrying ...", *instance.Status))
			}
			if *instance.Status == "running" {
				return nil
			}
		}

		if operate == "delete" {
			if instance == nil {
				return nil
			}

			if *instance.Status == "destroying" {
				return resource.RetryableError(fmt.Errorf("delete instance status is %v,start retrying ...", *instance.Status))
			}
		}

		return resource.NonRetryableError(fmt.Errorf("instance status is %v,we won't wait for it finish", *instance.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb fail, reason:%s\n ", logId, err.Error())
		errRet = err
		return
	}

	return
}

func (me *TseService) DeleteTseInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := tse.NewDeleteEngineRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteEngine(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseAccessAddressByFilter(ctx context.Context, param map[string]interface{}) (accessAddress *tse.DescribeSREInstanceAccessAddressResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tse.NewDescribeSREInstanceAccessAddressRequest()
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
		if k == "VpcId" {
			request.VpcId = v.(*string)
		}
		if k == "SubnetId" {
			request.SubnetId = v.(*string)
		}
		if k == "Workload" {
			request.Workload = v.(*string)
		}
		if k == "EngineRegion" {
			request.EngineRegion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTseClient().DescribeSREInstanceAccessAddress(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	accessAddress = response.Response

	return
}

func (me *TseService) DescribeTseNacosReplicasByFilter(ctx context.Context, param map[string]interface{}) (nacosReplicas []*tse.NacosReplica, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tse.NewDescribeNacosReplicasRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeNacosReplicas(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Replicas) < 1 {
			break
		}
		nacosReplicas = append(nacosReplicas, response.Response.Replicas...)
		if len(response.Response.Replicas) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseZookeeperReplicasByFilter(ctx context.Context, param map[string]interface{}) (zookeeperReplicas []*tse.ZookeeperReplica, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tse.NewDescribeZookeeperReplicasRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeZookeeperReplicas(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Replicas) < 1 {
			break
		}
		zookeeperReplicas = append(zookeeperReplicas, response.Response.Replicas...)
		if len(response.Response.Replicas) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseZookeeperServerInterfacesByFilter(ctx context.Context, param map[string]interface{}) (zookeeperServerInterfaces []*tse.ZookeeperServerInterface, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tse.NewDescribeZookeeperServerInterfacesRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeZookeeperServerInterfaces(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Content) < 1 {
			break
		}
		zookeeperServerInterfaces = append(zookeeperServerInterfaces, response.Response.Content...)
		if len(response.Response.Content) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseNacosServerInterfacesByFilter(ctx context.Context, instanceId string) (nacosServerInterfaces []*tse.NacosServerInterface, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tse.NewDescribeNacosServerInterfacesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeNacosServerInterfaces(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Content) < 1 {
			break
		}
		nacosServerInterfaces = append(nacosServerInterfaces, response.Response.Content...)
		if len(response.Response.Content) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseGatewayNodesByFilter(ctx context.Context, param map[string]interface{}) (gatewayNodes []*tse.CloudNativeAPIGatewayNode, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewayNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "GroupId" {
			request.GroupId = v.(*string)
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
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayNodes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.NodeList) < 1 {
			break
		}
		gatewayNodes = append(gatewayNodes, response.Response.Result.NodeList...)
		if len(response.Response.Result.NodeList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
