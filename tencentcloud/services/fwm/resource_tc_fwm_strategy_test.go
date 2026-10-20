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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

type mockMetaFwmStrategy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmStrategy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmStrategy{}

func newMockMetaFwmStrategy() *mockMetaFwmStrategy {
	return &mockMetaFwmStrategy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwm(s string) *string { return &s }

// TestFwmStrategy_Create tests the Create operation
func TestFwmStrategy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(ctx context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		assert.Equal(t, "cfm", *request.Product)
		assert.Equal(t, 1, len(request.ReceiveAccount))
		assert.Equal(t, "10000xxx", *request.ReceiveAccount[0])
		assert.Equal(t, 1, len(request.PreStrategy))
		assert.Equal(t, "rule-group-id-1", *request.PreStrategy[0].GroupId)
		assert.Equal(t, int64(1), *request.PreStrategy[0].Sequence)
		assert.Equal(t, 1, len(request.PostStrategy))
		assert.Equal(t, "rule-group-id-2", *request.PostStrategy[0].GroupId)
		assert.Equal(t, int64(1), *request.PostStrategy[0].Sequence)
		assert.Equal(t, 1, len(request.ReceiveGroup))
		assert.Equal(t, "group-xxx", *request.ReceiveGroup[0])

		resp := fwmv20250611.NewCreateStrategyResponse()
		resp.Response = &fwmv20250611.CreateStrategyResponseParams{
			StrategyIds: []*string{ptrStringFwm("s1"), ptrStringFwm("s2")},
			RequestId:   ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy: &fwmv20250611.StrategyResp{
				StrategyId: ptrStringFwm(*request.StrategyId),
				GroupId:    ptrStringFwm("rule-group-id-1"),
			},
			RequestId: ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwm.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product":         "cfm",
		"receive_account": []interface{}{"10000xxx"},
		"pre_strategy": []interface{}{
			map[string]interface{}{
				"group_id": "rule-group-id-1",
				"sequence": 1,
			},
		},
		"post_strategy": []interface{}{
			map[string]interface{}{
				"group_id": "rule-group-id-2",
				"sequence": 1,
			},
		},
		"receive_group": []interface{}{"group-xxx"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "s1#s2", d.Id())
	strategyIds := d.Get("strategy_ids").([]interface{})
	assert.Equal(t, 2, len(strategyIds))
	assert.Equal(t, "s1", strategyIds[0].(string))
	assert.Equal(t, "s2", strategyIds[1].(string))
}

// TestFwmStrategy_Create_EmptyStrategyIds tests Create when StrategyIds is empty
func TestFwmStrategy_Create_EmptyStrategyIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(ctx context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		resp := fwmv20250611.NewCreateStrategyResponse()
		resp.Response = &fwmv20250611.CreateStrategyResponseParams{
			StrategyIds: []*string{},
			RequestId:   ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwm.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "cfm",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmStrategy_Read tests the Read operation
func TestFwmStrategy_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy: &fwmv20250611.StrategyResp{
				StrategyId: ptrStringFwm(*request.StrategyId),
				GroupId:    ptrStringFwm("rule-group-id-1"),
			},
			RequestId: ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwm.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "cfm",
	})
	d.SetId("s1#s2")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "s1#s2", d.Id())
}

// TestFwmStrategy_Read_Empty tests Read when all strategies are not found
func TestFwmStrategy_Read_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy:  nil,
			RequestId: ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwm.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("s1#s2")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmStrategy_Delete tests the Delete operation
func TestFwmStrategy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	deleteCalled := false
	patches.ApplyMethodFunc(fwmClient, "DeleteStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DeleteStrategyRequest) (*fwmv20250611.DeleteStrategyResponse, error) {
		deleteCalled = true
		assert.Equal(t, 2, len(request.StrategyIds))
		assert.Equal(t, "s1", *request.StrategyIds[0])
		assert.Equal(t, "s2", *request.StrategyIds[1])

		resp := fwmv20250611.NewDeleteStrategyResponse()
		resp.Response = &fwmv20250611.DeleteStrategyResponseParams{
			RequestId: ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwm.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("s1#s2")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

// TestFwmStrategy_Update_Immutable tests Update rejects immutable field changes
func TestFwmStrategy_Update_Immutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(ctx context.Context, request *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy: &fwmv20250611.StrategyResp{
				StrategyId: ptrStringFwm(*request.StrategyId),
			},
			RequestId: ptrStringFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwm.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "cfm",
	})
	d.SetId("s1")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "pre_strategy"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pre_strategy is immutable")
}
