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

type mockMetaFwmRuleGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmRuleGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmRuleGroup{}

func newMockMetaFwmRuleGroup() *mockMetaFwmRuleGroup {
	return &mockMetaFwmRuleGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwmRG(s string) *string {
	return &s
}

func ptrInt64FwmRG(i int64) *int64 {
	return &i
}

func buildFwmRuleGroupResponse(groupId string) *fwmv20250611.DescribeRuleGroupsResponse {
	resp := fwmv20250611.NewDescribeRuleGroupsResponse()
	resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
		TotalCount: ptrInt64FwmRG(1),
		RuleGroups: []*fwmv20250611.RuleGroupResp{
			{
				GroupId:   ptrStringFwmRG(groupId),
				GroupName: ptrStringFwmRG("tf-example"),
				Product:   ptrStringFwmRG("2"),
				AppId:     ptrStringFwmRG("1300000001"),
				RuleCount: ptrInt64FwmRG(5),
			},
		},
	}
	return resp
}

func TestFwmRuleGroupResource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "tf-example", *request.GroupName)
		assert.NotNil(t, request.Product)
		assert.Equal(t, "enterprise_sg", *request.Product)

		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			GroupId:   ptrStringFwmRG("fwmrg-wsrrhuvrib"),
			RequestId: ptrStringFwmRG("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "GroupId", *request.Filters[0].Name)
		assert.Equal(t, 1, len(request.Filters[0].Values))
		assert.Equal(t, "fwmrg-wsrrhuvrib", *request.Filters[0].Values[0])
		return buildFwmRuleGroupResponse("fwmrg-wsrrhuvrib"), nil
	})

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Get("group_id").(string))
	assert.Equal(t, "tf-example", d.Get("group_name").(string))
}

func TestFwmRuleGroupResource_CreateEmptyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			RequestId: ptrStringFwmRG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GroupId")
	assert.Equal(t, "", d.Id())
}

func TestFwmRuleGroupResource_ReadProductNormalize(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeRuleGroupsResponse()
		resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
			TotalCount: ptrInt64FwmRG(1),
			RuleGroups: []*fwmv20250611.RuleGroupResp{
				{
					GroupId:    ptrStringFwmRG("fwmrg-wsrrhuvrib"),
					GroupName:  ptrStringFwmRG("tf-example"),
					Product:    ptrStringFwmRG("2"),
					AppId:      ptrStringFwmRG("1300000001"),
					RuleCount:  ptrInt64FwmRG(3),
					CreateTime: ptrStringFwmRG("2025-01-01 00:00:00"),
					UpdateTime: ptrStringFwmRG("2025-01-02 00:00:00"),
					CreateBy:   ptrStringFwmRG("100000000001"),
					CreateName: ptrStringFwmRG("creator"),
					UpdateBy:   ptrStringFwmRG("100000000002"),
					UpdateName: ptrStringFwmRG("updater"),
					ReceiveAccount: []*fwmv20250611.ReceiveAccount{
						{
							Uin:          ptrStringFwmRG("100000000003"),
							Nickname:     ptrStringFwmRG("account"),
							ReceiverType: ptrInt64FwmRG(1),
							Members: []*fwmv20250611.MemberInfo{
								{
									AppId:    ptrStringFwmRG("1300000002"),
									Uin:      ptrStringFwmRG("100000000004"),
									Nickname: ptrStringFwmRG("member"),
									MemberId: ptrStringFwmRG("m-xxxx"),
									NodeName: ptrStringFwmRG("dept-a"),
								},
							},
						},
					},
					ExecAreas: []*string{ptrStringFwmRG("vpc"), ptrStringFwmRG("edge")},
					Direction: ptrInt64FwmRG(-1),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "cfw_nat_acl",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Get("group_id").(string))
	assert.Equal(t, "tf-example", d.Get("group_name").(string))
	assert.Equal(t, "cfw_nat_acl", d.Get("product").(string))
	assert.Equal(t, "1300000001", d.Get("app_id").(string))
	assert.Equal(t, 3, d.Get("rule_count").(int))
	assert.Equal(t, -1, d.Get("direction").(int))

	execAreas := d.Get("exec_areas").([]interface{})
	assert.Equal(t, 2, len(execAreas))
	assert.Equal(t, "vpc", execAreas[0].(string))

	receiveAccount := d.Get("receive_account").([]interface{})
	assert.Equal(t, 1, len(receiveAccount))
	receiveAccountMap := receiveAccount[0].(map[string]interface{})
	assert.Equal(t, "100000000003", receiveAccountMap["uin"].(string))
	assert.Equal(t, 1, receiveAccountMap["receiver_type"].(int))

	members := receiveAccountMap["members"].([]interface{})
	assert.Equal(t, 1, len(members))
	memberMap := members[0].(map[string]interface{})
	assert.Equal(t, "100000000004", memberMap["uin"].(string))
	assert.Equal(t, "m-xxxx", memberMap["member_id"].(string))
	assert.Equal(t, "dept-a", memberMap["node_name"].(string))
}

func TestFwmRuleGroupResource_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeRuleGroupsResponse()
		resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
			TotalCount: ptrInt64FwmRG(0),
			RuleGroups: []*fwmv20250611.RuleGroupResp{},
		}
		return resp, nil
	})

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmRuleGroupResource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "fwmrg-wsrrhuvrib", *request.GroupId)
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "tf-example-renamed", *request.GroupName)

		resp := fwmv20250611.NewModifyRuleGroupResponse()
		resp.Response = &fwmv20250611.ModifyRuleGroupResponseParams{
			RequestId: ptrStringFwmRG("fake-request-id-modify"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		return buildFwmRuleGroupResponse("fwmrg-wsrrhuvrib"), nil
	})

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example-renamed",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "group_name"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestFwmRuleGroupResource_UpdateImmutableArg(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "cfw_nat_acl",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "product"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product")
	assert.Contains(t, err.Error(), "immutable")
}

func TestFwmRuleGroupResource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		assert.NotNil(t, request.GroupIds)
		assert.Equal(t, 1, len(request.GroupIds))
		assert.Equal(t, "fwmrg-wsrrhuvrib", *request.GroupIds[0])

		resp := fwmv20250611.NewDeleteRuleGroupResponse()
		resp.Response = &fwmv20250611.DeleteRuleGroupResponseParams{
			RequestId: ptrStringFwmRG("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
