package bdrc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"
)

func ptrStringBdrcRunCopyPairTasks(s string) *string {
	return &s
}

// bdrcRunCopyPairTasksSchema returns the framework action schema used to build
// the raw config value for InvokeRequest.
func bdrcRunCopyPairTasksSchema(t *testing.T) action_schema.Schema {
	a := bdrc.NewBdrcRunCopyPairTasks()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newBdrcRunCopyPairTasksInvokeRequest builds an action.InvokeRequest carrying
// the given copy_pair_ids and copy_pair_type as raw tftypes values.
func newBdrcRunCopyPairTasksInvokeRequest(t *testing.T, copyPairIds []string, copyPairType string) action.InvokeRequest {
	s := bdrcRunCopyPairTasksSchema(t)

	idVals := make([]tftypes.Value, 0, len(copyPairIds))
	for _, id := range copyPairIds {
		idVals = append(idVals, tftypes.NewValue(tftypes.String, id))
	}

	var idsVal tftypes.Value
	if len(idVals) == 0 {
		idsVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		idsVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, idVals)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"copy_pair_ids":  tftypes.List{ElementType: tftypes.String},
			"copy_pair_type": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"copy_pair_ids":  idsVal,
		"copy_pair_type": tftypes.NewValue(tftypes.String, copyPairType),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newBdrcRunCopyPairTasksInvokeRequestNullIds builds an InvokeRequest whose
// copy_pair_ids is null.
func newBdrcRunCopyPairTasksInvokeRequestNullIds(t *testing.T, copyPairType string) action.InvokeRequest {
	s := bdrcRunCopyPairTasksSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"copy_pair_ids":  tftypes.List{ElementType: tftypes.String},
			"copy_pair_type": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"copy_pair_ids":  tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
		"copy_pair_type": tftypes.NewValue(tftypes.String, copyPairType),
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setBdrcRunCopyPairTasksClient injects a mock connectivity client into the
// action instance via the Configure method.
func setBdrcRunCopyPairTasksClient(a *bdrc.BdrcRunCopyPairTasks, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/bdrc/ -run "TestBdrcRunCopyPairTasks" -v -count=1 -gcflags="all=-l"

// TestBdrcRunCopyPairTasks_Invoke_Success tests a successful invoke
func TestBdrcRunCopyPairTasks_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(client, "UseBdrcV20260330Client", bdrcClient)

	var gotRequest *bdrcv20260330.RunCopyPairTasksRequest
	patches.ApplyMethodFunc(bdrcClient, "RunCopyPairTasksWithContext", func(_ context.Context, request *bdrcv20260330.RunCopyPairTasksRequest) (*bdrcv20260330.RunCopyPairTasksResponse, error) {
		gotRequest = request
		resp := bdrcv20260330.NewRunCopyPairTasksResponse()
		resp.Response = &bdrcv20260330.RunCopyPairTasksResponseParams{
			CopyPairIds: []*string{
				ptrStringBdrcRunCopyPairTasks("pair-1111222233334444"),
				ptrStringBdrcRunCopyPairTasks("pair-5555666677778888"),
			},
			RequestId: ptrStringBdrcRunCopyPairTasks("fake-request-id"),
		}
		return resp, nil
	})

	a := bdrc.NewBdrcRunCopyPairTasks()
	setBdrcRunCopyPairTasksClient(a.(*bdrc.BdrcRunCopyPairTasks), client)

	req := newBdrcRunCopyPairTasksInvokeRequest(t,
		[]string{"pair-1111222233334444", "pair-5555666677778888"},
		"DISK")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, []*string{
		ptrStringBdrcRunCopyPairTasks("pair-1111222233334444"),
		ptrStringBdrcRunCopyPairTasks("pair-5555666677778888"),
	}, gotRequest.CopyPairIds)
	assert.Equal(t, "DISK", *gotRequest.CopyPairType)
}

// TestBdrcRunCopyPairTasks_Invoke_APIError tests API error handling
func TestBdrcRunCopyPairTasks_Invoke_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "RunCopyPairTasksWithContext", func(_ context.Context, request *bdrcv20260330.RunCopyPairTasksRequest) (*bdrcv20260330.RunCopyPairTasksResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.CopyPairNotExist, Message=copy pair not exist")
	})

	a := bdrc.NewBdrcRunCopyPairTasks()
	setBdrcRunCopyPairTasksClient(a.(*bdrc.BdrcRunCopyPairTasks), client)

	req := newBdrcRunCopyPairTasksInvokeRequest(t,
		[]string{"pair-invalid"},
		"DISK")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "launching BDRC copy pair tasks")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ResourceNotFound.CopyPairNotExist")
}

// TestBdrcRunCopyPairTasks_Invoke_MissingInput tests missing required inputs are
// rejected before any API call
func TestBdrcRunCopyPairTasks_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(client, "UseBdrcV20260330Client", bdrcClient)

	apiCalled := false
	patches.ApplyMethodFunc(bdrcClient, "RunCopyPairTasksWithContext", func(_ context.Context, request *bdrcv20260330.RunCopyPairTasksRequest) (*bdrcv20260330.RunCopyPairTasksResponse, error) {
		apiCalled = true
		return bdrcv20260330.NewRunCopyPairTasksResponse(), nil
	})

	a := bdrc.NewBdrcRunCopyPairTasks()
	setBdrcRunCopyPairTasksClient(a.(*bdrc.BdrcRunCopyPairTasks), client)

	// empty copy_pair_type
	req := newBdrcRunCopyPairTasksInvokeRequest(t, []string{"pair-1111222233334444"}, "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing copy_pair_type")
	assert.False(t, apiCalled)

	// null copy_pair_ids
	req = newBdrcRunCopyPairTasksInvokeRequestNullIds(t, "DISK")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing copy_pair_ids")
	assert.False(t, apiCalled)
}

// TestBdrcRunCopyPairTasks_Invoke_ClientNotConfigured tests the nil-client guard
func TestBdrcRunCopyPairTasks_Invoke_ClientNotConfigured(t *testing.T) {
	a := bdrc.NewBdrcRunCopyPairTasks()

	req := newBdrcRunCopyPairTasksInvokeRequest(t, []string{"pair-1111222233334444"}, "DISK")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestBdrcRunCopyPairTasks_MetadataAndSchema validates metadata and schema definition
func TestBdrcRunCopyPairTasks_MetadataAndSchema(t *testing.T) {
	a := bdrc.NewBdrcRunCopyPairTasks()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_bdrc_run_copy_pair_tasks", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "copy_pair_ids")
	assert.Contains(t, attrs, "copy_pair_type")

	assert.True(t, attrs["copy_pair_ids"].IsRequired())
	assert.True(t, attrs["copy_pair_type"].IsRequired())

	idsAttr, ok := attrs["copy_pair_ids"].(action_schema.ListAttribute)
	assert.True(t, ok)
	assert.Equal(t, types.StringType, idsAttr.ElementType)

	typeAttr, ok := attrs["copy_pair_type"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, typeAttr.Required)
}
