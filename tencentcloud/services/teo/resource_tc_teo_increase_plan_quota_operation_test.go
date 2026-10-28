package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestIncreasePlanQuota" -v -count=1 -gcflags="all=-l"

// mockClient returns a *connectivity.TencentCloudClient backed by an empty
// instance. Individual test cases patch the method on the concrete client to
// return a mock teo client.
func mockClient() *connectivity.TencentCloudClient {
	return &connectivity.TencentCloudClient{}
}

// actionSchema retrieves the schema from the action for test helpers.
func actionSchema(t *testing.T) actionschema.Schema {
	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	resp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema error: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// buildConfig builds an action.InvokeRequest.Config (tfsdk.Config) from the
// given plan_id / quota_type / quota_number values by constructing the
// underlying tftypes.Value directly.
func buildConfig(t *testing.T, planId, quotaType string, quotaNumber int64) tfsdk.Config {
	sch := actionSchema(t)
	val := tftypes.NewValue(
		sch.Type().TerraformType(context.Background()),
		map[string]tftypes.Value{
			"plan_id":      tftypes.NewValue(tftypes.String, planId),
			"quota_type":   tftypes.NewValue(tftypes.String, quotaType),
			"quota_number": tftypes.NewValue(tftypes.Number, quotaNumber),
			"deal_name":    tftypes.NewValue(tftypes.String, nil),
		},
	)
	return tfsdk.Config{
		Raw:    val,
		Schema: sch,
	}
}

// TestIncreasePlanQuota_Success tests a successful invocation.
func TestIncreasePlanQuota_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mockClient(), "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(ctx context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		assert.Equal(t, "edgeone-2unuvzjmmn2q", *request.PlanId)
		assert.Equal(t, "site", *request.QuotaType)
		assert.Equal(t, int64(10), *request.QuotaNumber)

		resp := teov20220901.NewIncreasePlanQuotaResponse()
		resp.Response = &teov20220901.IncreasePlanQuotaResponseParams{
			DealName:  ptrString("202312290001"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	a.SetClientForTest(mockClient())

	invReq := action.InvokeRequest{Config: buildConfig(t, "edgeone-2unuvzjmmn2q", "site", 10)}
	invResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invReq, invResp)

	assert.False(t, invResp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", invResp.Diagnostics)
}

// TestIncreasePlanQuota_APIError tests handling API errors.
func TestIncreasePlanQuota_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mockClient(), "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(ctx context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid plan_id")
	})

	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	a.SetClientForTest(mockClient())

	invReq := action.InvokeRequest{Config: buildConfig(t, "invalid-plan", "site", 10)}
	invResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invReq, invResp)

	assert.True(t, invResp.Diagnostics.HasError())
}

// TestIncreasePlanQuota_NilResponse tests handling a nil response.
func TestIncreasePlanQuota_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mockClient(), "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(ctx context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		return nil, nil
	})

	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	a.SetClientForTest(mockClient())

	invReq := action.InvokeRequest{Config: buildConfig(t, "edgeone-2unuvzjmmn2q", "site", 10)}
	invResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invReq, invResp)

	assert.True(t, invResp.Diagnostics.HasError())
}

// TestIncreasePlanQuota_NilDealName tests handling nil DealName in response.
func TestIncreasePlanQuota_NilDealName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mockClient(), "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IncreasePlanQuotaWithContext", func(ctx context.Context, request *teov20220901.IncreasePlanQuotaRequest) (*teov20220901.IncreasePlanQuotaResponse, error) {
		resp := teov20220901.NewIncreasePlanQuotaResponse()
		resp.Response = &teov20220901.IncreasePlanQuotaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	a.SetClientForTest(mockClient())

	invReq := action.InvokeRequest{Config: buildConfig(t, "edgeone-2unuvzjmmn2q", "site", 10)}
	invResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invReq, invResp)

	assert.True(t, invResp.Diagnostics.HasError())
}

// TestIncreasePlanQuota_MissingPlanId tests error when plan_id is missing.
func TestIncreasePlanQuota_MissingPlanId(t *testing.T) {
	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	a.SetClientForTest(mockClient())

	invReq := action.InvokeRequest{Config: buildConfig(t, "", "site", 10)}
	invResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invReq, invResp)

	assert.True(t, invResp.Diagnostics.HasError())
	foundPlanId := false
	for _, d := range invResp.Diagnostics {
		if contains(d.Summary(), "plan_id") {
			foundPlanId = true
		}
	}
	assert.True(t, foundPlanId, "expected an error mentioning plan_id")
}

// TestIncreasePlanQuota_MissingQuotaType tests error when quota_type is missing.
func TestIncreasePlanQuota_MissingQuotaType(t *testing.T) {
	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	a.SetClientForTest(mockClient())

	invReq := action.InvokeRequest{Config: buildConfig(t, "edgeone-2unuvzjmmn2q", "", 10)}
	invResp := &action.InvokeResponse{}
	a.Invoke(context.Background(), invReq, invResp)

	assert.True(t, invResp.Diagnostics.HasError())
	foundQuotaType := false
	for _, d := range invResp.Diagnostics {
		if contains(d.Summary(), "quota_type") {
			foundQuotaType = true
		}
	}
	assert.True(t, foundQuotaType, "expected an error mentioning quota_type")
}

// TestIncreasePlanQuota_Schema validates the schema definition.
func TestIncreasePlanQuota_Schema(t *testing.T) {
	sch := actionSchema(t)
	assert.Contains(t, sch.Attributes, "plan_id")
	assert.Contains(t, sch.Attributes, "quota_type")
	assert.Contains(t, sch.Attributes, "quota_number")
	assert.Contains(t, sch.Attributes, "deal_name")
}

// TestIncreasePlanQuota_Metadata validates the metadata TypeName.
func TestIncreasePlanQuota_Metadata(t *testing.T) {
	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	metaReq := action.MetadataRequest{ProviderTypeName: "tencentcloud"}
	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), metaReq, metaResp)
	assert.Equal(t, "tencentcloud_teo_increase_plan_quota", metaResp.TypeName)
}

// TestIncreasePlanQuota_ConfigureNilProviderData ensures Configure with nil
// ProviderData is a silent no-op.
func TestIncreasePlanQuota_ConfigureNilProviderData(t *testing.T) {
	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	resp := &action.ConfigureResponse{}
	a.Configure(context.Background(), action.ConfigureRequest{}, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

// TestIncreasePlanQuota_ConfigureWrongType ensures Configure adds an error
// diagnostic when ProviderData is of an unexpected type.
func TestIncreasePlanQuota_ConfigureWrongType(t *testing.T) {
	a := teo.NewTeoIncreasePlanQuotaAction().(*teo.TeoIncreasePlanQuotaAction)
	resp := &action.ConfigureResponse{}
	a.Configure(context.Background(), action.ConfigureRequest{ProviderData: "not-a-meta"}, resp)
	assert.True(t, resp.Diagnostics.HasError())
}

// contains is a small helper to check substrings in diagnostics.
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
