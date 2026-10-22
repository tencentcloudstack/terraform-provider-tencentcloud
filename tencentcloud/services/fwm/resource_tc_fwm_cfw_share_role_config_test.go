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

func ptrStringFwm(s string) *string {
	return &s
}

func buildOrganMembersResponse(role, sharerAppId string) *fwmv20250611.DescribeOrganMembersResponse {
	resp := fwmv20250611.NewDescribeOrganMembersResponse()
	total := int64(1)
	resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
		RequestId:  ptrStringFwm("fake-request-id"),
		TotalCount: &total,
		Members: []*fwmv20250611.OrganMemberItem{
			{
				MemberId:       ptrStringFwm("member-123"),
				AppId:          ptrStringFwm("app-456"),
				CfwShareRole:   ptrStringFwm(role),
				CfwSharerAppId: ptrStringFwm(sharerAppId),
			},
		},
	}
	return resp
}

func buildEmptyOrganMembersResponse() *fwmv20250611.DescribeOrganMembersResponse {
	resp := fwmv20250611.NewDescribeOrganMembersResponse()
	total := int64(0)
	resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
		RequestId:  ptrStringFwm("fake-request-id"),
		TotalCount: &total,
		Members:    []*fwmv20250611.OrganMemberItem{},
	}
	return resp
}

func TestFwmCfwShareRole_Read_MemberExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildOrganMembersResponse("none", ""), nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "none",
	})
	d.SetId("app-456#member-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "app-456#member-123", d.Id())
	assert.Equal(t, "none", d.Get("cfw_share_role"))
	assert.Equal(t, "member-123", d.Get("member_id"))
	assert.Equal(t, "app-456", d.Get("app_id"))
}

func TestFwmCfwShareRole_Read_MemberNotExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildEmptyOrganMembersResponse(), nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "none",
	})
	d.SetId("app-456#member-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmCfwShareRole_Create_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildOrganMembersResponse("none", ""), nil
	})

	patches.ApplyMethodFunc(fwmClient, "ModifyCfwShareRoleWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwShareRoleRequest) (*fwmv20250611.ModifyCfwShareRoleResponse, error) {
		resp := fwmv20250611.NewModifyCfwShareRoleResponse()
		resp.Response = &fwmv20250611.ModifyCfwShareRoleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-create"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "none",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "app-456#member-123", d.Id())
}

func TestFwmCfwShareRole_Update_DirectChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	callCount := 0
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwShareRoleWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwShareRoleRequest) (*fwmv20250611.ModifyCfwShareRoleResponse, error) {
		callCount++
		assert.NotNil(t, request.CfwShareRole)
		assert.Equal(t, "user", *request.CfwShareRole)
		assert.NotNil(t, request.UserMemberId)
		assert.Equal(t, "member-123", *request.UserMemberId)

		resp := fwmv20250611.NewModifyCfwShareRoleResponse()
		resp.Response = &fwmv20250611.ModifyCfwShareRoleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildOrganMembersResponse("none", ""), nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "user",
	})
	d.SetId("app-456#member-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestFwmCfwShareRole_Update_SharerToUser_Relay(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	roles := []string{}
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwShareRoleWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwShareRoleRequest) (*fwmv20250611.ModifyCfwShareRoleResponse, error) {
		roles = append(roles, *request.CfwShareRole)
		resp := fwmv20250611.NewModifyCfwShareRoleResponse()
		resp.Response = &fwmv20250611.ModifyCfwShareRoleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-relay"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildOrganMembersResponse("sharer", ""), nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "user",
	})
	d.SetId("app-456#member-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(roles))
	assert.Equal(t, "none", roles[0])
	assert.Equal(t, "user", roles[1])
}

func TestFwmCfwShareRole_Update_UserToSharer_Relay(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	roles := []string{}
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwShareRoleWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwShareRoleRequest) (*fwmv20250611.ModifyCfwShareRoleResponse, error) {
		roles = append(roles, *request.CfwShareRole)
		resp := fwmv20250611.NewModifyCfwShareRoleResponse()
		resp.Response = &fwmv20250611.ModifyCfwShareRoleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-relay"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildOrganMembersResponse("user", "app-789"), nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "sharer",
	})
	d.SetId("app-456#member-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(roles))
	assert.Equal(t, "none", roles[0])
	assert.Equal(t, "sharer", roles[1])
}

func TestFwmCfwShareRole_Update_SameRole_SkipWrite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmShareRole().client, "UseFwmClient", fwmClient)

	callCount := 0
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwShareRoleWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwShareRoleRequest) (*fwmv20250611.ModifyCfwShareRoleResponse, error) {
		callCount++
		resp := fwmv20250611.NewModifyCfwShareRoleResponse()
		resp.Response = &fwmv20250611.ModifyCfwShareRoleResponseParams{
			RequestId: ptrStringFwm("fake-request-id-skip"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return buildOrganMembersResponse("none", ""), nil
	})

	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "none",
	})
	d.SetId("app-456#member-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, callCount)
}

func TestFwmCfwShareRole_Delete_Empty(t *testing.T) {
	meta := newMockMetaFwmShareRole()
	res := fwm.ResourceTencentCloudFwmCfwShareRole()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":      "member-123",
		"app_id":         "app-456",
		"cfw_share_role": "none",
	})
	d.SetId("app-456#member-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "app-456#member-123", d.Id())
}
