package ags

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	// AGS_SANDBOX_TOOL_DESCRIBE_LIMIT is the page size used when locating a
	// single sandbox tool by id. The API accepts at most 100 ids per request.
	AGS_SANDBOX_TOOL_DESCRIBE_LIMIT = 100

	// Sandbox tool status values documented by DescribeSandboxToolList.
	AGS_SANDBOX_TOOL_STATUS_ACTIVE   = "ACTIVE"
	AGS_SANDBOX_TOOL_STATUS_FAILED   = "FAILED"
	AGS_SANDBOX_TOOL_STATUS_CREATING = "CREATING"
	AGS_SANDBOX_TOOL_STATUS_DELETING = "DELETING"

	// Sandbox instance status values documented by DescribeSandboxInstances.
	AGS_SANDBOX_INSTANCE_STATUS_STARTING    = "STARTING"
	AGS_SANDBOX_INSTANCE_STATUS_RUNNING     = "RUNNING"
	AGS_SANDBOX_INSTANCE_STATUS_STOPPING    = "STOPPING"
	AGS_SANDBOX_INSTANCE_STATUS_STOPPED     = "STOPPED"
	AGS_SANDBOX_INSTANCE_STATUS_STOP_FAILED = "STOP_FAILED"
	AGS_SANDBOX_INSTANCE_STATUS_FAILED      = "FAILED"
	AGS_SANDBOX_INSTANCE_STATUS_RESIZING    = "RESIZING"
)

func NewAgsService(client *connectivity.TencentCloudClient) AgsService {
	return AgsService{client: client}
}

type AgsService struct {
	client *connectivity.TencentCloudClient
}

// CreateSandboxTool creates a sandbox tool and returns the id assigned by the
// cloud side. The id is the only resource identity: RequestId is only a request
// trace id and must never be used as an identity.
func (me *AgsService) CreateSandboxTool(ctx context.Context, request *agsv20250920.CreateSandboxToolRequest) (toolId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.CreateSandboxToolResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().CreateSandboxToolWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateSandboxTool returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ToolId == nil || *response.Response.ToolId == "" {
		errRet = fmt.Errorf("CreateSandboxTool returned an empty tool id")
		return
	}

	toolId = *response.Response.ToolId
	return
}

// DescribeSandboxToolById returns the sandbox tool with the given id.
//
// A nil result with a nil error means "queried successfully, no such tool"; an
// API failure is always reported as an error so the caller never confuses a
// failed request with a deleted resource.
func (me *AgsService) DescribeSandboxToolById(ctx context.Context, toolId string) (sandboxTool *agsv20250920.SandboxTool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDescribeSandboxToolListRequest()
	request.ToolIds = []*string{helper.String(toolId)}
	request.Limit = helper.Int64(AGS_SANDBOX_TOOL_DESCRIBE_LIMIT)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DescribeSandboxToolListResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DescribeSandboxToolListWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeSandboxToolList returned an empty response for sandbox tool %s", toolId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, tool := range response.Response.SandboxToolSet {
		if tool != nil && tool.ToolId != nil && *tool.ToolId == toolId {
			sandboxTool = tool
			return
		}
	}

	return
}

// UpdateSandboxTool applies an in-place update to an existing sandbox tool.
func (me *AgsService) UpdateSandboxTool(ctx context.Context, request *agsv20250920.UpdateSandboxToolRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.UpdateSandboxToolResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().UpdateSandboxToolWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("UpdateSandboxTool returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DeleteSandboxTool deletes the sandbox tool itself, identified by toolId.
func (me *AgsService) DeleteSandboxTool(ctx context.Context, toolId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDeleteSandboxToolRequest()
	request.ToolId = helper.String(toolId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DeleteSandboxToolResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DeleteSandboxToolWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DeleteSandboxTool returned an empty response for sandbox tool %s", toolId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// StartSandboxInstance starts a sandbox instance and returns the created
// instance. The instance id carried by the response is the only resource
// identity: RequestId is a request trace id and must never be used as identity.
func (me *AgsService) StartSandboxInstance(ctx context.Context, request *agsv20250920.StartSandboxInstanceRequest) (instance *agsv20250920.SandboxInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.StartSandboxInstanceResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().StartSandboxInstanceWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("StartSandboxInstance returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Instance == nil || response.Response.Instance.InstanceId == nil || *response.Response.Instance.InstanceId == "" {
		errRet = fmt.Errorf("StartSandboxInstance returned an empty instance id")
		return
	}

	instance = response.Response.Instance
	return
}

// DescribeSandboxInstanceById returns the sandbox instance with the given id.
//
// A nil result with a nil error means "queried successfully, no such instance";
// an API failure is always reported as an error so the caller never confuses a
// failed request with a deleted resource.
func (me *AgsService) DescribeSandboxInstanceById(ctx context.Context, instanceId string) (sandboxInstance *agsv20250920.SandboxInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDescribeSandboxInstancesRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DescribeSandboxInstancesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DescribeSandboxInstancesWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeSandboxInstances returned an empty response for sandbox instance %s", instanceId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, instance := range response.Response.InstanceSet {
		if instance != nil && instance.InstanceId != nil && *instance.InstanceId == instanceId {
			sandboxInstance = instance
			return
		}
	}

	return
}

// UpdateSandboxInstance applies an in-place update to an existing sandbox
// instance.
func (me *AgsService) UpdateSandboxInstance(ctx context.Context, request *agsv20250920.UpdateSandboxInstanceRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.UpdateSandboxInstanceResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().UpdateSandboxInstanceWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("UpdateSandboxInstance returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// StopSandboxInstance stops the sandbox instance itself, identified by
// instanceId.
func (me *AgsService) StopSandboxInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewStopSandboxInstanceRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.StopSandboxInstanceResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().StopSandboxInstanceWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("StopSandboxInstance returned an empty response for sandbox instance %s", instanceId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
