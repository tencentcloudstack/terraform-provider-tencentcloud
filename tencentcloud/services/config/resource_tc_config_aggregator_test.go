package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	config_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcconfig "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

type mockMetaForConfigAggregator struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForConfigAggregator) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForConfigAggregator{}

func newMockMetaForConfigAggregator() *mockMetaForConfigAggregator {
	return &mockMetaForConfigAggregator{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func buildConfigAggregatorDescribeResponse() *config_sdk.DescribeAggregatorResponse {
	resp := config_sdk.NewDescribeAggregatorResponse()
	resp.Response = &config_sdk.DescribeAggregatorResponseParams{
		Name:        helper.String("tf-example-aggregator"),
		Description: helper.String("tf example aggregator"),
		Type:        helper.String("CUSTOM"),
		AggregatorStatus: func() *uint64 {
			s := uint64(1)
			return &s
		}(),
		AggregatorAccounts: []*config_sdk.AggregatorAccount{
			{
				MemberUin:  helper.IntUint64(100012345679),
				MemberName: helper.String("member-1"),
			},
			{
				MemberUin:  helper.IntUint64(100012345680),
				MemberName: helper.String("member-2"),
			},
		},
	}
	return resp
}

// TestConfigAggregator_Schema verifies the schema fields match the design.
func TestConfigAggregator_Schema(t *testing.T) {
	res := svcconfig.ResourceTencentCloudConfigAggregator()

	cases := []struct {
		field    string
		typ      schema.ValueType
		required bool
		forceNew bool
	}{
		{"name", schema.TypeString, true, false},
		{"description", schema.TypeString, true, false},
		{"type", schema.TypeString, true, true},
		{"owner_uin", schema.TypeString, true, true},
		{"account_group_id", schema.TypeString, false, false},
		{"aggregator_status", schema.TypeInt, false, false},
	}

	for _, c := range cases {
		s, ok := res.Schema[c.field]
		assert.True(t, ok, "field %s should exist in schema", c.field)
		assert.Equal(t, c.typ, s.Type, "field %s type mismatch", c.field)
		if c.required {
			assert.True(t, s.Required, "field %s should be required", c.field)
		} else {
			assert.True(t, s.Computed, "field %s should be computed", c.field)
		}
		if c.forceNew {
			assert.True(t, s.ForceNew, "field %s should be ForceNew", c.field)
		}
	}

	_, ok := res.Schema["aggregator_accounts"]
	assert.True(t, ok, "aggregator_accounts should exist in schema")
	assert.Equal(t, schema.TypeList, res.Schema["aggregator_accounts"].Type)
	assert.True(t, res.Schema["aggregator_accounts"].Optional)
}

// TestConfigAggregator_Create covers the Create business logic end-to-end with mocked cloud API.
func TestConfigAggregator_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAggregator().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.CreateAggregatorRequest
	patches.ApplyMethodFunc(configClient, "CreateAggregatorWithContext", func(ctx context.Context, request *config_sdk.CreateAggregatorRequest) (*config_sdk.CreateAggregatorResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewCreateAggregatorResponse()
		resp.Response = &config_sdk.CreateAggregatorResponseParams{
			AccountGroupId: helper.String("ca-xxxxxxxx"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(configClient, "DescribeAggregatorWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregatorRequest) (*config_sdk.DescribeAggregatorResponse, error) {
		return buildConfigAggregatorDescribeResponse(), nil
	})

	meta := newMockMetaForConfigAggregator()
	res := svcconfig.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf-example-aggregator",
		"description": "tf example aggregator",
		"type":        "CUSTOM",
		"owner_uin":   "100012345678",
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"member_uin":  100012345679,
				"member_name": "member-1",
			},
			map[string]interface{}{
				"member_uin":  100012345680,
				"member_name": "member-2",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ca-xxxxxxxx#100012345678", d.Id())
	assert.Equal(t, "ca-xxxxxxxx", d.Get("account_group_id"))
	assert.Equal(t, "100012345678", d.Get("owner_uin"))

	assert.NotNil(t, capturedRequest, "CreateAggregatorWithContext should be called")
	assert.Equal(t, "tf-example-aggregator", *capturedRequest.Name)
	assert.Equal(t, "CUSTOM", *capturedRequest.Type)
	assert.Len(t, capturedRequest.AggregatorAccounts, 2)
	assert.Equal(t, uint64(100012345679), *capturedRequest.AggregatorAccounts[0].MemberUin)
	assert.Equal(t, "member-2", *capturedRequest.AggregatorAccounts[1].MemberName)
}

// TestConfigAggregator_Read covers the Read business logic (populate state from cloud).
func TestConfigAggregator_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAggregator().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DescribeAggregatorWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregatorRequest) (*config_sdk.DescribeAggregatorResponse, error) {
		return buildConfigAggregatorDescribeResponse(), nil
	})

	meta := newMockMetaForConfigAggregator()
	res := svcconfig.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ca-xxxxxxxx#100012345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ca-xxxxxxxx#100012345678", d.Id())
	assert.Equal(t, "tf-example-aggregator", d.Get("name"))
	assert.Equal(t, "tf example aggregator", d.Get("description"))
	assert.Equal(t, "CUSTOM", d.Get("type"))
	assert.Equal(t, 1, d.Get("aggregator_status"))
	assert.Equal(t, "ca-xxxxxxxx", d.Get("account_group_id"))
	assert.Equal(t, "100012345678", d.Get("owner_uin"))

	accounts := d.Get("aggregator_accounts").([]interface{})
	assert.Len(t, accounts, 2)
	first := accounts[0].(map[string]interface{})
	assert.Equal(t, 100012345679, first["member_uin"])
	assert.Equal(t, "member-1", first["member_name"])
}

// TestConfigAggregator_Update covers the Update business logic (name/description/members change).
func TestConfigAggregator_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAggregator().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.UpdateAggregatorRequest
	patches.ApplyMethodFunc(configClient, "UpdateAggregatorWithContext", func(ctx context.Context, request *config_sdk.UpdateAggregatorRequest) (*config_sdk.UpdateAggregatorResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewUpdateAggregatorResponse()
		resp.Response = &config_sdk.UpdateAggregatorResponseParams{}
		return resp, nil
	})
	patches.ApplyMethodFunc(configClient, "DescribeAggregatorWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregatorRequest) (*config_sdk.DescribeAggregatorResponse, error) {
		return buildConfigAggregatorDescribeResponse(), nil
	})

	meta := newMockMetaForConfigAggregator()
	res := svcconfig.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf-example-aggregator",
		"description": "tf example aggregator",
		"type":        "CUSTOM",
		"owner_uin":   "100012345678",
		"aggregator_accounts": []interface{}{
			map[string]interface{}{
				"member_uin":  100012345679,
				"member_name": "member-1",
			},
		},
	})
	d.SetId("ca-xxxxxxxx#100012345678")

	// Simulate a name + member change.
	assert.NoError(t, d.Set("name", "tf-example-aggregator-updated"))
	assert.NoError(t, d.Set("aggregator_accounts", []interface{}{
		map[string]interface{}{
			"member_uin":  100012345681,
			"member_name": "member-3",
		},
	}))

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest, "UpdateAggregatorWithContext should be called")
	assert.Equal(t, "ca-xxxxxxxx", *capturedRequest.AccountGroupId)
	assert.Equal(t, "tf-example-aggregator-updated", *capturedRequest.Name)
	assert.Equal(t, uint64(100012345678), *capturedRequest.OwnerUin)
	assert.Len(t, capturedRequest.AggregatorAccounts, 1)
	assert.Equal(t, uint64(100012345681), *capturedRequest.AggregatorAccounts[0].MemberUin)
	assert.Equal(t, "member-3", *capturedRequest.AggregatorAccounts[0].MemberName)
}

// TestConfigAggregator_Delete covers the Delete business logic.
func TestConfigAggregator_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAggregator().client, "UseConfigV20220802Client", configClient)

	var capturedRequest *config_sdk.DeleteAggregatorsRequest
	patches.ApplyMethodFunc(configClient, "DeleteAggregatorsWithContext", func(ctx context.Context, request *config_sdk.DeleteAggregatorsRequest) (*config_sdk.DeleteAggregatorsResponse, error) {
		capturedRequest = request
		resp := config_sdk.NewDeleteAggregatorsResponse()
		resp.Response = &config_sdk.DeleteAggregatorsResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForConfigAggregator()
	res := svcconfig.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ca-xxxxxxxx#100012345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest, "DeleteAggregatorsWithContext should be called")
	assert.Equal(t, "ca-xxxxxxxx", *capturedRequest.AccountGroupId)
	assert.Equal(t, uint64(100012345678), *capturedRequest.OwnerUin)
}

// TestConfigAggregator_Create_NilResponse verifies the create nil-response guard returns an error.
func TestConfigAggregator_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAggregator().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "CreateAggregatorWithContext", func(ctx context.Context, request *config_sdk.CreateAggregatorRequest) (*config_sdk.CreateAggregatorResponse, error) {
		resp := config_sdk.NewCreateAggregatorResponse()
		resp.Response = &config_sdk.CreateAggregatorResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForConfigAggregator()
	res := svcconfig.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":        "tf-example-aggregator",
		"description": "tf example aggregator",
		"type":        "CUSTOM",
		"owner_uin":   "100012345678",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestConfigAggregator_Read_NotFound verifies that a nil response clears the state id.
func TestConfigAggregator_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &config_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForConfigAggregator().client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "DescribeAggregatorWithContext", func(ctx context.Context, request *config_sdk.DescribeAggregatorRequest) (*config_sdk.DescribeAggregatorResponse, error) {
		resp := config_sdk.NewDescribeAggregatorResponse()
		resp.Response = &config_sdk.DescribeAggregatorResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForConfigAggregator()
	res := svcconfig.ResourceTencentCloudConfigAggregator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ca-xxxxxxxx#100012345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}
