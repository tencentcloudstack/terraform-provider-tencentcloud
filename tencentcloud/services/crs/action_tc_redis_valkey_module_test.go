package crs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	redisv20180412 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
)

func ptrStringRedisValkeyModule(s string) *string {
	return &s
}

// redisValkeyModuleSchema returns the framework action schema used to build
// the raw config value for InvokeRequest.
func redisValkeyModuleSchema(t *testing.T) action_schema.Schema {
	a := crs.NewRedisValkeyModule()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newRedisValkeyModuleInvokeRequest builds an action.InvokeRequest carrying the
// given instance_id, name and status as raw tftypes values.
func newRedisValkeyModuleInvokeRequest(t *testing.T, instanceId, name, status string) action.InvokeRequest {
	s := redisValkeyModuleSchema(t)

	attrTypes := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"instance_id": tftypes.String,
			"name":        tftypes.String,
			"status":      tftypes.String,
		},
	}

	raw := tftypes.NewValue(attrTypes, map[string]tftypes.Value{
		"instance_id": tftypes.NewValue(tftypes.String, instanceId),
		"name":        tftypes.NewValue(tftypes.String, name),
		"status":      tftypes.NewValue(tftypes.String, status),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newRedisValkeyModuleInvokeRequestNullField builds an InvokeRequest whose
// given field is null (all other fields filled with valid values).
func newRedisValkeyModuleInvokeRequestNullField(t *testing.T, nullField string) action.InvokeRequest {
	s := redisValkeyModuleSchema(t)

	values := map[string]tftypes.Value{
		"instance_id": tftypes.NewValue(tftypes.String, "crs-1111222233334444"),
		"name":        tftypes.NewValue(tftypes.String, "ValkeyBloom"),
		"status":      tftypes.NewValue(tftypes.String, "ON"),
	}
	values[nullField] = tftypes.NewValue(tftypes.String, nil)

	attrTypes := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"instance_id": tftypes.String,
			"name":        tftypes.String,
			"status":      tftypes.String,
		},
	}

	req := action.InvokeRequest{}
	req.Config.Raw = tftypes.NewValue(attrTypes, values)
	req.Config.Schema = s
	return req
}

// setRedisValkeyModuleClient injects a mock connectivity client into the action
// instance via the Configure method.
func setRedisValkeyModuleClient(a *crs.RedisValkeyModule, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/crs/ -run "TestRedisValkeyModule" -v -count=1 -gcflags="all=-l"

// TestRedisValkeyModule_Invoke_Success tests a successful invoke where the
// asynchronous task is already in succeed status on the first poll.
func TestRedisValkeyModule_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	redisClient := &redisv20180412.Client{}
	patches.ApplyMethodReturn(client, "UseRedisClient", redisClient)

	var gotModifyRequest *redisv20180412.ModifyValkeyModuleRequest
	patches.ApplyMethodFunc(redisClient, "ModifyValkeyModuleWithContext", func(_ context.Context, request *redisv20180412.ModifyValkeyModuleRequest) (*redisv20180412.ModifyValkeyModuleResponse, error) {
		gotModifyRequest = request
		resp := &redisv20180412.ModifyValkeyModuleResponse{}
		resp.Response = &redisv20180412.ModifyValkeyModuleResponseParams{
			TaskId:    func() *int64 { v := int64(12345); return &v }(),
			RequestId: ptrStringRedisValkeyModule("fake-request-id"),
		}
		return resp, nil
	})

	var gotTaskRequest *redisv20180412.DescribeTaskInfoRequest
	patches.ApplyMethodFunc(redisClient, "DescribeTaskInfo", func(request *redisv20180412.DescribeTaskInfoRequest) (*redisv20180412.DescribeTaskInfoResponse, error) {
		gotTaskRequest = request
		resp := &redisv20180412.DescribeTaskInfoResponse{}
		resp.Response = &redisv20180412.DescribeTaskInfoResponseParams{
			Status: ptrStringRedisValkeyModule("succeed"),
		}
		return resp, nil
	})

	a := crs.NewRedisValkeyModule()
	setRedisValkeyModuleClient(a.(*crs.RedisValkeyModule), client)

	req := newRedisValkeyModuleInvokeRequest(t, "crs-1111222233334444", "ValkeyBloom", "ON")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotModifyRequest)
	assert.Equal(t, "crs-1111222233334444", *gotModifyRequest.InstanceId)
	assert.Equal(t, "ValkeyBloom", *gotModifyRequest.Name)
	assert.Equal(t, "ON", *gotModifyRequest.Status)
	assert.NotNil(t, gotTaskRequest)
	assert.Equal(t, uint64(12345), *gotTaskRequest.TaskId)
}

// TestRedisValkeyModule_Invoke_TaskFailed tests the action fails when the
// asynchronous task ends with a failed status. The polling budget is shrunk to
// milliseconds so the retry loop does not really wait for the full timeout.
func TestRedisValkeyModule_Invoke_TaskFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	redisClient := &redisv20180412.Client{}
	patches.ApplyMethodReturn(client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyValkeyModuleWithContext", func(_ context.Context, request *redisv20180412.ModifyValkeyModuleRequest) (*redisv20180412.ModifyValkeyModuleResponse, error) {
		resp := &redisv20180412.ModifyValkeyModuleResponse{}
		resp.Response = &redisv20180412.ModifyValkeyModuleResponseParams{
			TaskId: func() *int64 { v := int64(67890); return &v }(),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(redisClient, "DescribeTaskInfo", func(request *redisv20180412.DescribeTaskInfoRequest) (*redisv20180412.DescribeTaskInfoResponse, error) {
		resp := &redisv20180412.DescribeTaskInfoResponse{}
		resp.Response = &redisv20180412.DescribeTaskInfoResponseParams{
			Status: ptrStringRedisValkeyModule("failed"),
		}
		return resp, nil
	})

	patches.ApplyGlobalVar(&tccommon.WriteRetryTimeout, 1)

	a := crs.NewRedisValkeyModule()
	setRedisValkeyModuleClient(a.(*crs.RedisValkeyModule), client)

	req := newRedisValkeyModuleInvokeRequest(t, "crs-1111222233334444", "ValkeyJSON", "ON")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for Valkey module")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "redis task exe fail")
}

// TestRedisValkeyModule_Invoke_ModifyAPIError tests API error handling of the
// ModifyValkeyModule call.
func TestRedisValkeyModule_Invoke_ModifyAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	redisClient := &redisv20180412.Client{}
	patches.ApplyMethodReturn(client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyValkeyModuleWithContext", func(_ context.Context, request *redisv20180412.ModifyValkeyModuleRequest) (*redisv20180412.ModifyValkeyModuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.UnsupportedModule, Message=module not supported")
	})

	a := crs.NewRedisValkeyModule()
	setRedisValkeyModuleClient(a.(*crs.RedisValkeyModule), client)

	req := newRedisValkeyModuleInvokeRequest(t, "crs-1111222233334444", "ValkeySearch", "ON")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "modifying Valkey module")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "UnsupportedModule")
}

// TestRedisValkeyModule_Invoke_MissingInput tests missing required inputs are
// rejected before any API call.
func TestRedisValkeyModule_Invoke_MissingInput(t *testing.T) {
	client := &connectivity.TencentCloudClient{}
	redisClient := &redisv20180412.Client{}
	for _, field := range []string{"instance_id", "name", "status"} {
		patches := gomonkey.NewPatches()
		patches.ApplyMethodReturn(client, "UseRedisClient", redisClient)
		a := crs.NewRedisValkeyModule()
		setRedisValkeyModuleClient(a.(*crs.RedisValkeyModule), client)

		req := newRedisValkeyModuleInvokeRequestNullField(t, field)
		resp := &action.InvokeResponse{}
		a.Invoke(context.Background(), req, resp)

		assert.True(t, resp.Diagnostics.HasError(), "field %s: expected error diagnostics", field)
		patches.Reset()
	}
}

// TestRedisValkeyModule_Invoke_EmptyModifyResponse tests the empty response
// branch of the ModifyValkeyModule call.
func TestRedisValkeyModule_Invoke_EmptyModifyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	redisClient := &redisv20180412.Client{}
	patches.ApplyMethodReturn(client, "UseRedisClient", redisClient)

	patches.ApplyMethodFunc(redisClient, "ModifyValkeyModuleWithContext", func(_ context.Context, request *redisv20180412.ModifyValkeyModuleRequest) (*redisv20180412.ModifyValkeyModuleResponse, error) {
		return nil, nil
	})

	a := crs.NewRedisValkeyModule()
	setRedisValkeyModuleClient(a.(*crs.RedisValkeyModule), client)

	req := newRedisValkeyModuleInvokeRequest(t, "crs-1111222233334444", "ValkeyBloom", "ON")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "empty response")
}
