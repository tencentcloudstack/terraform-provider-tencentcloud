package fwm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

// buildStrategyResp builds a DescribeStrategy item. execArea selects whether the
// resource should write the result into pre_strategy or post_strategy.
func buildStrategyResp(strategyId, groupId string, sequence int64, execArea string, receivers ...*fwmv20250611.ReceiveAccount) *fwmv20250611.StrategyResp {
	resp := &fwmv20250611.StrategyResp{
		StrategyId: ptrStrFwm(strategyId),
		GroupId:    ptrStrFwm(groupId),
		Sequence:   ptrInt64Fwm(sequence),
	}
	if execArea != "" {
		resp.ExecArea = ptrStrFwm(execArea)
	}
	if len(receivers) > 0 {
		resp.ReceiveAccount = receivers
	}
	return resp
}

func buildReceiver(uin string, receiverType int64) *fwmv20250611.ReceiveAccount {
	return &fwmv20250611.ReceiveAccount{
		Uin:          ptrStrFwm(uin),
		ReceiverType: ptrInt64Fwm(receiverType),
	}
}

// mockDescribeStrategy makes DescribeStrategyWithContext return the given item.
func mockDescribeStrategy(patches *gomonkey.Patches, client *fwmv20250611.Client, strategy *fwmv20250611.StrategyResp) {
	patches.ApplyMethodFunc(client, "DescribeStrategyWithContext", func(_ context.Context, _ *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		resp := fwmv20250611.NewDescribeStrategyResponse()
		resp.Response = &fwmv20250611.DescribeStrategyResponseParams{
			Strategy:  strategy,
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})
}

// mockCreateStrategy makes CreateStrategyWithContext return the given ids and
// records the outgoing request.
func mockCreateStrategy(patches *gomonkey.Patches, client *fwmv20250611.Client, capture **fwmv20250611.CreateStrategyRequest, ids ...string) {
	patches.ApplyMethodFunc(client, "CreateStrategyWithContext", func(_ context.Context, request *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		*capture = request

		strategyIds := make([]*string, 0, len(ids))
		for _, id := range ids {
			strategyIds = append(strategyIds, ptrStrFwm(id))
		}

		resp := fwmv20250611.NewCreateStrategyResponse()
		resp.Response = &fwmv20250611.CreateStrategyResponseParams{
			StrategyIds: strategyIds,
			RequestId:   ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})
}

func derefStringPtrs(ptrs []*string) []string {
	out := make([]string, 0, len(ptrs))
	for _, p := range ptrs {
		if p != nil {
			out = append(out, *p)
		}
	}
	return out
}

func TestFwmStrategy_Schema(t *testing.T) {
	res := fwmsvc.ResourceTencentCloudFwmStrategy()

	product := res.Schema["product"]
	assert.NotNil(t, product)
	assert.True(t, product.Required)
	assert.True(t, product.ForceNew)

	receiveAccount := res.Schema["receive_account"]
	assert.NotNil(t, receiveAccount)
	assert.True(t, receiveAccount.Optional)
	assert.Equal(t, schema.TypeSet, receiveAccount.Type)

	receiveGroup := res.Schema["receive_group"]
	assert.NotNil(t, receiveGroup)
	assert.True(t, receiveGroup.Optional)
	assert.Equal(t, schema.TypeSet, receiveGroup.Type)

	pre := res.Schema["pre_strategy"]
	assert.NotNil(t, pre)
	assert.True(t, pre.Optional)
	assert.Equal(t, schema.TypeList, pre.Type)
	assert.Equal(t, 1, pre.MaxItems)
	assert.Equal(t, []string{"post_strategy"}, pre.ConflictsWith)
	preElem := pre.Elem.(*schema.Resource)
	assert.True(t, preElem.Schema["group_id"].Required)
	assert.True(t, preElem.Schema["sequence"].Computed)

	post := res.Schema["post_strategy"]
	assert.NotNil(t, post)
	assert.True(t, post.Optional)
	assert.Equal(t, schema.TypeList, post.Type)
	assert.Equal(t, 1, post.MaxItems)
	assert.Equal(t, []string{"pre_strategy"}, post.ConflictsWith)
	postElem := post.Elem.(*schema.Resource)
	assert.True(t, postElem.Schema["group_id"].Required)
	assert.True(t, postElem.Schema["sequence"].Computed)

	strategyId := res.Schema["strategy_id"]
	assert.NotNil(t, strategyId)
	assert.True(t, strategyId.Computed)

	// Attributes of the previous generation are gone for good.
	for _, key := range []string{"exec_area", "group_id", "sequence", "rule_count", "group_name", "rule_status", "sub_account_uin", "create_time", "update_time"} {
		_, ok := res.Schema[key]
		assert.False(t, ok, "schema %s should not exist", key)
	}

	assert.NotNil(t, res.Importer)
}

func TestFwmStrategy_Create_PreStrategy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var captured *fwmv20250611.CreateStrategyRequest
	mockCreateStrategy(patches, fwmClient, &captured, "fwms-xxxx")
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-pre", 1, "pre", buildReceiver("100028274628", 0)))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product":         "enterprise_sg",
		"receive_account": []interface{}{"100028274628"},
		"receive_group":   []interface{}{"fwmg-group-1"},
		"pre_strategy": []interface{}{
			map[string]interface{}{"group_id": "fwmg-pre"},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwms-xxxx", d.Id())

	assert.NotNil(t, captured)
	assert.Equal(t, "enterprise_sg", *captured.Product)
	assert.ElementsMatch(t, []string{"100028274628"}, derefStringPtrs(captured.ReceiveAccount))
	assert.ElementsMatch(t, []string{"fwmg-group-1"}, derefStringPtrs(captured.ReceiveGroup))
	assert.Len(t, captured.PreStrategy, 1)
	assert.Equal(t, "fwmg-pre", *captured.PreStrategy[0].GroupId)
	assert.Equal(t, int64(1), *captured.PreStrategy[0].Sequence)
	assert.Nil(t, captured.PostStrategy)

	// read back
	assert.Equal(t, "fwms-xxxx", d.Get("strategy_id").(string))
	preStrategies := d.Get("pre_strategy").([]interface{})
	assert.Len(t, preStrategies, 1)
	assert.Equal(t, "fwmg-pre", preStrategies[0].(map[string]interface{})["group_id"])
	assert.Equal(t, 1, preStrategies[0].(map[string]interface{})["sequence"])
}

func TestFwmStrategy_Create_PostStrategy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var captured *fwmv20250611.CreateStrategyRequest
	mockCreateStrategy(patches, fwmClient, &captured, "fwms-xxxx")
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-post", 1, "post"))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
		"post_strategy": []interface{}{
			map[string]interface{}{"group_id": "fwmg-post"},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwms-xxxx", d.Id())

	assert.NotNil(t, captured.PostStrategy)
	assert.Len(t, captured.PostStrategy, 1)
	assert.Equal(t, "fwmg-post", *captured.PostStrategy[0].GroupId)
	assert.Nil(t, captured.PreStrategy)

	postStrategies := d.Get("post_strategy").([]interface{})
	assert.Len(t, postStrategies, 1)
	assert.Equal(t, "fwmg-post", postStrategies[0].(map[string]interface{})["group_id"])
}

func TestFwmStrategy_Create_EmptyStrategyIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var captured *fwmv20250611.CreateStrategyRequest
	mockCreateStrategy(patches, fwmClient, &captured)

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

func TestFwmStrategy_Create_EmptyFirstId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var captured *fwmv20250611.CreateStrategyRequest
	mockCreateStrategy(patches, fwmClient, &captured, "")

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

func TestFwmStrategy_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "CreateStrategyWithContext", func(_ context.Context, _ *fwmv20250611.CreateStrategyRequest) (*fwmv20250611.CreateStrategyResponse, error) {
		return &fwmv20250611.CreateStrategyResponse{}, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

func TestFwmStrategy_Read_PreStrategy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-pre", 2, "pre",
		buildReceiver("100028274628", 0),
		buildReceiver("fwmg-group-1", 1),
	))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwms-xxxx", d.Id())
	assert.Equal(t, "fwms-xxxx", d.Get("strategy_id").(string))

	accounts := d.Get("receive_account").(*schema.Set)
	assert.Equal(t, 1, accounts.Len())
	assert.True(t, accounts.Contains("100028274628"))

	groups := d.Get("receive_group").(*schema.Set)
	assert.Equal(t, 1, groups.Len())
	assert.True(t, groups.Contains("fwmg-group-1"))

	preStrategies := d.Get("pre_strategy").([]interface{})
	assert.Len(t, preStrategies, 1)
	assert.Equal(t, "fwmg-pre", preStrategies[0].(map[string]interface{})["group_id"])
	assert.Equal(t, 2, preStrategies[0].(map[string]interface{})["sequence"])

	assert.Empty(t, d.Get("post_strategy").([]interface{}))
}

func TestFwmStrategy_Read_PostStrategy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-post", 3, "post"))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	postStrategies := d.Get("post_strategy").([]interface{})
	assert.Len(t, postStrategies, 1)
	assert.Equal(t, "fwmg-post", postStrategies[0].(map[string]interface{})["group_id"])
	assert.Equal(t, 3, postStrategies[0].(map[string]interface{})["sequence"])

	assert.Empty(t, d.Get("pre_strategy").([]interface{}))
}

// TestFwmStrategy_Read_ResourceNotFound verifies a ResourceNotFound error is
// swallowed and the resource is dropped from state.
func TestFwmStrategy_Read_ResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, _ *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		return nil, errors.New("ResourceNotFound: strategy not exist")
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmStrategy_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyWithContext", func(_ context.Context, _ *fwmv20250611.DescribeStrategyRequest) (*fwmv20250611.DescribeStrategyResponse, error) {
		return &fwmv20250611.DescribeStrategyResponse{}, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "fwms-xxxx", d.Id())
}

// TestFwmStrategy_Read_NoExecArea verifies a strategy without exec_area keeps
// both strategy blocks untouched.
func TestFwmStrategy_Read_NoExecArea(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-pre", 1, ""))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Get("pre_strategy").([]interface{}))
	assert.Empty(t, d.Get("post_strategy").([]interface{}))
}

func buildStrategyState(t *testing.T, res *schema.Resource, id string, attrs map[string]interface{}) *terraform.InstanceState {
	oldData := schema.TestResourceDataRaw(t, res.Schema, attrs)
	oldData.SetId(id)
	return oldData.State()
}

func TestFwmStrategy_Update_ReceiveAccount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var captured *fwmv20250611.ModifyStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "ModifyStrategyWithContext", func(_ context.Context, request *fwmv20250611.ModifyStrategyRequest) (*fwmv20250611.ModifyStrategyResponse, error) {
		captured = request
		resp := fwmv20250611.NewModifyStrategyResponse()
		resp.Response = &fwmv20250611.ModifyStrategyResponseParams{
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-pre", 1, "pre",
		buildReceiver("100028274628", 0),
		buildReceiver("100028274629", 0),
	))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()

	state := buildStrategyState(t, res, "fwms-xxxx", map[string]interface{}{
		"product":         "enterprise_sg",
		"receive_account": []interface{}{"100028274628"},
	})
	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"product":         "enterprise_sg",
		"receive_account": []interface{}{"100028274628", "100028274629"},
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, captured)
	assert.Equal(t, "fwms-xxxx", *captured.StrategyId)
	assert.ElementsMatch(t, []string{"100028274628", "100028274629"}, derefStringPtrs(captured.ReceiveAccount))
	assert.Nil(t, captured.GroupId)
	assert.Nil(t, captured.ReceiveGroup)

	accounts := d.Get("receive_account").(*schema.Set)
	assert.Equal(t, 2, accounts.Len())
}

func TestFwmStrategy_Update_SkipWhenUnchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	modifyCalled := false
	patches.ApplyMethodFunc(fwmClient, "ModifyStrategyWithContext", func(_ context.Context, _ *fwmv20250611.ModifyStrategyRequest) (*fwmv20250611.ModifyStrategyResponse, error) {
		modifyCalled = true
		return fwmv20250611.NewModifyStrategyResponse(), nil
	})
	mockDescribeStrategy(patches, fwmClient, buildStrategyResp("fwms-xxxx", "fwmg-pre", 1, "pre", buildReceiver("100028274628", 0)))

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()

	attrs := map[string]interface{}{
		"product":         "enterprise_sg",
		"receive_account": []interface{}{"100028274628"},
	}
	state := buildStrategyState(t, res, "fwms-xxxx", attrs)
	rawConfig := terraform.NewResourceConfigRaw(attrs)
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)

	accounts := d.Get("receive_account").(*schema.Set)
	assert.Equal(t, 1, accounts.Len())
}

func TestFwmStrategy_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)

	var captured *fwmv20250611.DeleteStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "DeleteStrategyWithContext", func(_ context.Context, request *fwmv20250611.DeleteStrategyRequest) (*fwmv20250611.DeleteStrategyResponse, error) {
		captured = request
		resp := fwmv20250611.NewDeleteStrategyResponse()
		resp.Response = &fwmv20250611.DeleteStrategyResponseParams{
			RequestId: ptrStrFwm("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, captured)
	assert.Equal(t, []string{"fwms-xxxx"}, derefStringPtrs(captured.StrategyIds))
}

func TestFwmStrategy_Delete_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmStrategy().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DeleteStrategyWithContext", func(_ context.Context, _ *fwmv20250611.DeleteStrategyRequest) (*fwmv20250611.DeleteStrategyResponse, error) {
		return &fwmv20250611.DeleteStrategyResponse{}, nil
	})

	meta := newMockMetaFwmStrategy()
	res := fwmsvc.ResourceTencentCloudFwmStrategy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"product": "enterprise_sg",
	})
	d.SetId("fwms-xxxx")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}
