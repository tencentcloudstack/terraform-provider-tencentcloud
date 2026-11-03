package fwm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func ptrStringFwmDispatch(s string) *string {
	return &s
}

func ptrInt64FwmDispatch(i int64) *int64 {
	return &i
}

func ptrFloat64FwmDispatch(f float64) *float64 {
	return &f
}

func fwmDispatchStrategySchema(t *testing.T) action_schema.Schema {
	a := fwm.NewFwmDispatchStrategy()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

func fwmDispatchStrategyTimeoutsObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}
}

func nullFwmDispatchStrategyTimeoutsValue() tftypes.Value {
	return tftypes.NewValue(fwmDispatchStrategyTimeoutsObjType(), nil)
}

func newFwmDispatchStrategyInvokeRequest(t *testing.T, product, uin, subAccountUin string) action.InvokeRequest {
	return newFwmDispatchStrategyInvokeRequestWithTimeouts(t, product, uin, subAccountUin, nullFwmDispatchStrategyTimeoutsValue())
}

func newFwmDispatchStrategyInvokeRequestWithTimeouts(t *testing.T, product, uin, subAccountUin string, timeoutsVal tftypes.Value) action.InvokeRequest {
	s := fwmDispatchStrategySchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"product":         tftypes.String,
			"uin":             tftypes.String,
			"sub_account_uin": tftypes.String,
			"timeouts":        fwmDispatchStrategyTimeoutsObjType(),
		},
	}, map[string]tftypes.Value{
		"product":         tftypes.NewValue(tftypes.String, product),
		"uin":             tftypes.NewValue(tftypes.String, uin),
		"sub_account_uin": tftypes.NewValue(tftypes.String, subAccountUin),
		"timeouts":        timeoutsVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func setFwmDispatchStrategyClient(a *fwm.FwmDispatchStrategy, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// shortenPollInterval shortens the package-level poll interval so the
// Wait-loop tests do not sleep for seconds.
func shortenPollInterval(t *testing.T) func() {
	t.Helper()
	orig := fwm.FwmDispatchStrategyPollIntervalForTest()
	fwm.SetFwmDispatchStrategyPollIntervalForTest(1 * time.Millisecond)
	return func() {
		fwm.SetFwmDispatchStrategyPollIntervalForTest(orig)
	}
}

func dispatchStatus(status int64, progress float64, errMsg string) *fwmv20250611.DescribeStrategyDispatchStatusResponseParams {
	return &fwmv20250611.DescribeStrategyDispatchStatusResponseParams{
		Status:   ptrInt64FwmDispatch(status),
		Progress: ptrFloat64FwmDispatch(progress),
		ErrorMsg: ptrStringFwmDispatch(errMsg),
	}
}

// patchFwmClient installs a UseFwmV20250611Client patch returning a stub client.
func patchFwmClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) {
	patches.ApplyMethodFunc(client, "UseFwmV20250611Client", func() *fwmv20250611.Client {
		return &fwmv20250611.Client{}
	})
}

// patchStatus installs a DescribeStrategyDispatchStatusWithContext patch.
func patchStatus(patches *gomonkey.Patches, fn func(ctx context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error)) {
	patches.ApplyMethodFunc(&fwmv20250611.Client{}, "DescribeStrategyDispatchStatusWithContext", fn)
}

// patchDispatch installs a DispatchStrategyWithContext patch.
func patchDispatch(patches *gomonkey.Patches, fn func(ctx context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error)) {
	patches.ApplyMethodFunc(&fwmv20250611.Client{}, "DispatchStrategyWithContext", fn)
}

// go test ./tencentcloud/services/fwm/ -run "TestFwmDispatchStrategy" -v -count=1 -gcflags="all=-l"

// TestFwmDispatchStrategy_Invoke_Success tests a successful dispatch: pre-check
// sees no in-progress task, dispatch is triggered, then polling reaches succeeded.
func TestFwmDispatchStrategy_Invoke_Success(t *testing.T) {
	defer shortenPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	var dispatched bool
	var gotReq *fwmv20250611.DispatchStrategyRequest
	patchDispatch(patches, func(_ context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error) {
		dispatched = true
		gotReq = req
		resp := fwmv20250611.NewDispatchStrategyResponse()
		resp.Response = &fwmv20250611.DispatchStrategyResponseParams{RequestId: ptrStringFwmDispatch("fake-request-id")}
		return resp, nil
	})

	// pre-check sees succeeded (old task) => not in progress; poll then succeeds.
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		if req.Product == nil || *req.Product != "enterprise_sg" {
			return nil, fmt.Errorf("expected Product=enterprise_sg in status query")
		}
		resp := fwmv20250611.NewDescribeStrategyDispatchStatusResponse()
		resp.Response = dispatchStatus(fwm.DispatchStatusSucceeded, 100, "")
		return resp, nil
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, "enterprise_sg", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.True(t, dispatched, "expected DispatchStrategy to be called")
	assert.NotNil(t, gotReq)
	assert.NotNil(t, gotReq.Status)
	assert.Equal(t, int64(1), *gotReq.Status)
	assert.NotNil(t, gotReq.Product)
	assert.Equal(t, "enterprise_sg", *gotReq.Product)
}

// TestFwmDispatchStrategy_Invoke_ReuseInProgress tests that an in-progress task
// is reused: pre-check returns dispatching, so DispatchStrategy is NOT called.
func TestFwmDispatchStrategy_Invoke_ReuseInProgress(t *testing.T) {
	defer shortenPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	dispatched := false
	patchDispatch(patches, func(_ context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error) {
		dispatched = true
		return fwmv20250611.NewDispatchStrategyResponse(), nil
	})

	// Every status query returns dispatching; polling stays in-progress but the
	// short timeout forces convergence-report via failure. To keep this a
	// success-path test, return succeeded after the first query.
	calls := 0
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		calls++
		st := fwm.DispatchStatusDispatching
		if calls >= 2 {
			st = fwm.DispatchStatusSucceeded
		}
		resp := fwmv20250611.NewDescribeStrategyDispatchStatusResponse()
		resp.Response = dispatchStatus(st, 50, "")
		return resp, nil
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, "cfw_vpc_acl", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.False(t, dispatched, "expected DispatchStrategy NOT to be called when a task is in progress")
}

// TestFwmDispatchStrategy_Invoke_DispatchError tests API error from DispatchStrategy.
func TestFwmDispatchStrategy_Invoke_DispatchError(t *testing.T) {
	defer shortenPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	// pre-check: old task succeeded => not in progress => dispatch needed.
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyDispatchStatusResponse()
		resp.Response = dispatchStatus(fwm.DispatchStatusSucceeded, 100, "")
		return resp, nil
	})

	patchDispatch(patches, func(_ context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Product")
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, "cfw_nat_acl", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "dispatching FWM strategy")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InvalidParameter")
}

// TestFwmDispatchStrategy_Invoke_StatusQueryError tests pre-check status query error.
func TestFwmDispatchStrategy_Invoke_StatusQueryError(t *testing.T) {
	defer shortenPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=boom")
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, "cfw_edge_acl", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "checking FWM dispatch status")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InternalError")
}

// TestFwmDispatchStrategy_Invoke_PollFailed tests dispatch succeeds but polling
// reaches Status=3 (failed) carrying ErrorMsg. The pre-check also sees failed
// (not in-progress 0/1), so dispatch is triggered; the poll then immediately
// observes the failed task.
func TestFwmDispatchStrategy_Invoke_PollFailed(t *testing.T) {
	defer shortenPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	// Every status query returns failed with an error message.
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyDispatchStatusResponse()
		resp.Response = dispatchStatus(fwm.DispatchStatusFailed, 0, "downstream sync error")
		return resp, nil
	})

	patchDispatch(patches, func(_ context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error) {
		resp := fwmv20250611.NewDispatchStrategyResponse()
		resp.Response = &fwmv20250611.DispatchStrategyResponseParams{RequestId: ptrStringFwmDispatch("ok")}
		return resp, nil
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, "enterprise_sg", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for FWM dispatch strategy")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "downstream sync error")
}

// TestFwmDispatchStrategy_Invoke_PollTimeout tests dispatch succeeds but polling
// never converges within the (very short) configured timeout.
func TestFwmDispatchStrategy_Invoke_PollTimeout(t *testing.T) {
	defer shortenPollInterval(t)()
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	// pre-check: old task succeeded => not in progress => dispatch.
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyDispatchStatusResponse()
		resp.Response = dispatchStatus(fwm.DispatchStatusSucceeded, 100, "")
		return resp, nil
	})

	patchDispatch(patches, func(_ context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error) {
		return fwmv20250611.NewDispatchStrategyResponse(), nil
	})

	// After dispatch, always report dispatching (never converges). Replace the
	// earlier status patch with one that always returns dispatching.
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyDispatchStatusResponse()
		resp.Response = dispatchStatus(fwm.DispatchStatusDispatching, 42, "")
		return resp, nil
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	timeoutVal := tftypes.NewValue(fwmDispatchStrategyTimeoutsObjType(), map[string]tftypes.Value{
		"invoke": tftypes.NewValue(tftypes.String, "50ms"),
	})
	req := newFwmDispatchStrategyInvokeRequestWithTimeouts(t, "cfw_vpc_acl", "100011949846", "", timeoutVal)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "waiting for FWM dispatch strategy")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "did not converge")
}

// TestFwmDispatchStrategy_Invoke_InvalidProduct tests invalid product rejection
// before any API call.
func TestFwmDispatchStrategy_Invoke_InvalidProduct(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	apiCalled := false
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		apiCalled = true
		return fwmv20250611.NewDescribeStrategyDispatchStatusResponse(), nil
	})
	patchDispatch(patches, func(_ context.Context, req *fwmv20250611.DispatchStrategyRequest) (*fwmv20250611.DispatchStrategyResponse, error) {
		apiCalled = true
		return fwmv20250611.NewDispatchStrategyResponse(), nil
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, "all", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid product")
	assert.False(t, apiCalled)
}

// TestFwmDispatchStrategy_Invoke_MissingInput tests missing required inputs are
// rejected before any API call.
func TestFwmDispatchStrategy_Invoke_MissingInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	patchFwmClient(patches, client)

	apiCalled := false
	patchStatus(patches, func(_ context.Context, req *fwmv20250611.DescribeStrategyDispatchStatusRequest) (*fwmv20250611.DescribeStrategyDispatchStatusResponse, error) {
		apiCalled = true
		return fwmv20250611.NewDescribeStrategyDispatchStatusResponse(), nil
	})

	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	// empty product
	req := newFwmDispatchStrategyInvokeRequest(t, "", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing product")
	assert.False(t, apiCalled)

	// empty uin
	req = newFwmDispatchStrategyInvokeRequest(t, "enterprise_sg", "", "")
	resp = &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing uin")
	assert.False(t, apiCalled)
}

// TestFwmDispatchStrategy_Invoke_ClientNotConfigured tests the nil-client guard.
func TestFwmDispatchStrategy_Invoke_ClientNotConfigured(t *testing.T) {
	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)

	req := newFwmDispatchStrategyInvokeRequest(t, "enterprise_sg", "100011949846", "")
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestFwmDispatchStrategy_MetadataAndSchema validates metadata and schema
// definition.
func TestFwmDispatchStrategy_MetadataAndSchema(t *testing.T) {
	aImpl := fwm.NewFwmDispatchStrategy()
	a := aImpl.(*fwm.FwmDispatchStrategy)

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_fwm_dispatch_strategy", metaResp.TypeName)

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "product")
	assert.Contains(t, attrs, "uin")
	assert.Contains(t, attrs, "sub_account_uin")
	assert.True(t, attrs["product"].IsRequired())
	assert.True(t, attrs["uin"].IsRequired())
	assert.True(t, attrs["sub_account_uin"].IsOptional())
	// Operation has no Computed / status output attribute.
	assert.NotContains(t, attrs, "status")

	blocks := schemaResp.Schema.Blocks
	assert.Contains(t, blocks, "timeouts")
	timeoutsBlock, ok := blocks["timeouts"].(action_schema.SingleNestedBlock)
	assert.True(t, ok)
	assert.IsType(t, actiontimeouts.Type{}, timeoutsBlock.CustomType)
	invokeAttr, ok := timeoutsBlock.Attributes["invoke"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, invokeAttr.Optional)
	assert.NotEmpty(t, invokeAttr.Validators)
}