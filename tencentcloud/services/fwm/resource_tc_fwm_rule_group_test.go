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

// buildFwmRuleGroupResp builds a DescribeRuleGroups item. The resource only
// consumes GroupId / GroupName / Product, so the other SDK fields are omitted.
func buildFwmRuleGroupResp(groupId, groupName, product string) *fwmv20250611.RuleGroupResp {
	return &fwmv20250611.RuleGroupResp{
		GroupId:   ptrStringFwmRuleGroup(groupId),
		GroupName: ptrStringFwmRuleGroup(groupName),
		Product:   ptrStringFwmRuleGroup(product),
	}
}

// mockDescribeRuleGroups makes DescribeRuleGroupsWithContext return a fixed
// result set. It also asserts the lookup filter built by DescribeFwmRuleGroupById.
func mockDescribeRuleGroups(t *testing.T, patches *gomonkey.Patches, client *fwmv20250611.Client, ruleGroups ...*fwmv20250611.RuleGroupResp) {
	patches.ApplyMethodFunc(client, "DescribeRuleGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		if len(request.Filters) > 0 {
			assert.Equal(t, "GroupId", *request.Filters[0].Name)
			assert.NotEmpty(t, request.Filters[0].Values)
		}

		resp := fwmv20250611.NewDescribeRuleGroupsResponse()
		resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
			TotalCount: func() *int64 { v := int64(len(ruleGroups)); return &v }(),
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

	groupId := res.Schema["group_id"]
	assert.NotNil(t, groupId)
	assert.True(t, groupId.Computed)
	assert.False(t, groupId.Required)

	// Only these three attributes are exposed; the previously published
	// read-only attributes are no longer part of the schema.
	assert.Len(t, res.Schema, 3)
	for _, key := range []string{
		"app_id", "rule_count", "create_time", "update_time", "create_by",
		"create_name", "update_by", "update_name", "receive_account",
		"exec_areas", "direction",
	} {
		_, ok := res.Schema[key]
		assert.False(t, ok, "schema %s should not exist", key)
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

	mockDescribeRuleGroups(t, patches, fwmClient, buildFwmRuleGroupResp("fwmrg-wsrrhuvrib", "tf-example", "enterprise_sg"))

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
	assert.Equal(t, "tf-example", d.Get("group_name"))
	assert.Equal(t, "enterprise_sg", d.Get("product"))
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Get("group_id"))
}

// TestFwmRuleGroup_Create_EmptyGroupId verifies an empty GroupId in the create
// response is treated as a failure instead of producing an empty resource id.
func TestFwmRuleGroup_Create_EmptyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			GroupId:   ptrStringFwmRuleGroup(""),
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

// TestFwmRuleGroup_Create_NilResponse verifies a nil Response body fails fast.
func TestFwmRuleGroup_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		return &fwmv20250611.CreateRuleGroupResponse{}, nil
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

// TestFwmRuleGroup_Create_WaitVisible verifies Create polls DescribeRuleGroups
// until the new group becomes visible.
func TestFwmRuleGroup_Create_WaitVisible(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.CreateRuleGroupRequest) (*fwmv20250611.CreateRuleGroupResponse, error) {
		resp := fwmv20250611.NewCreateRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateRuleGroupResponseParams{
			GroupId:   ptrStringFwmRuleGroup("fwmrg-wsrrhuvrib"),
			RequestId: ptrStringFwmRuleGroup("fake-request-id"),
		}
		return resp, nil
	})

	calls := 0
	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, _ *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		calls++
		resp := fwmv20250611.NewDescribeRuleGroupsResponse()
		resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
			RequestId: ptrStringFwmRuleGroup("fake-request-id"),
		}
		if calls > 1 {
			resp.Response.RuleGroups = []*fwmv20250611.RuleGroupResp{
				buildFwmRuleGroupResp("fwmrg-wsrrhuvrib", "tf-example", "enterprise_sg"),
			}
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
	assert.NoError(t, err)
	assert.Greater(t, calls, 1)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
	assert.Equal(t, "tf-example", d.Get("group_name"))
}

func TestFwmRuleGroup_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeRuleGroups(t, patches, fwmClient, buildFwmRuleGroupResp("fwmrg-wsrrhuvrib", "tf-example", "enterprise_sg"))

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "old-name",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Get("group_id"))
	assert.Equal(t, "tf-example", d.Get("group_name"))
	assert.Equal(t, "enterprise_sg", d.Get("product"))
}

// TestFwmRuleGroup_Read_NotFound verifies the resource is dropped from state
// when the group no longer exists.
func TestFwmRuleGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeRuleGroups(t, patches, fwmClient)

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

// TestFwmRuleGroup_Read_NilResponse verifies a nil Response body surfaces as an
// error instead of clearing the resource.
func TestFwmRuleGroup_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeRuleGroupsWithContext", func(_ context.Context, _ *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		return &fwmv20250611.DescribeRuleGroupsResponse{}, nil
	})

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "fwmrg-wsrrhuvrib", d.Id())
}

// TestFwmRuleGroup_Read_IgnoresOtherGroups ensures the lookup picks the group
// matching the resource id rather than the first item of the response.
func TestFwmRuleGroup_Read_IgnoresOtherGroups(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeRuleGroups(t, patches, fwmClient,
		buildFwmRuleGroupResp("fwmrg-other", "other-name", "cfw_nat_acl"),
		buildFwmRuleGroupResp("fwmrg-wsrrhuvrib", "tf-example", "enterprise_sg"),
	)

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf-example", d.Get("group_name"))
	assert.Equal(t, "enterprise_sg", d.Get("product"))
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

	mockDescribeRuleGroups(t, patches, fwmClient, buildFwmRuleGroupResp("fwmrg-wsrrhuvrib", "tf-example-new", "enterprise_sg"))

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
	assert.Equal(t, "tf-example-new", d.Get("group_name"))
}

// TestFwmRuleGroup_Update_SkipWhenUnchanged verifies ModifyRuleGroup is not
// called when group_name does not change (Update still refreshes state).
func TestFwmRuleGroup_Update_SkipWhenUnchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)

	modifyCalled := false
	patches.ApplyMethodFunc(fwmClient, "ModifyRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
		modifyCalled = true
		return fwmv20250611.NewModifyRuleGroupResponse(), nil
	})

	mockDescribeRuleGroups(t, patches, fwmClient, buildFwmRuleGroupResp("fwmrg-wsrrhuvrib", "tf-example", "enterprise_sg"))

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
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
	assert.Equal(t, "tf-example", d.Get("group_name"))
}

func TestFwmRuleGroup_Delete(t *testing.T) {
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

// TestFwmRuleGroup_Delete_NilResponse verifies a nil Response body is reported
// as an error so a failed deletion is not silently swallowed.
func TestFwmRuleGroup_Delete_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmRuleGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DeleteRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		return &fwmv20250611.DeleteRuleGroupResponse{}, nil
	})

	meta := newMockMetaForFwmRuleGroup()
	res := fwm.ResourceTencentCloudFwmRuleGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "tf-example",
		"product":    "enterprise_sg",
	})
	d.SetId("fwmrg-wsrrhuvrib")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}
