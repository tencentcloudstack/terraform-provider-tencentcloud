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

type mockMetaFwmShareRole struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmShareRole) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmShareRole{}

func newMockMetaFwmShareRole() *mockMetaFwmShareRole {
	return &mockMetaFwmShareRole{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwmShareRole(s string) *string {
	return &s
}

// mockShareRoleClient wires both client accessors used by this resource to the
// same fake SDK client:
//   - UseFwmClient          -> ModifyCfwShareRole (Create/Update)
//   - UseFwmV20250611Client -> DescribeOrganMembers (Read via DescribeFwmMemberGroup)
func mockShareRoleClient(patches *gomonkey.Patches, meta *mockMetaFwmShareRole, client *fwmv20250611.Client) {
	patches.ApplyMethodReturn(meta.client, "UseFwmClient", client)
	patches.ApplyMethodReturn(meta.client, "UseFwmV20250611Client", client)
}

func buildShareRoleMember(memberId, role string) *fwmv20250611.OrganMemberItem {
	return &fwmv20250611.OrganMemberItem{
		MemberId:     ptrStringFwmShareRole(memberId),
		AppId:        ptrStringFwmShareRole("app-456"),
		CfwShareRole: ptrStringFwmShareRole(role),
	}
}

// mockDescribeOrganMembersForShareRole hooks the Context variant on purpose: the
// non-Context wrapper is inlined, so patching it has no effect on darwin/arm64.
func mockDescribeOrganMembersForShareRole(patches *gomonkey.Patches, client *fwmv20250611.Client, members []*fwmv20250611.OrganMemberItem) {
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		total := int64(len(members))
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId:  ptrStringFwmShareRole("fake-request-id"),
			TotalCount: &total,
			Members:    members,
		}
		return resp, nil
	})
}

func mockDescribeOrganMembersForShareRolePaged(patches *gomonkey.Patches, client *fwmv20250611.Client, pages [][]*fwmv20250611.OrganMemberItem) {
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		var (
			offset int64
			limit  int64 = 100
		)
		if request != nil && request.Offset != nil {
			offset = *request.Offset
		}
		if request != nil && request.Limit != nil {
			limit = *request.Limit
		}

		var page []*fwmv20250611.OrganMemberItem
		if limit > 0 && offset%limit == 0 {
			if idx := int(offset / limit); idx >= 0 && idx < len(pages) {
				page = pages[idx]
			}
		}

		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		total := int64(len(page))
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId:  ptrStringFwmShareRole("fake-request-id"),
			TotalCount: &total,
			Members:    page,
		}
		return resp, nil
	})
}

type modifyShareRoleCapture struct {
	Called         bool
	UserMemberId   string
	CfwShareRole   string
	SharerMemberId string
	SharerSet      bool
}

func mockModifyCfwShareRole(patches *gomonkey.Patches, client *fwmv20250611.Client, capture *modifyShareRoleCapture) {
	patches.ApplyMethodFunc(client, "ModifyCfwShareRoleWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwShareRoleRequest) (*fwmv20250611.ModifyCfwShareRoleResponse, error) {
		capture.Called = true
		if request.UserMemberId != nil {
			capture.UserMemberId = *request.UserMemberId
		}
		if request.CfwShareRole != nil {
			capture.CfwShareRole = *request.CfwShareRole
		}
		if request.SharerMemberId != nil {
			capture.SharerSet = true
			capture.SharerMemberId = *request.SharerMemberId
		}

		resp := fwmv20250611.NewModifyCfwShareRoleResponse()
		resp.Response = &fwmv20250611.ModifyCfwShareRoleResponseParams{
			RequestId: ptrStringFwmShareRole("fake-request-id"),
		}
		return resp, nil
	})
}

func TestFwmCfwShareRoleConfig_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()

	assert.Len(t, res.Schema, 3)

	role := res.Schema["cfw_share_role"]
	assert.NotNil(t, role)
	assert.True(t, role.Required)
	assert.False(t, role.ForceNew)

	userMemberId := res.Schema["user_member_id"]
	assert.NotNil(t, userMemberId)
	assert.True(t, userMemberId.Required)
	assert.True(t, userMemberId.ForceNew)

	sharerMemberId := res.Schema["sharer_member_id"]
	assert.NotNil(t, sharerMemberId)
	assert.True(t, sharerMemberId.Optional)
	assert.False(t, sharerMemberId.Required)

	// Attributes of the previous generation are gone.
	for _, key := range []string{"member_id", "app_id"} {
		_, ok := res.Schema[key]
		assert.False(t, ok, "schema %s should not exist", key)
	}

	assert.NotNil(t, res.Importer)
}

func TestFwmCfwShareRoleConfig_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)

	capture := &modifyShareRoleCapture{}
	mockModifyCfwShareRole(patches, fwmClient, capture)
	mockDescribeOrganMembersForShareRole(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-123", "none"),
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-123",
		"cfw_share_role": "none",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "member-123", d.Id())

	assert.True(t, capture.Called)
	assert.Equal(t, "member-123", capture.UserMemberId)
	assert.Equal(t, "none", capture.CfwShareRole)
	assert.False(t, capture.SharerSet)

	// Read back
	assert.Equal(t, "member-123", d.Get("user_member_id"))
	assert.Equal(t, "none", d.Get("cfw_share_role"))
}

// TestFwmCfwShareRoleConfig_Create_WithSharer verifies sharer_member_id is
// forwarded when cfw_share_role is "user".
func TestFwmCfwShareRoleConfig_Create_WithSharer(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)

	capture := &modifyShareRoleCapture{}
	mockModifyCfwShareRole(patches, fwmClient, capture)
	mockDescribeOrganMembersForShareRole(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-123", "user"),
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id":   "member-123",
		"cfw_share_role":   "user",
		"sharer_member_id": "app-789",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "member-123", d.Id())
	assert.True(t, capture.Called)
	assert.Equal(t, "user", capture.CfwShareRole)
	assert.True(t, capture.SharerSet)
	assert.Equal(t, "app-789", capture.SharerMemberId)
}

func TestFwmCfwShareRoleConfig_Read_MemberExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)
	mockDescribeOrganMembersForShareRole(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-123", "user"),
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-123",
		"cfw_share_role": "none",
	})
	d.SetId("member-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "member-123", d.Id())
	assert.Equal(t, "member-123", d.Get("user_member_id"))
	assert.Equal(t, "user", d.Get("cfw_share_role"))
}

func TestFwmCfwShareRoleConfig_Read_MemberNotExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)
	mockDescribeOrganMembersForShareRole(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-other", "none"),
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-123",
		"cfw_share_role": "none",
	})
	d.SetId("member-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmCfwShareRoleConfig_Read_EmptyResponse verifies a nil API response is
// surfaced as an error instead of silently clearing the resource.
func TestFwmCfwShareRoleConfig_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return &fwmv20250611.DescribeOrganMembersResponse{}, nil
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-123",
		"cfw_share_role": "none",
	})
	d.SetId("member-123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "member-123", d.Id())
}

// TestFwmCfwShareRoleConfig_Read_Pagination verifies Read can find the member on
// the second page returned by DescribeOrganMembers.
func TestFwmCfwShareRoleConfig_Read_Pagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	firstPage := make([]*fwmv20250611.OrganMemberItem, 0, 100)
	for i := 0; i < 100; i++ {
		firstPage = append(firstPage, buildShareRoleMember("member-other", "none"))
	}
	secondPage := []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-target", "sharer"),
	}

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)
	mockDescribeOrganMembersForShareRolePaged(patches, fwmClient, [][]*fwmv20250611.OrganMemberItem{firstPage, secondPage})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-target",
		"cfw_share_role": "none",
	})
	d.SetId("member-target")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "member-target", d.Id())
	assert.Equal(t, "sharer", d.Get("cfw_share_role"))
}

func TestFwmCfwShareRoleConfig_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)

	capture := &modifyShareRoleCapture{}
	mockModifyCfwShareRole(patches, fwmClient, capture)
	mockDescribeOrganMembersForShareRole(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-123", "user"),
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id":   "member-123",
		"cfw_share_role":   "user",
		"sharer_member_id": "app-789",
	})
	d.SetId("member-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.Equal(t, "member-123", capture.UserMemberId)
	assert.Equal(t, "user", capture.CfwShareRole)
	assert.Equal(t, "app-789", capture.SharerMemberId)
	assert.Equal(t, "user", d.Get("cfw_share_role"))
}

// TestFwmCfwShareRoleConfig_Update_AlwaysCallsModify locks in the current
// behavior: Update does not diff cfw_share_role, it always pushes the value.
func TestFwmCfwShareRoleConfig_Update_AlwaysCallsModify(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)

	capture := &modifyShareRoleCapture{}
	mockModifyCfwShareRole(patches, fwmClient, capture)
	mockDescribeOrganMembersForShareRole(patches, fwmClient, []*fwmv20250611.OrganMemberItem{
		buildShareRoleMember("member-123", "none"),
	})

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-123",
		"cfw_share_role": "none",
	})
	d.SetId("member-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.Equal(t, "none", capture.CfwShareRole)
}

// TestFwmCfwShareRoleConfig_Delete verifies delete only drops the resource from
// state: the share role itself must not be reset.
func TestFwmCfwShareRoleConfig_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	meta := newMockMetaFwmShareRole()
	mockShareRoleClient(patches, meta, fwmClient)

	capture := &modifyShareRoleCapture{}
	mockModifyCfwShareRole(patches, fwmClient, capture)

	res := fwm.ResourceTencentCloudFwmCfwShareRoleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_member_id": "member-123",
		"cfw_share_role": "none",
	})
	d.SetId("member-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.False(t, capture.Called)
	assert.Equal(t, "member-123", d.Id())
}
