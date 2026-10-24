package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

type mockMetaForFwmRuleGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForFwmRuleGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForFwmRuleGroup{}

func newMockMetaForFwmRuleGroup() *mockMetaForFwmRuleGroup {
	return &mockMetaForFwmRuleGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwmRuleGroup(s string) *string {
	return &s
}

func ptrInt64FwmRuleGroup(v int64) *int64 {
	return &v
}

func buildFwmRuleGroupResp() *fwmv20250611.RuleGroupResp {
	return &fwmv20250611.RuleGroupResp{
		AppId:      ptrStringFwmRuleGroup("1300000000"),
		GroupId:    ptrStringFwmRuleGroup("fwmrg-wsrrhuvrib"),
		GroupName:  ptrStringFwmRuleGroup("tf-example"),
		RuleCount:  ptrInt64FwmRuleGroup(3),
		Product:    ptrStringFwmRuleGroup("1"),
		CreateTime: ptrStringFwmRuleGroup("2025-10-21 10:00:00"),
		UpdateTime: ptrStringFwmRuleGroup("2025-10-21 11:00:00"),
		CreateBy:   ptrStringFwmRuleGroup("100000000001"),
		CreateName: ptrStringFwmRuleGroup("tf-creator"),
		UpdateBy:   ptrStringFwmRuleGroup("100000000002"),
		UpdateName: ptrStringFwmRuleGroup("tf-updater"),
		ReceiveAccount: []*fwmv20250611.ReceiveAccount{
			{
				Uin:          ptrStringFwmRuleGroup("100000000003"),
				Nickname:     ptrStringFwmRuleGroup("tf-receiver"),
				ReceiverType: ptrInt64FwmRuleGroup(1),
				Members: []*fwmv20250611.MemberInfo{
					{
						AppId:    ptrStringFwmRuleGroup("1300000001"),
						Uin:      ptrStringFwmRuleGroup("100000000004"),
						Nickname: ptrStringFwmRuleGroup("tf-member"),
						MemberId: ptrStringFwmRuleGroup("member-001"),
						NodeName: ptrStringFwmRuleGroup("node-001"),
					},
				},
			},
		},
		ExecAreas: []*string{ptrStringFwmRuleGroup("post_official")},
		Direction: ptrInt64FwmRuleGroup(1),
	}
}

func mockFwmRuleGroupsWithContext(patches *gomonkey.Patches, fwmClient *fwmv20250611.Client, ruleGroups ...*fwmv20250611.RuleGroupResp) {
	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, _ *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeRuleGroupsResponse()
		resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
			TotalCount: ptrInt64FwmRuleGroup(int64(len(ruleGroups))),
			RuleGroups: ruleGroups,
			RequestId:  ptrStringFwmRuleGroup("fake-request-id"),
		}
		return resp, nil
	})
}

func TestFwmRuleGroup_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmRuleGroup()

	groupName := res.Schema["group_name"]
	assert.NotNil(t, groupName)
	assert.True(t, groupName.Required)
	assert.False(t, groupName.ForceNew)
	assert.False(t, groupName.Computed)

	product := res.Schema["product"]
	assert.NotNil(t, product)
	assert.True(t, product.Required)
	assert.True(t, product.ForceNew)
	assert.False(t, product.Computed)
	if product.ValidateFunc != nil {
		_, errs := product.ValidateFunc("cfw_vpc_acl", "product")
		assert.Empty(t, errs)
		_, errs = product.ValidateFunc("all", "product")
		assert.NotEmpty(t, errs)
	}

	for _, key := range []string{
		"group_id", "app_id", "rule_count", "create_time", "update_time",
		"create_by", "create_name", "update_by", "update_name", "receive_account",
		"exec_areas", "direction",
	} {
		item, ok := res.Schema[key]
		assert.True(t, ok, "schema %s should exist", key)
		assert.True(t, item.Computed, "schema %s should be computed", key)
		assert.False(t, item.Required, "schema %s should not be required", key)
	}

	assert.NotNil(t, res.Importer)
}

func TestFwmRuleGroup_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "tf-example", *request.GroupName)
		assert.NotNil(t, request.Product)
		assert.Equal(t, "enterprise_sg", *request.Product)

		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			GroupId:   ptrStringFwmRuleGroup("fwmrg-wsrrhuvrib"),
			RequestId: ptrStringFwmRuleGroup("fake-request-id"),
		}
		return resp, nil
	})

	mockFwmRuleGroupsWithContext(patches, fwmClient, buildFwmRuleGroupResp())

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
	assert.Equal(t, "tf-example", d.Get("group_name").(string))
	assert.Equal(t, "enterprise_sg", d.Get("product").(string))
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Get("group_id").(string))
}

func TestFwmRuleGroup_Create_NilGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			RequestId: ptrStringFwmRuleGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmRuleGroup_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	mockFwmRuleGroupsWithContext(patches, fwmClient, buildFwmRuleGroupResp())

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Get("group_id").(string))
	assert.Equal(t, "tf-example", d.Get("group_name").(string))
	assert.Equal(t, "enterprise_sg", d.Get("product").(string))
	assert.Equal(t, "1300000000", d.Get("app_id").(string))
	assert.Equal(t, 3, d.Get("rule_count").(int))
	assert.Equal(t, "2025-10-21 10:00:00", d.Get("create_time").(string))
	assert.Equal(t, "2025-10-21 11:00:00", d.Get("update_time").(string))
	assert.Equal(t, "100000000001", d.Get("create_by").(string))
	assert.Equal(t, "tf-creator", d.Get("create_name").(string))
	assert.Equal(t, "100000000002", d.Get("update_by").(string))
	assert.Equal(t, "tf-updater", d.Get("update_name").(string))
	assert.Equal(t, 1, d.Get("direction").(int))

	receiveAccounts := d.Get("receive_account").([]interface{})
	assert.Len(t, receiveAccounts, 1)
	receiveAccount := receiveAccounts[0].(map[string]interface{})
	assert.Equal(t, "100000000003", receiveAccount["uin"].(string))
	assert.Equal(t, "tf-receiver", receiveAccount["nickname"].(string))
	assert.Equal(t, 1, receiveAccount["receiver_type"].(int))

	members := receiveAccount["members"].([]interface{})
	assert.Len(t, members, 1)
	member := members[0].(map[string]interface{})
	assert.Equal(t, "1300000001", member["app_id"].(string))
	assert.Equal(t, "100000000004", member["uin"].(string))
	assert.Equal(t, "tf-member", member["nickname"].(string))
	assert.Equal(t, "member-001", member["member_id"].(string))
	assert.Equal(t, "node-001", member["node_name"].(string))

	execAreas := d.Get("exec_areas").([]interface{})
	assert.Len(t, execAreas, 1)
	assert.Equal(t, "post_official", execAreas[0].(string))
}

func TestFwmRuleGroup_Read_ProductCodeMapping(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	ruleGroup := buildFwmRuleGroupResp()
	ruleGroup.Product = ptrStringFwmRuleGroup("2")
	mockFwmRuleGroupsWithContext(patches, fwmClient, ruleGroup)

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "cfw_nat_acl",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cfw_nat_acl", d.Get("product").(string))
}

func TestFwmRuleGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	mockFwmRuleGroupsWithContext(patches, fwmClient)

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmRuleGroup_Update_GroupName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	modifyCalled := false
	patches.ApplyMethodFunc(fwmClient, "ModifyRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-wsrrhuvrib", *request.GroupId)
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "tf-example-new", *request.GroupName)
		modifyCalled = true

		resp := fwmv20250611.NewModifyRuleGroupResponse()
		resp.Response = &fwmv20250611.ModifyRuleGroupResponseParams{
			RequestId: ptrStringFwmRuleGroup("fake-request-id"),
		}
		return resp, nil
	})

	mockFwmRuleGroupsWithContext(patches, fwmClient, buildFwmRuleGroupResp())

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	state := &terraform.InstanceState{
		ID: "fwmrg-wsrrhuvrib",
		Attributes: map[string]string{
			"id":         "fwmrg-wsrrhuvrib",
			"group_id":   "fwmrg-wsrrhuvrib",
			"group_name": "tf-example",
			"product":    "enterprise_sg",
		},
	}
	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"group_name": "tf-example-new",
		"product":    "enterprise_sg",
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
}

func TestFwmRuleGroup_Delete_SingleGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		assert.Len(t, request.GroupIds, 1)
		assert.NotNil(t, request.GroupIds[0])
		assert.Equal(t, "fwmrg-wsrrhuvrib", *request.GroupIds[0])

		resp := fwmv20250611.NewDeleteRuleGroupResponse()
		resp.Response = &fwmv20250611.DeleteRuleGroupResponseParams{
			RequestId: ptrStringFwmRuleGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
