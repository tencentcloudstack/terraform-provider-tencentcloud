package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CdwchService struct {
	client *connectivity.TencentCloudClient
}

func (me *CdwchService) DescribeInstance(ctx context.Context, instanceId string) (InstanceInfo *cdwch.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeInstanceRequest()
	request.IsOpenApi = helper.Bool(true)
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	InstanceInfo = response.Response.InstanceInfo

	return
}

func (me *CdwchService) DestroyInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDestroyInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DestroyInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) ResizeDisk(ctx context.Context, instanceId string, nodeType string, resizeDisk int) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewResizeDiskRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Type = &nodeType
	request.DiskSize = helper.IntInt64(resizeDisk)
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().ResizeDisk(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) ScaleUpInstance(ctx context.Context, instanceId, nodeType, specName string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewScaleUpInstanceRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.ScaleUpEnableRolling = helper.Bool(true)
	request.Type = &nodeType
	request.SpecName = &specName
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().ScaleUpInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) ScaleOutInstance(ctx context.Context, instanceId string, nodeType string, scaleOutCluster string, nodeCount int, userSubnetIPNum int, shardIps []*string) (errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewScaleOutInstanceRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Type = &nodeType
	request.NodeCount = helper.IntInt64(nodeCount)
	request.ScaleOutCluster = &scaleOutCluster
	request.UserSubnetIPNum = helper.IntInt64(userSubnetIPNum)
	if shardIps != nil {
		request.ReduceShardInfo = shardIps
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().ScaleOutInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response.Response.ErrorMsg != nil && *response.Response.ErrorMsg != "" {
		errRet = fmt.Errorf(*response.Response.ErrorMsg)
	}
	return
}

func (me *CdwchService) DescribeInstanceClusters(ctx context.Context, instanceId string) (clusterInfos []*cdwch.ClusterInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeInstanceClustersRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstanceClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusterInfos = response.Response.Clusters

	return
}

func (me *CdwchService) DescribeInstancesNew(ctx context.Context, instanceId string) (instancesList []*cdwch.InstanceInfo, errRet error) {
	logId := getLogId(ctx)

	request := cdwch.NewDescribeInstancesNewRequest()
	request.SearchInstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCdwchClient().DescribeInstancesNew(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instancesList = response.Response.InstancesList

	return
}
