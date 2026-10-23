package fwm_test

import (
	"context"
	"testing"

	"github.com/agiledagon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

type mockMetaFwmMemberGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmMemberGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmMemberGroup{}

func newMockMetaFwmMemberGroup() *mockMetaFwmMemberGroup {
	return &mockMetaFwmMemberGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwm(s string) *string {
	return &s
}

func buildOrganMemberItem(memberId, appId, groupId string) *fwmv20250611.OrganMemberItem {
	item := &fwmv20250611.OrganMemberItem{
		MemberId: ptrStringFwm(memberId),
		AppId:    ptrStringFwm(appId),
	}
	if groupId != "" {
		item.AccountGroup = &fwmv20250611.AccountGroupInfo{
			GroupId:   ptrStringFwm(groupId),
			GroupName: ptrStringFwm("group-name"),
		}
	}
	return item
}

func mockDescribeOrganMembers(patches *gomonkey.Patches, fwmClient *fwmv20250611.Client, total int64, members []*fwmv20250611.OrganMemberItem) {
	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId:  ptrStringFwm("fake-request-id"),
			TotalCount: &total,
			Members:    members,
		}
		return resp, nil
	})
}

// buildFwmMemberGroupData builds a *schema.ResourceData with the given old state
// (key->string) and new raw config (key->interface{}) so that d.HasChange
// returns true for fields that differ between old and new.
func buildFwmMemberGroupData(t *testing.T, id string, oldState map[string]string, newConfig map[string]interface{}) *schema.ResourceData {
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	state := &terraform.InstanceState{
		ID:         id,
		Attributes: oldState,
	}
	rawConfig := terraform.NewResourceConfigRaw(newConfig)
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)
	return d
}

func TestFwmMemberGroup_Read_Grouped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	mockDescribeOrganMembers(patches, fwmClient, int64(1), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id": "mem-xxxx",
		"app_id":    "1234567890",
		"group_id":  "",
	})
	d.SetId("1234567890#mem-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890#mem-xxxx", d.Id())
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

func TestFwmMemberGroup_Read_Ungrouped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	mockDescribeOrganMembers(patches, fwmClient, int64(1), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", ""),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id": "mem-xxxx",
		"app_id":    "1234567890",
		"group_id":  "group-yyyy",
	})
	d.SetId("1234567890#mem-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890#mem-xxxx", d.Id())
	assert.Equal(t, "", d.Get("group_id"))
}

func TestFwmMemberGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	mockDescribeOrganMembers(patches, fwmClient, int64(0), []*fwmv20250611.OrganMemberItem{})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id": "mem-notexist",
		"app_id":    "1234567890",
		"group_id":  "",
	})
	d.SetId("1234567890#mem-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmMemberGroup_Read_NotFound_MultipleMembers(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	mockDescribeOrganMembers(patches, fwmClient, int64(2), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-aaaa", "1234567890", "group-yyyy"),
		buildOrganMemberItem("mem-bbbb", "1234567890", ""),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id": "mem-notexist",
		"app_id":    "1234567890",
		"group_id":  "",
	})
	d.SetId("1234567890#mem-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmMemberGroup_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	var capturedUserMemberId, capturedGroupId string
	patches.ApplyMethodFunc(fwmClient, "ModifyMemberGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyMemberGroupRequest) (*fwmv20250611.ModifyMemberGroupResponse, error) {
		if request.UserMemberId != nil {
			capturedUserMemberId = *request.UserMemberId
		}
		if request.GroupId != nil {
			capturedGroupId = *request.GroupId
		}
		resp := fwmv20250611.NewModifyMemberGroupResponse()
		resp.Response = &fwmv20250611.ModifyMemberGroupResponseParams{
			RequestId: ptrStringFwm("fake-request-id-create"),
		}
		return resp, nil
	})

	mockDescribeOrganMembers(patches, fwmClient, int64(1), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id": "mem-xxxx",
		"app_id":    "1234567890",
		"group_id":  "group-yyyy",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890#mem-xxxx", d.Id())
	assert.Equal(t, "mem-xxxx", capturedUserMemberId)
	assert.Equal(t, "group-yyyy", capturedGroupId)
}

func TestFwmMemberGroup_Update_MoveIntoGroup(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	var capturedUserMemberId, capturedGroupId string
	patches.ApplyMethodFunc(fwmClient, "ModifyMemberGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyMemberGroupRequest) (*fwmv20250611.ModifyMemberGroupResponse, error) {
		if request.UserMemberId != nil {
			capturedUserMemberId = *request.UserMemberId
		}
		if request.GroupId != nil {
			capturedGroupId = *request.GroupId
		}
		resp := fwmv20250611.NewModifyMemberGroupResponse()
		resp.Response = &fwmv20250611.ModifyMemberGroupResponseParams{
			RequestId: ptrStringFwm("fake-request-id-move-in"),
		}
		return resp, nil
	})

	mockDescribeOrganMembers(patches, fwmClient, int64(1), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	d := buildFwmMemberGroupData(t, "1234567890#mem-xxxx",
		map[string]string{
			"member_id": "mem-xxxx",
			"app_id":    "1234567890",
			"group_id":  "",
		},
		map[string]interface{}{
			"member_id": "mem-xxxx",
			"app_id":    "1234567890",
			"group_id":  "group-yyyy",
		},
	)

	err := fwm.ResourceTencentCloudFwmMemberGroup().Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", capturedUserMemberId)
	assert.Equal(t, "group-yyyy", capturedGroupId)
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

func TestFwmMemberGroup_Update_MoveOut(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	var capturedUserMemberId, capturedGroupId string
	patches.ApplyMethodFunc(fwmClient, "ModifyMemberGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyMemberGroupRequest) (*fwmv20250611.ModifyMemberGroupResponse, error) {
		if request.UserMemberId != nil {
			capturedUserMemberId = *request.UserMemberId
		}
		if request.GroupId != nil {
			capturedGroupId = *request.GroupId
		}
		resp := fwmv20250611.NewModifyMemberGroupResponse()
		resp.Response = &fwmv20250611.ModifyMemberGroupResponseParams{
			RequestId: ptrStringFwm("fake-request-id-move-out"),
		}
		return resp, nil
	})

	mockDescribeOrganMembers(patches, fwmClient, int64(1), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", ""),
	})

	meta := newMockMetaFwmMemberGroup()
	d := buildFwmMemberGroupData(t, "1234567890#mem-xxxx",
		map[string]string{
			"member_id": "mem-xxxx",
			"app_id":    "1234567890",
			"group_id":  "group-yyyy",
		},
		map[string]interface{}{
			"member_id": "mem-xxxx",
			"app_id":    "1234567890",
			"group_id":  "",
		},
	)

	err := fwm.ResourceTencentCloudFwmMemberGroup().Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", capturedUserMemberId)
	assert.Equal(t, "", capturedGroupId)
	assert.Equal(t, "", d.Get("group_id"))
}

func TestFwmMemberGroup_Update_SkipUnchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmClient", fwmClient)

	modifyCalled := false
	patches.ApplyMethodFunc(fwmClient, "ModifyMemberGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyMemberGroupRequest) (*fwmv20250611.ModifyMemberGroupResponse, error) {
		modifyCalled = true
		resp := fwmv20250611.NewModifyMemberGroupResponse()
		resp.Response = &fwmv20250611.ModifyMemberGroupResponseParams{
			RequestId: ptrStringFwm("fake-request-id-skip"),
		}
		return resp, nil
	})

	mockDescribeOrganMembers(patches, fwmClient, int64(1), []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	d := buildFwmMemberGroupData(t, "1234567890#mem-xxxx",
		map[string]string{
			"member_id": "mem-xxxx",
			"app_id":    "1234567890",
			"group_id":  "group-yyyy",
		},
		map[string]interface{}{
			"member_id": "mem-xxxx",
			"app_id":    "1234567890",
			"group_id":  "group-yyyy",
		},
	)

	err := fwm.ResourceTencentCloudFwmMemberGroup().Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

func TestFwmMemberGroup_Delete(t *testing.T) {
	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id": "mem-xxxx",
		"app_id":    "1234567890",
		"group_id":  "group-yyyy",
	})
	d.SetId("1234567890#mem-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
