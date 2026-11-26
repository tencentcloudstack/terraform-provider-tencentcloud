package cngw_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	basetypes "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	cngwv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cngw/v20230418"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcngw "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cngw"
)

const (
	cngwAIServiceSourceTestId          = "source-12345"
	cngwAIServiceSourceTestGatewayId   = "gw-12345"
	cngwAIServiceSourceTestSourceType  = "MCPRegistry"
	cngwAIServiceSourceTestSourceName  = "tf-example-source"
	cngwAIServiceSourceTestNewName     = "tf-example-source-renamed"
	cngwAIServiceSourceTestNormalName  = "tf-example-source-normalized"
	cngwAIServiceSourceTestDescription = "created by terraform"
	cngwAIServiceSourceTestInstanceId  = "ins-12345"
	cngwAIServiceSourceTestAddressA    = "10.0.0.1:8848"
	cngwAIServiceSourceTestAddressB    = "10.0.0.2:8848"
	cngwAIServiceSourceTestNewAddress  = "10.0.0.9:8848"
)

// cngwAIServiceSourceTestStr is a pointer helper for fake string values.
func cngwAIServiceSourceTestStr(v string) *string { return &v }

// cngwAIServiceSourceTestSourceInfoAttrTypes mirrors the attribute types of
// the source_info object attribute declared by the resource schema.
func cngwAIServiceSourceTestSourceInfoAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instance_id": types.StringType,
		"auth": types.ObjectType{AttrTypes: map[string]attr.Type{
			"username":     types.StringType,
			"password":     types.StringType,
			"access_token": types.StringType,
		}},
		"addresses": types.ListType{ElemType: types.StringType},
	}
}

// cngwAIServiceSourceTestSourceInfo builds a source_info object attribute value.
func cngwAIServiceSourceTestSourceInfo(instanceId string, addresses ...string) types.Object {
	addrValues := make([]attr.Value, 0, len(addresses))
	for _, a := range addresses {
		addrValues = append(addrValues, types.StringValue(a))
	}

	elements := types.ListValueMust(types.StringType, addrValues)

	authObject := types.ObjectValueMust(
		map[string]attr.Type{
			"username":     types.StringType,
			"password":     types.StringType,
			"access_token": types.StringType,
		},
		map[string]attr.Value{
			"username":     types.StringNull(),
			"password":     types.StringNull(),
			"access_token": types.StringNull(),
		},
	)

	return types.ObjectValueMust(
		cngwAIServiceSourceTestSourceInfoAttrTypes(),
		map[string]attr.Value{
			"instance_id": types.StringValue(instanceId),
			"auth":        authObject,
			"addresses":   elements,
		},
	)
}

// cngwAIServiceSourceTestItem builds the describe response item.
func cngwAIServiceSourceTestItem() *cngwv20230418.CNAPIGwAIServiceSource {
	return &cngwv20230418.CNAPIGwAIServiceSource{
		SourceName:  cngwAIServiceSourceTestStr(cngwAIServiceSourceTestSourceName),
		SourceId:    cngwAIServiceSourceTestStr(cngwAIServiceSourceTestId),
		SourceType:  cngwAIServiceSourceTestStr(cngwAIServiceSourceTestSourceType),
		SourceInfo:  cngwAIServiceSourceTestSourceInfoSdk(),
		Description: cngwAIServiceSourceTestStr(cngwAIServiceSourceTestDescription),
		CreateTime:  cngwAIServiceSourceTestStr("2024-01-01 00:00:00"),
		UpdateTime:  cngwAIServiceSourceTestStr("2024-01-01 00:00:00"),
	}
}

// cngwAIServiceSourceTestSourceInfoSdk builds the SDK-side source info that
// mirrors cngwAIServiceSourceTestSourceInfo.
func cngwAIServiceSourceTestSourceInfoSdk() *cngwv20230418.CNAPIGwAIServiceSourceInfo {
	return &cngwv20230418.CNAPIGwAIServiceSourceInfo{
		InstanceId: cngwAIServiceSourceTestStr(cngwAIServiceSourceTestInstanceId),
		Auth:       &cngwv20230418.CNAPIGwAIServiceSourceAuth{},
		Addresses: []*string{
			cngwAIServiceSourceTestStr(cngwAIServiceSourceTestAddressA),
			cngwAIServiceSourceTestStr(cngwAIServiceSourceTestAddressB),
		},
	}
}

// cngwAIServiceSourceTestModel builds the terraform model used for plan/state;
// every unset attribute carries an EXPLICIT typed null value (a raw zero value
// breaks tfsdk Plan/State reflection).
func cngwAIServiceSourceTestModel() svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel {
	return svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel{
		Id:            types.StringValue(cngwAIServiceSourceTestId),
		GatewayId:     types.StringValue(cngwAIServiceSourceTestGatewayId),
		SourceType:    types.StringValue(cngwAIServiceSourceTestSourceType),
		SourceName:    types.StringValue(cngwAIServiceSourceTestSourceName),
		SourceProduct: types.StringNull(),
		SourceInfo:    cngwAIServiceSourceTestSourceInfo(cngwAIServiceSourceTestInstanceId, cngwAIServiceSourceTestAddressA, cngwAIServiceSourceTestAddressB),
		Description:   types.StringValue(cngwAIServiceSourceTestDescription),
	}
}

// assertCngwAIServiceSourceModel asserts every model field.
func assertCngwAIServiceSourceModel(t *testing.T, m svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel) {
	t.Helper()

	if m.Id.ValueString() != cngwAIServiceSourceTestId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), cngwAIServiceSourceTestId)
	}
	if m.GatewayId.ValueString() != cngwAIServiceSourceTestGatewayId {
		t.Errorf("gateway_id = %q, want %q", m.GatewayId.ValueString(), cngwAIServiceSourceTestGatewayId)
	}
	if m.SourceType.ValueString() != cngwAIServiceSourceTestSourceType {
		t.Errorf("source_type = %q, want %q", m.SourceType.ValueString(), cngwAIServiceSourceTestSourceType)
	}
	if m.SourceName.ValueString() != cngwAIServiceSourceTestSourceName {
		t.Errorf("source_name = %q, want %q", m.SourceName.ValueString(), cngwAIServiceSourceTestSourceName)
	}
	if m.SourceInfo.IsNull() {
		t.Fatal("source_info should not be null")
	}
	var sourceInfo struct {
		InstanceId types.String `tfsdk:"instance_id"`
		Auth       types.Object `tfsdk:"auth"`
		Addresses  types.List   `tfsdk:"addresses"`
	}
	if diags := m.SourceInfo.As(context.Background(), &sourceInfo, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("source_info: %v", diags.Errors())
	}
	if sourceInfo.InstanceId.ValueString() != cngwAIServiceSourceTestInstanceId {
		t.Errorf("source_info.instance_id = %q, want %q", sourceInfo.InstanceId.ValueString(), cngwAIServiceSourceTestInstanceId)
	}
	if sourceInfo.Auth.IsNull() {
		t.Fatal("source_info.auth should not be null")
	}
	var addresses []string
	if diags := sourceInfo.Addresses.ElementsAs(context.Background(), &addresses, false); diags.HasError() {
		t.Fatalf("addresses: %v", diags.Errors())
	}
	if len(addresses) != 2 || addresses[0] != cngwAIServiceSourceTestAddressA || addresses[1] != cngwAIServiceSourceTestAddressB {
		t.Errorf("addresses = %v, want [%q %q]", addresses, cngwAIServiceSourceTestAddressA, cngwAIServiceSourceTestAddressB)
	}
	if m.Description.ValueString() != cngwAIServiceSourceTestDescription {
		t.Errorf("description = %q, want %q", m.Description.ValueString(), cngwAIServiceSourceTestDescription)
	}
}

// patchCngwAIServiceSourceSdkClient stubs the versioned client accessor and
// returns the mock SDK client whose methods the test patches next (the
// two-level mode).
func patchCngwAIServiceSourceSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *cngwv20230418.Client {
	cngwClient := &cngwv20230418.Client{}
	patches.ApplyMethodReturn(client, "UseCngwV20230418Client", cngwClient)

	return cngwClient
}

// stubCngwAIServiceSourceDescribe answers the list describe API with the
// given items; a captured filter validates the SourceId filter propagation.
func stubCngwAIServiceSourceDescribe(patches *gomonkey.Patches, cngwClient *cngwv20230418.Client, items []*cngwv20230418.CNAPIGwAIServiceSource, captured **cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListRequest) {
	patches.ApplyMethodFunc(cngwClient, "DescribeCloudNativeAPIGatewayAIServiceSourceList", func(req *cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListRequest) (*cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListResponse, error) {
		if captured != nil {
			*captured = req
		}

		var dataList []*cngwv20230418.CNAPIGwAIServiceSource
		if req.Filters != nil && len(req.Filters) > 0 {
			wantId := ""
			if len(req.Filters[0].Values) > 0 && req.Filters[0].Values[0] != nil {
				wantId = *req.Filters[0].Values[0]
			}
			for _, item := range items {
				if item.SourceId != nil && *item.SourceId == wantId {
					dataList = append(dataList, item)
				}
			}
		} else {
			dataList = items
		}

		return &cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListResponse{Response: &cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListResponseParams{
			Result: &cngwv20230418.CNAPIGwAIServiceSourceList{
				DataList:   dataList,
				TotalCount: func() *uint64 { v := uint64(len(dataList)); return &v }(),
			},
			RequestId: cngwAIServiceSourceTestStr("fake-request-id"),
		}}, nil
	})
}

// cngwAIServiceSourceTestContext wires the resource to a configured client;
// the SDK client methods are stubbed per test case.
type cngwAIServiceSourceTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newCngwAIServiceSourceTestContext(t *testing.T) *cngwAIServiceSourceTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svcngw.NewCngwCloudNativeAPIGatewayAIServiceSourceResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &cngwAIServiceSourceTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *cngwAIServiceSourceTestContext) plan(t *testing.T, m svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *cngwAIServiceSourceTestContext) state(t *testing.T, m svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// TestCngwAIServiceSourceCreateSuccess verifies that Create sends all request
// fields, tracks the returned source id and refreshes the state.
func TestCngwAIServiceSourceCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)

	var captured *cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceRequest
	patches.ApplyMethodFunc(cngwClient, "CreateCloudNativeAPIGatewayAIServiceSource", func(req *cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceRequest) (*cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponse, error) {
		captured = req

		return &cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponse{Response: &cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponseParams{
			Result: &cngwv20230418.CNAPIGwCreateCommonResult{
				ID:      cngwAIServiceSourceTestStr(cngwAIServiceSourceTestId),
				Success: func() *bool { v := true; return &v }(),
			},
			RequestId: cngwAIServiceSourceTestStr("fake-request-id"),
		}}, nil
	})
	stubCngwAIServiceSourceDescribe(patches, cngwClient, []*cngwv20230418.CNAPIGwAIServiceSource{cngwAIServiceSourceTestItem()}, nil)

	model := cngwAIServiceSourceTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateCloudNativeAPIGatewayAIServiceSource request is nil")
	}
	if captured.GatewayId == nil || *captured.GatewayId != cngwAIServiceSourceTestGatewayId {
		t.Errorf("create request gateway_id = %v, want %q", captured.GatewayId, cngwAIServiceSourceTestGatewayId)
	}
	if captured.SourceType == nil || *captured.SourceType != cngwAIServiceSourceTestSourceType {
		t.Errorf("create request source_type = %v, want %q", captured.SourceType, cngwAIServiceSourceTestSourceType)
	}
	if captured.SourceName == nil || *captured.SourceName != cngwAIServiceSourceTestSourceName {
		t.Errorf("create request source_name = %v, want %q", captured.SourceName, cngwAIServiceSourceTestSourceName)
	}
	if captured.Description == nil || *captured.Description != cngwAIServiceSourceTestDescription {
		t.Errorf("create request description = %v, want %q", captured.Description, cngwAIServiceSourceTestDescription)
	}
	if captured.SourceInfo == nil {
		t.Fatal("create request source_info is nil")
	}
	if captured.SourceInfo.InstanceId == nil || *captured.SourceInfo.InstanceId != cngwAIServiceSourceTestInstanceId {
		t.Errorf("create request source_info.instance_id = %v, want %q", captured.SourceInfo.InstanceId, cngwAIServiceSourceTestInstanceId)
	}
	if len(captured.SourceInfo.Addresses) != 2 || *captured.SourceInfo.Addresses[0] != cngwAIServiceSourceTestAddressA || *captured.SourceInfo.Addresses[1] != cngwAIServiceSourceTestAddressB {
		t.Errorf("create request addresses = %v, want [%q %q]", captured.SourceInfo.Addresses, cngwAIServiceSourceTestAddressA, cngwAIServiceSourceTestAddressB)
	}

	var got svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCngwAIServiceSourceModel(t, got)
}

// TestCngwAIServiceSourceCreateEmptyId verifies that Create fails loudly when
// the API returns no source id.
func TestCngwAIServiceSourceCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(cngwClient, "CreateCloudNativeAPIGatewayAIServiceSource", func(req *cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceRequest) (*cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponse, error) {
		return &cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponse{Response: &cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponseParams{
			RequestId: cngwAIServiceSourceTestStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, cngwAIServiceSourceTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty source id")
	}
}

// TestCngwAIServiceSourceCreateApiError verifies API failures surface as errors.
func TestCngwAIServiceSourceCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(cngwClient, "CreateCloudNativeAPIGatewayAIServiceSource", func(req *cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceRequest) (*cngwv20230418.CreateCloudNativeAPIGatewayAIServiceSourceResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, cngwAIServiceSourceTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestCngwAIServiceSourceReadSuccess verifies Read refreshes all leaf fields
// and that the read passes the tracked source id through the SourceId filter.
func TestCngwAIServiceSourceReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)

	var captured *cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListRequest
	stubCngwAIServiceSourceDescribe(patches, cngwClient, []*cngwv20230418.CNAPIGwAIServiceSource{cngwAIServiceSourceTestItem()}, &captured)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cngwAIServiceSourceTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DescribeCloudNativeAPIGatewayAIServiceSourceList request is nil")
	}
	if captured.Filters == nil || len(captured.Filters) != 1 || len(captured.Filters[0].Values) != 1 || *captured.Filters[0].Values[0] != cngwAIServiceSourceTestId {
		t.Errorf("read request SourceId filter = %v, want [%q]", captured.Filters, cngwAIServiceSourceTestId)
	}

	var got svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertCngwAIServiceSourceModel(t, got)
}

// TestCngwAIServiceSourceReadDetectsAddressesDrift verifies a user-configurable
// collection inside source_info is refreshed from the cloud: the state holds
// two addresses, the cloud now answers one different address, and Read must
// record the cloud value so a change made outside Terraform shows up in the
// next plan.
func TestCngwAIServiceSourceReadDetectsAddressesDrift(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)

	drifted := cngwAIServiceSourceTestItem()
	drifted.SourceInfo.Addresses = []*string{cngwAIServiceSourceTestStr(cngwAIServiceSourceTestNewAddress)}
	stubCngwAIServiceSourceDescribe(patches, cngwClient, []*cngwv20230418.CNAPIGwAIServiceSource{drifted}, nil)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cngwAIServiceSourceTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}

	var sourceInfo struct {
		InstanceId types.String `tfsdk:"instance_id"`
		Auth       types.Object `tfsdk:"auth"`
		Addresses  types.List   `tfsdk:"addresses"`
	}
	if diags := got.SourceInfo.As(context.Background(), &sourceInfo, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("source_info: %v", diags.Errors())
	}
	var addresses []string
	if diags := sourceInfo.Addresses.ElementsAs(context.Background(), &addresses, false); diags.HasError() {
		t.Fatalf("addresses: %v", diags.Errors())
	}
	if len(addresses) != 1 || addresses[0] != cngwAIServiceSourceTestNewAddress {
		t.Errorf("addresses = %v, want the cloud value [%q]", addresses, cngwAIServiceSourceTestNewAddress)
	}
}

// TestCngwAIServiceSourceReadNotFound verifies that a describe list without
// the tracked source drops the resource from state (confirmed not-found).
func TestCngwAIServiceSourceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)
	stubCngwAIServiceSourceDescribe(patches, cngwClient, nil, nil)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cngwAIServiceSourceTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the source is confirmed gone")
	}
}

// TestCngwAIServiceSourceReadApiError verifies describe failures surface as errors.
func TestCngwAIServiceSourceReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(cngwClient, "DescribeCloudNativeAPIGatewayAIServiceSourceList", func(req *cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListRequest) (*cngwv20230418.DescribeCloudNativeAPIGatewayAIServiceSourceListResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, cngwAIServiceSourceTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
}

// TestCngwAIServiceSourceUpdateSuccess verifies Update sends the PLANNED value
// to the modify request and records the POST-update cloud state from the
// read-back — not the plan.
func TestCngwAIServiceSourceUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)

	var captured *cngwv20230418.ModifyCloudNativeAPIGatewayAIServiceSourceRequest
	patches.ApplyMethodFunc(cngwClient, "ModifyCloudNativeAPIGatewayAIServiceSource", func(req *cngwv20230418.ModifyCloudNativeAPIGatewayAIServiceSourceRequest) (*cngwv20230418.ModifyCloudNativeAPIGatewayAIServiceSourceResponse, error) {
		captured = req

		return &cngwv20230418.ModifyCloudNativeAPIGatewayAIServiceSourceResponse{Response: &cngwv20230418.ModifyCloudNativeAPIGatewayAIServiceSourceResponseParams{
			RequestId: cngwAIServiceSourceTestStr("fake-request-id"),
		}}, nil
	})

	normalized := cngwAIServiceSourceTestItem()
	normalized.SourceName = cngwAIServiceSourceTestStr(cngwAIServiceSourceTestNormalName)
	stubCngwAIServiceSourceDescribe(patches, cngwClient, []*cngwv20230418.CNAPIGwAIServiceSource{normalized}, nil)

	plannedModel := cngwAIServiceSourceTestModel()
	plannedModel.SourceName = types.StringValue(cngwAIServiceSourceTestNewName)
	plannedModel.SourceInfo = cngwAIServiceSourceTestSourceInfo(cngwAIServiceSourceTestInstanceId, cngwAIServiceSourceTestNewAddress)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, cngwAIServiceSourceTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, cngwAIServiceSourceTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyCloudNativeAPIGatewayAIServiceSource request is nil")
	}
	if captured.SourceId == nil || *captured.SourceId != cngwAIServiceSourceTestId {
		t.Errorf("update request source_id = %v, want %q", captured.SourceId, cngwAIServiceSourceTestId)
	}
	if captured.SourceName == nil || *captured.SourceName != cngwAIServiceSourceTestNewName {
		t.Errorf("update request source_name = %v, want the planned %q", captured.SourceName, cngwAIServiceSourceTestNewName)
	}
	if captured.SourceInfo == nil || len(captured.SourceInfo.Addresses) != 1 || *captured.SourceInfo.Addresses[0] != cngwAIServiceSourceTestNewAddress {
		t.Errorf("update request addresses = %v, want the planned [%q]", captured.SourceInfo, cngwAIServiceSourceTestNewAddress)
	}

	var got svcngw.CngwCloudNativeAPIGatewayAIServiceSourceResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != cngwAIServiceSourceTestId {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), cngwAIServiceSourceTestId)
	}
	if got.SourceName.ValueString() != cngwAIServiceSourceTestNormalName {
		t.Errorf("source_name = %q, want the read-back cloud value %q", got.SourceName.ValueString(), cngwAIServiceSourceTestNormalName)
	}
}

// TestCngwAIServiceSourceDeleteSuccess verifies Delete sends the delete
// request with the tracked id and the stored source type.
func TestCngwAIServiceSourceDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)

	var captured *cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceRequest
	patches.ApplyMethodFunc(cngwClient, "DeleteCloudNativeAPIGatewayAIServiceSource", func(req *cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceRequest) (*cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceResponse, error) {
		captured = req

		return &cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceResponse{Response: &cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceResponseParams{
			RequestId: cngwAIServiceSourceTestStr("fake-request-id"),
		}}, nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cngwAIServiceSourceTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteCloudNativeAPIGatewayAIServiceSource request is nil")
	}
	if captured.SourceId == nil || *captured.SourceId != cngwAIServiceSourceTestId {
		t.Errorf("delete request source_id = %v, want %q", captured.SourceId, cngwAIServiceSourceTestId)
	}
	if captured.GatewayId == nil || *captured.GatewayId != cngwAIServiceSourceTestGatewayId {
		t.Errorf("delete request gateway_id = %v, want %q", captured.GatewayId, cngwAIServiceSourceTestGatewayId)
	}
	if captured.SourceType == nil || *captured.SourceType != cngwAIServiceSourceTestSourceType {
		t.Errorf("delete request source_type = %v, want %q", captured.SourceType, cngwAIServiceSourceTestSourceType)
	}
}

// TestCngwAIServiceSourceDeleteApiError verifies delete failures surface as errors.
func TestCngwAIServiceSourceDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newCngwAIServiceSourceTestContext(t)
	cngwClient := patchCngwAIServiceSourceSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(cngwClient, "DeleteCloudNativeAPIGatewayAIServiceSource", func(req *cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceRequest) (*cngwv20230418.DeleteCloudNativeAPIGatewayAIServiceSourceResponse, error) {
		return nil, fmt.Errorf("delete failed")
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, cngwAIServiceSourceTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}
