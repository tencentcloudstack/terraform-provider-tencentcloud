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

type mockMetaFwmAccountGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmAccountGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmAccountGroup{}

func newMockMetaFwmAccountGroup() *mockMetaFwmAccountGroup {
	return &mockMetaFwmAccountGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFAG(s string) *string {
	return &s
}

func TestFwmAccountGroup_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, "test account group", d.Get("description"))
	assert.Equal(t, 0, d.Get("member_count"))
	assert.Equal(t, true, d.Get("can_delete"))
	assert.Equal(t, "2026-10-22 10:00:00", d.Get("create_time"))
}

func TestFwmAccountGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(0)
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:     ptrStringFAG("fake-request-id"),
			TotalCount:    &total,
			AccountGroups: []*fwmv20250611.AccountGroup{},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmAccountGroup_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateAccountGroupRequest) (*fwmv20250611.CreateAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "test-account-group", *request.GroupName)
		assert.NotNil(t, request.Description)
		assert.Equal(t, "test account group", *request.Description)

		resp := fwmv20250611.NewCreateAccountGroupResponse()
		resp.Response = &fwmv20250611.CreateAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			GroupId:   ptrStringFAG("acg-1a2b3c4d"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group",
		"description": "test account group",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "acg-1a2b3c4d", d.Id())
	assert.Equal(t, "test-account-group", d.Get("group_name"))
	assert.Equal(t, "test account group", d.Get("description"))
}

func TestFwmAccountGroup_Create_EmptyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateAccountGroupRequest) (*fwmv20250611.CreateAccountGroupResponse, error) {
		resp := fwmv20250611.NewCreateAccountGroupResponse()
		resp.Response = &fwmv20250611.CreateAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-create"),
			GroupId:   nil,
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group",
		"description": "test account group",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestFwmAccountGroup_Update_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyAccountGroupRequest) (*fwmv20250611.ModifyAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)
		assert.NotNil(t, request.GroupName)
		assert.Equal(t, "test-account-group-updated", *request.GroupName)

		resp := fwmv20250611.NewModifyAccountGroupResponse()
		resp.Response = &fwmv20250611.ModifyAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeAccountGroupsWithContext", func(_ context.Context, request *fwmv20250611.DescribeAccountGroupsRequest) (*fwmv20250611.DescribeAccountGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeAccountGroupsResponse()
		total := int64(1)
		memberCount := int64(0)
		canDelete := true
		resp.Response = &fwmv20250611.DescribeAccountGroupsResponseParams{
			RequestId:  ptrStringFAG("fake-request-id"),
			TotalCount: &total,
			AccountGroups: []*fwmv20250611.AccountGroup{
				{
					GroupId:     ptrStringFAG("acg-1a2b3c4d"),
					GroupName:   ptrStringFAG("test-account-group-updated"),
					Description: ptrStringFAG("test account group"),
					MemberCount: &memberCount,
					CanDelete:   &canDelete,
					CreateTime:  ptrStringFAG("2026-10-22 10:00:00"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name":  "test-account-group-updated",
		"description": "test account group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-account-group-updated", d.Get("group_name"))
}

func TestFwmAccountGroup_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteAccountGroupRequest) (*fwmv20250611.DeleteAccountGroupResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "acg-1a2b3c4d", *request.GroupId)

		resp := fwmv20250611.NewDeleteAccountGroupResponse()
		resp.Response = &fwmv20250611.DeleteAccountGroupResponseParams{
			RequestId: ptrStringFAG("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestFwmAccountGroup_Delete_ErrorPropagated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmAccountGroup().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteAccountGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteAccountGroupRequest) (*fwmv20250611.DeleteAccountGroupResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaFwmAccountGroup()
	res := fwm.ResourceTencentCloudFwmAccountGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_name": "test-account-group",
	})
	d.SetId("acg-1a2b3c4d")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}
