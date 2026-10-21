package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	fwmsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

// ---- gomonkey-based unit tests for tencentcloud_fwm_strategy ----
// Run with: go test ./tencentcloud/services/fwm/ -run "TestFwmStrategy" -v -count=1 -gcflags="all=-l"

type mockMetaFwmStrategy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmStrategy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmStrategy{}

func newMockMetaFwmStrategy() *mockMetaFwmStrategy {
	return &mockMetaFwmStrategy{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStrFwm(s string) *string {
	return &s
}

func ptrInt64Fwm(i int64) *int64 {
	return &i
}

// TestFwmStrategy_Schema validates the schema definition.
func TestFwmStrategy_Schema(t *testing.T) {
	res := fwmsvc.ResourceTencentCloudFwmStrategy()

	productField, ok := res.Schema["product"]
	assert.True(t, ok, "product should exist in schema")
	assert.True(t, productField.Required)
	assert.True(t, productField.ForceNew)

	execAreaField, ok := res.Schema["exec_area"]
	assert.True(t, ok, "exec_area should exist in schema")
	assert.True(t, execAreaField.Required)
	assert.True(t, execAreaField.ForceNew)

	groupIdField, ok := res.Schema["group_id"]
	assert.True(t, ok, "group_id should exist in schema")
	assert.False(t, groupIdField.Required)

	sequenceField, ok := res.Schema["sequence"]
	assert.True(t, ok, "sequence should exist in schema")
	assert.Equal(t, schema.TypeInt, sequenceField.Type)

	receiveAccountField, ok := res.Schema["receive_account"]
	assert.True(t, ok, "receive_account should exist in schema")
	assert.Equal(t, schema.TypeList, receiveAccountField.Type)

	strategyIdField, ok := res.Schema["strategy_id"]
	assert.True(t, ok, "strategy_id should exist in schema")
	assert.True(t, strategyIdField.Computed)
}

// TestFwmStrategy_Create verifies Create sets the strategy id and reads back state.
func TestFwmStrategy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.CreateStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(ctx context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		capturedRequest = request
		resp := fwmv20250611.NewCreateStrategyResponse()
		resp.Response = &fwmv20250611.CreateStrategyResponseParams{
			StrategyIds: []*string{ptrStrFwm("fwms-xxxx")},
			RequestId:   ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeStrategy called by the read after create.
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy: &fwmv20250611.StrategyResp{
				StrategyId: ptrStrFwm("fwms-xxxx"),
				GroupId:    ptrStrFwm("fwmg-yyyy"),
				GroupName:  ptrStrFwm("tf-example-group"),
				RuleCount:  ptrInt64Fwm(10),
				Sequence:   ptrInt64Fwm(1),
				ReceiveAccount: []*fwmv20250611.ReceiveAccount{
					{
						Uin:          ptrStrFwm("100028274628"),
						ReceiverType: ptrInt64Fwm(0),
					},
				},
				CreateTime: ptrStrFwm("2025-01-01 00:00:00"),
			},
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product":         "enterprise_sg",
		"sub_account_uin": "100028274628",
		"exec_area":       "pre",
		"group_id":        "fwmg-xxxxxx",
		"sequence":        1,
		"receive_account": []interface{}{"100028274628"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// The id should be the strategy id from the response.
	assert.Equal(t, "fwms-xxxx", d.Id())

	// Verify exec_area=pre caused the strategy to be sent via PreStrategy.
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.PreStrategy)
	assert.Equal(t, 1, len(capturedRequest.PreStrategy))
	assert.Equal(t, "fwmg-xxxxxx", *capturedRequest.PreStrategy[0].GroupId)

	// Verify read-back values.
	assert.Equal(t, "fwms-xxxx", d.Get("strategy_id").(string))
	assert.Equal(t, "fwmg-yyyy", d.Get("group_id").(string))
	assert.Equal(t, 10, d.Get("rule_count").(int))
}

// TestFwmStrategy_Create_EmptyId verifies Create fails when StrategyIds is empty.
func TestFwmStrategy_Create_EmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(ctx context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		resp := fwmv20250611.NewCreateStrategyResponse()
		resp.Response = &fwmv20250611.CreateStrategyResponseParams{
			StrategyIds: []*string{},
			RequestId:   ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product":         "enterprise_sg",
		"sub_account_uin": "100028274628",
		"exec_area":       "pre",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestFwmStrategy_ReadValue verifies Read sets fields when API returns values.
func TestFwmStrategy_ReadValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy: &fwmv20250611.StrategyResp{
				StrategyId: ptrStrFwm("fwms-xxxx"),
				GroupId:    ptrStrFwm("fwmg-yyyy"),
				GroupName:  ptrStrFwm("tf-example-group"),
				RuleCount:  ptrInt64Fwm(10),
				RuleStatus: ptrInt64Fwm(1),
				Sequence:   ptrInt64Fwm(1),
				ReceiveAccount: []*fwmv20250611.ReceiveAccount{
					{
						Uin:          ptrStrFwm("100028274628"),
						ReceiverType: ptrInt64Fwm(0),
					},
					{
						Uin:          ptrStrFwm("fwmg-group-1"),
						ReceiverType: ptrInt64Fwm(1),
					},
				},
				CreateTime: ptrStrFwm("2025-01-01 00:00:00"),
				UpdateTime: ptrStrFwm("2025-01-02 00:00:00"),
				CreateBy:   ptrStrFwm("creator"),
				UpdateBy:   ptrStrFwm("updater"),
			},
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify receive_account and receive_group are split by ReceiverType.
	receiveAccount := d.Get("receive_account").([]interface{})
	assert.Equal(t, 1, len(receiveAccount))
	assert.Equal(t, "100028274628", receiveAccount[0].(string))

	receiveGroup := d.Get("receive_group").([]interface{})
	assert.Equal(t, 1, len(receiveGroup))
	assert.Equal(t, "fwmg-group-1", receiveGroup[0].(string))

	assert.Equal(t, "tf-example-group", d.Get("group_name").(string))
	assert.Equal(t, 1, d.Get("rule_status").(int))
}

// TestFwmStrategy_ReadNil verifies Read clears id when strategy is nil.
func TestFwmStrategy_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy:  nil,
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestFwmStrategy_Update verifies field change triggers ModifyStrategy with incremental semantics.
func TestFwmStrategy_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.ModifyStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "ModifyStrategyWithContext", func(ctx context.Context, request *fwmv20250611.ModifyStrategyRequest) (*fwmv20250611.ModifyStrategyResponse, error) {
		capturedRequest = request
		resp := fwmv20250611.NewModifyStrategyResponse()
		resp.Response = &fwmv20250611.ModifyStrategyResponseParams{
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeStrategy called by the read after update.
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy: &fwmv20250611.StrategyResp{
				StrategyId: ptrStrFwm("fwms-xxxx"),
				GroupId:    ptrStrFwm("fwmg-new"),
				Sequence:   ptrInt64Fwm(2),
				ReceiveAccount: []*fwmv20250611.ReceiveAccount{
					{
						Uin:          ptrStrFwm("100028274628"),
						ReceiverType: ptrInt64Fwm(0),
					},
				},
			},
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product":         "enterprise_sg",
		"sub_account_uin": "100028274628",
		"exec_area":       "pre",
		"group_id":        "fwmg-new",
		"sequence":        2,
		"receive_account": []interface{}{"100028274628"},
	})
	d.SetId("fwms-xxxx")

	// Simulate that only group_id changed (incremental semantics: only GroupId should be sent).
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "group_id"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify incremental semantics: only GroupId should be set (sequence/receive_account unchanged).
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.GroupId)
	assert.Equal(t, "fwmg-new", *capturedRequest.GroupId)
	// Sequence was not changed, so it should be nil.
	assert.Nil(t, capturedRequest.Sequence)
	// ReceiveAccount was not changed, so it should be nil.
	assert.Nil(t, capturedRequest.ReceiveAccount)
}

// TestFwmStrategy_Delete verifies Delete calls DeleteStrategy with the strategy id.
func TestFwmStrategy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.DeleteStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "DeleteStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DeleteStrategyRequest) (*fwmv20250611.DeleteStrategyResponse, error) {
		capturedRequest = request
		resp := fwmv20250611.NewDeleteStrategyResponse()
		resp.Response = &fwmv20250611.DeleteStrategyResponseParams{
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("fwms-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	// Verify StrategyIds contains a single element with the id.
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, 1, len(capturedRequest.StrategyIds))
	assert.Equal(t, "fwms-xxxx", *capturedRequest.StrategyIds[0])
}
