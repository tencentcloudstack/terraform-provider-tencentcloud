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

// buildOrganMemberItem builds the API item returned by DescribeOrganMembers.
// groupId is left empty for a member that does not belong to any group.
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

// mockDescribeOrganMembers makes DescribeOrganMembers return a single page.
// The mock hooks DescribeOrganMembersWithContext: the non-Context variant
// simply forwards to it, and the latter is not inlined, so the patch takes
// effect on darwin/arm64 without disabling inlining for the whole build.
func mockDescribeOrganMembers(patches *gomonkey.Patches, client *fwmv20250611.Client, members []*fwmv20250611.OrganMemberItem) {
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId: ptrStringFwm("fake-request-id"),
			Members:   members,
		}
		return resp, nil
	})
}

// mockDescribeOrganMembersPaged makes DescribeOrganMembers return the page
// selected by request.Offset, which lets tests cover the pagination loop.
func mockDescribeOrganMembersPaged(patches *gomonkey.Patches, client *fwmv20250611.Client, pages [][]*fwmv20250611.OrganMemberItem) {
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		var offset int64
		if request != nil && request.Offset != nil {
			offset = *request.Offset
		}

		var (
			limit int64 = 100
			page  []*fwmv20250611.OrganMemberItem
		)
		if request != nil && request.Limit != nil {
			limit = *request.Limit
		}
		if limit > 0 && offset%limit == 0 {
			idx := int(offset / limit)
			if idx >= 0 && idx < len(pages) {
				page = pages[idx]
			}
		}

		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId: ptrStringFwm("fake-request-id"),
			Members:   page,
		}
		return resp, nil
	})
}

type modifyMemberGroupCapture struct {
	Called       bool
	UserMemberId string
	GroupId      string
}

func mockModifyMemberGroup(patches *gomonkey.Patches, client *fwmv20250611.Client, capture *modifyMemberGroupCapture) {
	patches.ApplyMethodFunc(client, "ModifyMemberGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyMemberGroupRequest) (*fwmv20250611.ModifyMemberGroupResponse, error) {
		capture.Called = true
		if request.UserMemberId != nil {
			capture.UserMemberId = *request.UserMemberId
		}
		if request.GroupId != nil {
			capture.GroupId = *request.GroupId
		}

		resp := fwmv20250611.NewModifyMemberGroupResponse()
		resp.Response = &fwmv20250611.ModifyMemberGroupResponseParams{
			RequestId: ptrStringFwm("fake-request-id-modify"),
		}
		return resp, nil
	})
}

// buildFwmMemberGroupConfigData builds a *schema.ResourceData with the given
// old state (key->string) and new raw config (key->interface{}) so that
// d.HasChange returns true for fields that differ between old and new.
func buildFwmMemberGroupConfigData(t *testing.T, id string, oldState map[string]string, newConfig map[string]interface{}) *schema.ResourceData {
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
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

func TestFwmMemberGroupConfig_Read_Grouped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-xxxx",
		"group_id":       "",
	})
	d.SetId("mem-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
	assert.Equal(t, "mem-xxxx", d.Get("user_member_id"))
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

func TestFwmMemberGroupConfig_Read_Ungrouped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", ""),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-xxxx",
		"group_id":       "group-yyyy",
	})
	d.SetId("mem-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
	assert.Equal(t, "", d.Get("group_id"))
}

func TestFwmMemberGroupConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-aaaa", "1234567890", "group-yyyy"),
		buildOrganMemberItem("mem-bbbb", "1234567890", ""),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-notexist",
		"group_id":       "",
	})
	d.SetId("mem-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmMemberGroupConfig_Read_Pagination verifies the member is found on the
// second page: the first page is full (100 items) so the service must advance
// the offset instead of stopping.
func TestFwmMemberGroupConfig_Read_Pagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	firstPage := make([]*fwmv20250611.OrganMemberItem, 0, 100)
	for i := 0; i < 100; i++ {
		firstPage = append(firstPage, buildOrganMemberItem("mem-other", "1234567890", ""))
	}
	secondPage := []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-target", "1234567890", "group-yyyy"),
	}

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeOrganMembersPaged(patches, fwmClient, [][]*fwmv20250611.OrganMemberItem{firstPage, secondPage})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-target",
		"group_id":       "",
	})
	d.SetId("mem-target")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-target", d.Id())
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

// TestFwmMemberGroupConfig_Read_EmptyResponse verifies a nil API response is
// surfaced as an error instead of silently clearing the resource.
func TestFwmMemberGroupConfig_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return &fwmv20250611.DescribeOrganMembersResponse{}, nil
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-xxxx",
		"group_id":       "",
	})
	d.SetId("mem-xxxx")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
}

func TestFwmMemberGroupConfig_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)

	capture := &modifyMemberGroupCapture{}
	mockModifyMemberGroup(patches, fwmClient, capture)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-xxxx",
		"group_id":       "group-yyyy",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.Equal(t, "mem-xxxx", d.Id())
	assert.Equal(t, "mem-xxxx", capture.UserMemberId)
	assert.Equal(t, "group-yyyy", capture.GroupId)
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

func TestFwmMemberGroupConfig_Update_MoveIntoGroup(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)

	capture := &modifyMemberGroupCapture{}
	mockModifyMemberGroup(patches, fwmClient, capture)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	d := buildFwmMemberGroupConfigData(t, "mem-xxxx",
		map[string]string{
			"user_member_id": "mem-xxxx",
			"group_id":       "",
		},
		map[string]interface{}{
			"user_member_id": "mem-xxxx",
			"group_id":       "group-yyyy",
		},
	)

	err := fwm.ResourceTencentCloudFwmMemberGroupConfig().Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.Equal(t, "mem-xxxx", capture.UserMemberId)
	assert.Equal(t, "group-yyyy", capture.GroupId)
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

func TestFwmMemberGroupConfig_Update_MoveOut(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)

	capture := &modifyMemberGroupCapture{}
	mockModifyMemberGroup(patches, fwmClient, capture)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", ""),
	})

	meta := newMockMetaFwmMemberGroup()
	d := buildFwmMemberGroupConfigData(t, "mem-xxxx",
		map[string]string{
			"user_member_id": "mem-xxxx",
			"group_id":       "group-yyyy",
		},
		map[string]interface{}{
			"user_member_id": "mem-xxxx",
			"group_id":       "",
		},
	)

	err := fwm.ResourceTencentCloudFwmMemberGroupConfig().Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.Equal(t, "mem-xxxx", capture.UserMemberId)
	assert.Equal(t, "", capture.GroupId)
	assert.Equal(t, "", d.Get("group_id"))
}

func TestFwmMemberGroupConfig_Update_SkipUnchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)

	capture := &modifyMemberGroupCapture{}
	mockModifyMemberGroup(patches, fwmClient, capture)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	d := buildFwmMemberGroupConfigData(t, "mem-xxxx",
		map[string]string{
			"user_member_id": "mem-xxxx",
			"group_id":       "group-yyyy",
		},
		map[string]interface{}{
			"user_member_id": "mem-xxxx",
			"group_id":       "group-yyyy",
		},
	)

	err := fwm.ResourceTencentCloudFwmMemberGroupConfig().Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, capture.Called)
	assert.Equal(t, "group-yyyy", d.Get("group_id"))
}

// TestFwmMemberGroupConfig_Delete verifies delete only drops the resource from
// state: the member itself must not be removed from the organization.
func TestFwmMemberGroupConfig_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmMemberGroup().client, "UseFwmV20250611Client", fwmClient)

	capture := &modifyMemberGroupCapture{}
	mockModifyMemberGroup(patches, fwmClient, capture)
	mockDescribeOrganMembers(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildOrganMemberItem("mem-xxxx", "1234567890", "group-yyyy"),
	})

	meta := newMockMetaFwmMemberGroup()
	res := fwm.ResourceTencentCloudFwmMemberGroupConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "mem-xxxx",
		"group_id":       "group-yyyy",
	})
	d.SetId("mem-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.False(t, capture.Called)
	assert.Equal(t, "mem-xxxx", d.Id())
}
