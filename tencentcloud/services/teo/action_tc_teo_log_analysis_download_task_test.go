package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringTeoLogAnalysisDownloadTask(s string) *string {
	return &s
}

func teoLogAnalysisDownloadTaskSchema(t *testing.T) action_schema.Schema {
	a := teo.NewTeoLogAnalysisDownloadTask()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

func newTeoLogAnalysisDownloadTaskInvokeRequest(
	t *testing.T,
	zoneId, area, startTime, endTime, logType, condition, format, sort string,
) action.InvokeRequest {
	s := teoLogAnalysisDownloadTaskSchema(t)

	optional := func(v string) tftypes.Value {
		if v == "" {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":    tftypes.String,
			"area":       tftypes.String,
			"start_time": tftypes.String,
			"end_time":   tftypes.String,
			"log_type":   tftypes.String,
			"condition":  tftypes.String,
			"format":     tftypes.String,
			"sort":       tftypes.String,
		},
	}, map[string]tftypes.Value{
		"zone_id":    tftypes.NewValue(tftypes.String, zoneId),
		"area":       tftypes.NewValue(tftypes.String, area),
		"start_time": tftypes.NewValue(tftypes.String, startTime),
		"end_time":   tftypes.NewValue(tftypes.String, endTime),
		"log_type":   optional(logType),
		"condition":  optional(condition),
		"format":     optional(format),
		"sort":       optional(sort),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func setTeoLogAnalysisDownloadTaskClient(a *teo.TeoLogAnalysisDownloadTask, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoLogAnalysisDownloadTask" -v -count=1 -gcflags="all=-l"

func TestTeoLogAnalysisDownloadTask_MetadataAndSchema(t *testing.T) {
	a := teo.NewTeoLogAnalysisDownloadTask()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_log_analysis_download_task", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	for _, name := range []string{"zone_id", "area", "start_time", "end_time"} {
		assert.Contains(t, attrs, name)
		assert.True(t, attrs[name].IsRequired(), "%s should be required", name)
	}
	for _, name := range []string{"log_type", "condition", "format", "sort"} {
		assert.Contains(t, attrs, name)
		assert.True(t, attrs[name].IsOptional(), "%s should be optional", name)
	}
}

func TestTeoLogAnalysisDownloadTask_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotCreateRequest *teov20220901.CreateLogAnalysisDownloadTaskRequest
	patches.ApplyMethodFunc(teoClient, "CreateLogAnalysisDownloadTaskWithContext", func(_ context.Context, request *teov20220901.CreateLogAnalysisDownloadTaskRequest) (*teov20220901.CreateLogAnalysisDownloadTaskResponse, error) {
		gotCreateRequest = request
		resp := teov20220901.NewCreateLogAnalysisDownloadTaskResponse()
		resp.Response = &teov20220901.CreateLogAnalysisDownloadTaskResponseParams{
			TaskId:    ptrStringTeoLogAnalysisDownloadTask("task-3fkff38fyw8s"),
			RequestId: ptrStringTeoLogAnalysisDownloadTask("fake-request-id"),
		}
		return resp, nil
	})

	var gotDescribeRequest *teov20220901.DescribeLogAnalysisDownloadTasksRequest
	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasksWithContext", func(_ context.Context, request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		gotDescribeRequest = request
		resp := teov20220901.NewDescribeLogAnalysisDownloadTasksResponse()
		resp.Response = &teov20220901.DescribeLogAnalysisDownloadTasksResponseParams{
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			Tasks: []*teov20220901.LogAnalysisDownloadTask{
				{
					TaskId:     ptrStringTeoLogAnalysisDownloadTask("task-3fkff38fyw8s"),
					ZoneId:     ptrStringTeoLogAnalysisDownloadTask("zone-3fkff38fyw8s"),
					Area:       ptrStringTeoLogAnalysisDownloadTask("mainland"),
					Status:     ptrStringTeoLogAnalysisDownloadTask("completed"),
					Url:        ptrStringTeoLogAnalysisDownloadTask("https://example.com/log.csv"),
					CreateTime: ptrStringTeoLogAnalysisDownloadTask("2020-04-30T00:10:00Z"),
				},
			},
			RequestId: ptrStringTeoLogAnalysisDownloadTask("fake-request-id"),
		}
		return resp, nil
	})

	a := teo.NewTeoLogAnalysisDownloadTask()
	setTeoLogAnalysisDownloadTaskClient(a.(*teo.TeoLogAnalysisDownloadTask), client)

	req := newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"zone-3fkff38fyw8s", "mainland", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"l7-access-logs", "", "csv", "desc")

	var messages []string
	resp := &action.InvokeResponse{}
	resp.SendProgress = func(event action.InvokeProgressEvent) {
		messages = append(messages, event.Message)
	}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotCreateRequest)
	assert.Equal(t, "zone-3fkff38fyw8s", *gotCreateRequest.ZoneId)
	assert.Equal(t, "mainland", *gotCreateRequest.Area)
	assert.Equal(t, "2020-04-29T00:00:00Z", *gotCreateRequest.StartTime)
	assert.Equal(t, "2020-04-30T00:00:00Z", *gotCreateRequest.EndTime)
	assert.Equal(t, "l7-access-logs", *gotCreateRequest.LogType)
	assert.Nil(t, gotCreateRequest.Condition)
	assert.Equal(t, "csv", *gotCreateRequest.Format)
	assert.Equal(t, "desc", *gotCreateRequest.Sort)

	assert.NotNil(t, gotDescribeRequest)
	assert.Equal(t, "zone-3fkff38fyw8s", *gotDescribeRequest.ZoneId)
	assert.Equal(t, "mainland", *gotDescribeRequest.Area)
	assert.Equal(t, 1, len(gotDescribeRequest.Filters))
	assert.Equal(t, "task-id", *gotDescribeRequest.Filters[0].Name)
	assert.Equal(t, "task-3fkff38fyw8s", *gotDescribeRequest.Filters[0].Values[0])

	assert.Equal(t, 3, len(messages))
	assert.Contains(t, messages[0], "Creating TEO log analysis download task")
	assert.Contains(t, messages[1], "task-3fkff38fyw8s")
	assert.Contains(t, messages[2], "https://example.com/log.csv")
}

func TestTeoLogAnalysisDownloadTask_Invoke_TaskFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateLogAnalysisDownloadTaskWithContext", func(_ context.Context, request *teov20220901.CreateLogAnalysisDownloadTaskRequest) (*teov20220901.CreateLogAnalysisDownloadTaskResponse, error) {
		resp := teov20220901.NewCreateLogAnalysisDownloadTaskResponse()
		resp.Response = &teov20220901.CreateLogAnalysisDownloadTaskResponseParams{
			TaskId: ptrStringTeoLogAnalysisDownloadTask("task-3fkff38fyw8s"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasksWithContext", func(_ context.Context, request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		resp := teov20220901.NewDescribeLogAnalysisDownloadTasksResponse()
		resp.Response = &teov20220901.DescribeLogAnalysisDownloadTasksResponseParams{
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			Tasks: []*teov20220901.LogAnalysisDownloadTask{
				{
					TaskId: ptrStringTeoLogAnalysisDownloadTask("task-3fkff38fyw8s"),
					Status: ptrStringTeoLogAnalysisDownloadTask("failed"),
				},
			},
		}
		return resp, nil
	})

	a := teo.NewTeoLogAnalysisDownloadTask()
	setTeoLogAnalysisDownloadTaskClient(a.(*teo.TeoLogAnalysisDownloadTask), client)

	req := newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"zone-3fkff38fyw8s", "mainland", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"", "", "", "")

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "task-3fkff38fyw8s")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "failed")
}

func TestTeoLogAnalysisDownloadTask_Invoke_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateLogAnalysisDownloadTaskWithContext", func(_ context.Context, request *teov20220901.CreateLogAnalysisDownloadTaskRequest) (*teov20220901.CreateLogAnalysisDownloadTaskResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter.ZoneNotFound, Message=zone not found")
	})

	describeCalled := false
	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasksWithContext", func(_ context.Context, request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		describeCalled = true
		return teov20220901.NewDescribeLogAnalysisDownloadTasksResponse(), nil
	})

	a := teo.NewTeoLogAnalysisDownloadTask()
	setTeoLogAnalysisDownloadTaskClient(a.(*teo.TeoLogAnalysisDownloadTask), client)

	req := newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"zone-3fkff38fyw8s", "mainland", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"", "", "", "")

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "creating TEO log analysis download task")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InvalidParameter.ZoneNotFound")
	assert.False(t, describeCalled)
}

func TestTeoLogAnalysisDownloadTask_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	apiCalled := false
	patches.ApplyMethodFunc(teoClient, "CreateLogAnalysisDownloadTaskWithContext", func(_ context.Context, request *teov20220901.CreateLogAnalysisDownloadTaskRequest) (*teov20220901.CreateLogAnalysisDownloadTaskResponse, error) {
		apiCalled = true
		return teov20220901.NewCreateLogAnalysisDownloadTaskResponse(), nil
	})

	a := teo.NewTeoLogAnalysisDownloadTask()
	setTeoLogAnalysisDownloadTaskClient(a.(*teo.TeoLogAnalysisDownloadTask), client)

	req := newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"", "mainland", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"", "", "", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing zone_id")
	assert.False(t, apiCalled)

	req = newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"zone-3fkff38fyw8s", "invalid-area", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"", "", "", "")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid area")
	assert.False(t, apiCalled)

	req = newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"zone-3fkff38fyw8s", "mainland", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"", "", "", "unknown")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid sort")
	assert.False(t, apiCalled)
}

func TestTeoLogAnalysisDownloadTask_Invoke_ClientNotConfigured(t *testing.T) {
	a := teo.NewTeoLogAnalysisDownloadTask()

	req := newTeoLogAnalysisDownloadTaskInvokeRequest(t,
		"zone-3fkff38fyw8s", "mainland", "2020-04-29T00:00:00Z", "2020-04-30T00:00:00Z",
		"", "", "", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}
