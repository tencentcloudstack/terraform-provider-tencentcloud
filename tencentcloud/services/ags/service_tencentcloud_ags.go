package ags

import (
	"context"
	"fmt"
	"log"

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
	AGS_SANDBOX_TOOL_STATUS_ACTIVE  = "ACTIVE"
	AGS_SANDBOX_TOOL_STATUS_FAILED  = "FAILED"
	AGS_SANDBOX_TOOL_STATUS_CREATING = "CREATING"
	AGS_SANDBOX_TOOL_STATUS_DELETING = "DELETING"
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

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAgsClient().CreateSandboxToolWithContext(ctx, request)
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

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAgsClient().DescribeSandboxToolListWithContext(ctx, request)
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

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAgsClient().UpdateSandboxToolWithContext(ctx, request)
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

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAgsClient().DeleteSandboxToolWithContext(ctx, request)
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
