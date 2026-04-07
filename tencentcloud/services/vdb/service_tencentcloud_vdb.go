package vdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// VdbClientInterface abstracts the VDB SDK client methods used by VdbService.
// This interface enables mock-based unit testing without calling real cloud APIs.
type VdbClientInterface interface {
	DescribeInstancesWithContext(ctx context.Context, request *vdb.DescribeInstancesRequest) (response *vdb.DescribeInstancesResponse, err error)
	DescribeInstanceNodesWithContext(ctx context.Context, request *vdb.DescribeInstanceNodesRequest) (response *vdb.DescribeInstanceNodesResponse, err error)
	DescribeDBSecurityGroupsWithContext(ctx context.Context, request *vdb.DescribeDBSecurityGroupsRequest) (response *vdb.DescribeDBSecurityGroupsResponse, err error)
	CreateInstanceWithContext(ctx context.Context, request *vdb.CreateInstanceRequest) (response *vdb.CreateInstanceResponse, err error)
	ScaleUpInstanceWithContext(ctx context.Context, request *vdb.ScaleUpInstanceRequest) (response *vdb.ScaleUpInstanceResponse, err error)
	ScaleOutInstanceWithContext(ctx context.Context, request *vdb.ScaleOutInstanceRequest) (response *vdb.ScaleOutInstanceResponse, err error)
	IsolateInstanceWithContext(ctx context.Context, request *vdb.IsolateInstanceRequest) (response *vdb.IsolateInstanceResponse, err error)
	DestroyInstancesWithContext(ctx context.Context, request *vdb.DestroyInstancesRequest) (response *vdb.DestroyInstancesResponse, err error)
	AssociateSecurityGroupsWithContext(ctx context.Context, request *vdb.AssociateSecurityGroupsRequest) (response *vdb.AssociateSecurityGroupsResponse, err error)
	DisassociateSecurityGroupsWithContext(ctx context.Context, request *vdb.DisassociateSecurityGroupsRequest) (response *vdb.DisassociateSecurityGroupsResponse, err error)
	ModifyDBInstanceSecurityGroupsWithContext(ctx context.Context, request *vdb.ModifyDBInstanceSecurityGroupsRequest) (response *vdb.ModifyDBInstanceSecurityGroupsResponse, err error)
}

type VdbService struct {
	client    *connectivity.TencentCloudClient
	vdbClient VdbClientInterface
}

func NewVdbService(client *connectivity.TencentCloudClient) VdbService {
	return VdbService{client: client}
}

// getVdbClient returns the VDB client, preferring the injected mock if available.
func (me *VdbService) getVdbClient() VdbClientInterface {
	if me.vdbClient != nil {
		return me.vdbClient
	}
	return me.client.UseVdbV20230616Client()
}

func (me *VdbService) DescribeVdbInstanceById(ctx context.Context, instanceId string) (instance *vdb.InstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vdb.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *vdb.DescribeInstancesResponse
	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.getVdbClient().DescribeInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeInstances response is nil"))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	if response.Response.Items == nil || len(response.Response.Items) == 0 {
		return nil, nil
	}

	instance = response.Response.Items[0]
	return
}

func (me *VdbService) DescribeVdbInstanceNodesById(ctx context.Context, instanceId string) (nodes []*vdb.NodeInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vdb.NewDescribeInstanceNodesRequest()
	request.InstanceId = helper.String(instanceId)
	var limit int64 = 100
	request.Limit = &limit

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.getVdbClient().DescribeInstanceNodesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeInstanceNodes response is nil"))
		}

		nodes = result.Response.Items
		return nil
	})

	return
}

func (me *VdbService) DescribeDBSecurityGroupsByInstanceId(ctx context.Context, instanceId string) (groups []*vdb.SecurityGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vdb.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.getVdbClient().DescribeDBSecurityGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeDBSecurityGroups response is nil"))
		}

		groups = result.Response.Groups
		return nil
	})

	return
}

func (me *VdbService) WaitForInstanceStatus(ctx context.Context, instanceId string, targetStatus string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	return resource.Retry(timeout, func() *resource.RetryError {
		instance, e := me.DescribeVdbInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instance == nil {
			if strings.EqualFold(targetStatus, "isolated") || strings.EqualFold(targetStatus, "destroyed") {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("instance %s not found, waiting for %s", instanceId, targetStatus))
		}

		if instance.Status != nil {
			log.Printf("[DEBUG]%s instance %s current status: [%s] (waiting for %s)\n", logId, instanceId, *instance.Status, targetStatus)
			if strings.EqualFold(*instance.Status, targetStatus) {
				return nil
			}
		}

		status := "unknown"
		if instance.Status != nil {
			status = *instance.Status
		}
		return resource.RetryableError(fmt.Errorf("instance %s status is %s, waiting for %s", instanceId, status, targetStatus))
	})
}

func (me *VdbService) WaitForInstanceScaleUp(ctx context.Context, instanceId string, targetCpu float64, targetMemory float64, targetDiskSize uint64, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	return resource.Retry(timeout, func() *resource.RetryError {
		instance, e := me.DescribeVdbInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instance == nil {
			return resource.RetryableError(fmt.Errorf("instance %s not found, waiting for scale up", instanceId))
		}

		currentCpu := float64(0)
		if instance.Cpu != nil {
			currentCpu = *instance.Cpu
		}

		currentMemory := float64(0)
		if instance.Memory != nil {
			currentMemory = *instance.Memory
		}

		currentDisk := uint64(0)
		if instance.Disk != nil {
			currentDisk = *instance.Disk
		}

		log.Printf("[DEBUG]%s instance %s scale up progress: cpu=%.0f/%.0f, memory=%.0f/%.0f, disk=%d/%d\n",
			logId, instanceId, currentCpu, targetCpu, currentMemory, targetMemory, currentDisk, targetDiskSize)

		if currentCpu == targetCpu && currentMemory == targetMemory && currentDisk == targetDiskSize {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("instance %s scale up not complete: cpu=%.0f/%.0f, memory=%.0f/%.0f, disk=%d/%d",
			instanceId, currentCpu, targetCpu, currentMemory, targetMemory, currentDisk, targetDiskSize))
	})
}

func (me *VdbService) WaitForInstanceScaleOut(ctx context.Context, instanceId string, targetReplicaNum uint64, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	return resource.Retry(timeout, func() *resource.RetryError {
		instance, e := me.DescribeVdbInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instance == nil {
			return resource.RetryableError(fmt.Errorf("instance %s not found, waiting for scale out", instanceId))
		}

		currentReplicaNum := uint64(0)
		if instance.ReplicaNum != nil {
			currentReplicaNum = *instance.ReplicaNum
		}

		log.Printf("[DEBUG]%s instance %s scale out progress: worker_node_num=%d/%d\n",
			logId, instanceId, currentReplicaNum, targetReplicaNum)

		if currentReplicaNum == targetReplicaNum {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("instance %s scale out not complete: worker_node_num=%d/%d",
			instanceId, currentReplicaNum, targetReplicaNum))
	})
}

func (me *VdbService) WaitForInstanceNotFound(ctx context.Context, instanceId string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	return resource.Retry(timeout, func() *resource.RetryError {
		instance, e := me.DescribeVdbInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if instance == nil {
			return nil
		}

		status := "unknown"
		if instance.Status != nil {
			status = *instance.Status
		}
		log.Printf("[DEBUG]%s instance %s still exists, status: [%s], waiting for deletion\n", logId, instanceId, status)
		return resource.RetryableError(fmt.Errorf("instance %s still exists with status %s, waiting for deletion", instanceId, status))
	})
}

func (me *VdbService) WaitForSecurityGroupsMatch(ctx context.Context, instanceId string, targetSgIds []string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	targetSet := make(map[string]bool, len(targetSgIds))
	for _, id := range targetSgIds {
		targetSet[id] = true
	}

	return resource.Retry(timeout, func() *resource.RetryError {
		groups, e := me.DescribeDBSecurityGroupsByInstanceId(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		currentSet := make(map[string]bool, len(groups))
		for _, g := range groups {
			if g.SecurityGroupId != nil {
				currentSet[*g.SecurityGroupId] = true
			}
		}

		// Check if sets are equal
		if len(currentSet) == len(targetSet) {
			match := true
			for id := range targetSet {
				if !currentSet[id] {
					match = false
					break
				}
			}
			if match {
				return nil
			}
		}

		log.Printf("[DEBUG]%s instance %s security groups: current=%d, target=%d, waiting for match\n",
			logId, instanceId, len(currentSet), len(targetSet))
		return resource.RetryableError(fmt.Errorf("instance %s security groups not matched yet: current=%d, target=%d",
			instanceId, len(currentSet), len(targetSet)))
	})
}
