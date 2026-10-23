package fwm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcfwm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

type mockMetaForFwmCfwManaged struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForFwmCfwManaged) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForFwmCfwManaged{}

func newMockMetaForFwmCfwManaged() *mockMetaForFwmCfwManaged {
	return &mockMetaForFwmCfwManaged{client: &connectivity.TencentCloudClient{}}
}

func ptrStrFCM(s string) *string { return &s }

// TestFwmCfwManaged_Create verifies that Create sets the composite id and calls ModifyCfwManaged.
func TestFwmCfwManaged_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	var capturedMemberId *string
	var capturedCfwManaged *int64
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwManaged", func(request *fwm.ModifyCfwManagedRequest) (*fwm.ModifyCfwManagedResponse, error) {
		capturedMemberId = request.MemberId
		capturedCfwManaged = request.CfwManaged
		resp := fwm.NewModifyCfwManagedResponse()
		resp.Response = &fwm.ModifyCfwManagedResponseParams{
			RequestId: ptrStrFCM("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembers", func(request *fwm.DescribeOrganMembersRequest) (*fwm.DescribeOrganMembersResponse, error) {
		resp := fwm.NewDescribeOrganMembersResponse()
		resp.Response = &fwm.DescribeOrganMembersResponseParams{
			Members: []*fwm.OrganMemberItem{
				{
					MemberId:   ptrStrFCM("100012345678"),
					AppId:      ptrStrFCM("1234567890"),
					CfwManaged: func() *int64 { v := int64(0); return &v }(),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"app_id":      "1234567890",
		"member_id":   "100012345678",
		"cfw_managed": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890#100012345678", d.Id())

	// Verify ModifyCfwManaged was called with the right values.
	assert.NotNil(t, capturedMemberId)
	assert.Equal(t, "100012345678", *capturedMemberId)
	assert.NotNil(t, capturedCfwManaged)
	assert.Equal(t, int64(1), *capturedCfwManaged)
}

// TestFwmCfwManaged_Read_MemberExists verifies that Read parses the composite id and sets fields.
func TestFwmCfwManaged_Read_MemberExists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembers", func(request *fwm.DescribeOrganMembersRequest) (*fwm.DescribeOrganMembersResponse, error) {
		resp := fwm.NewDescribeOrganMembersResponse()
		resp.Response = &fwm.DescribeOrganMembersResponseParams{
			Members: []*fwm.OrganMemberItem{
				{
					MemberId:   ptrStrFCM("100012345678"),
					AppId:      ptrStrFCM("1234567890"),
					CfwManaged: func() *int64 { v := int64(1); return &v }(),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("1234567890#100012345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890#100012345678", d.Id())
	assert.Equal(t, "100012345678", d.Get("member_id"))
	assert.Equal(t, "1234567890", d.Get("app_id"))
	assert.Equal(t, 1, d.Get("cfw_managed"))
}

// TestFwmCfwManaged_Read_MemberNotFound verifies that Read clears state when the member is not matched.
func TestFwmCfwManaged_Read_MemberNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembers", func(request *fwm.DescribeOrganMembersRequest) (*fwm.DescribeOrganMembersResponse, error) {
		resp := fwm.NewDescribeOrganMembersResponse()
		resp.Response = &fwm.DescribeOrganMembersResponseParams{
			Members: []*fwm.OrganMemberItem{},
		}
		return resp, nil
	})

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("1234567890#100012345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmCfwManaged_Read_BrokenId verifies that Read returns an error when the id is broken.
func TestFwmCfwManaged_Read_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("brokenid")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestFwmCfwManaged_Update_Changed verifies that Update calls ModifyCfwManaged when cfw_managed changes.
func TestFwmCfwManaged_Update_Changed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	called := false
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwManaged", func(request *fwm.ModifyCfwManagedRequest) (*fwm.ModifyCfwManagedResponse, error) {
		called = true
		resp := fwm.NewModifyCfwManagedResponse()
		resp.Response = &fwm.ModifyCfwManagedResponseParams{
			RequestId: ptrStrFCM("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembers", func(request *fwm.DescribeOrganMembersRequest) (*fwm.DescribeOrganMembersResponse, error) {
		resp := fwm.NewDescribeOrganMembersResponse()
		resp.Response = &fwm.DescribeOrganMembersResponseParams{
			Members: []*fwm.OrganMemberItem{
				{
					MemberId:   ptrStrFCM("100012345678"),
					AppId:      ptrStrFCM("1234567890"),
					CfwManaged: func() *int64 { v := int64(0); return &v }(),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"app_id":      "1234567890",
		"member_id":   "100012345678",
		"cfw_managed": 1,
	})
	d.SetId("1234567890#100012345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, called, "ModifyCfwManaged should be called when cfw_managed changes")
}

// TestFwmCfwManaged_Update_Unchanged verifies that Update skips ModifyCfwManaged when cfw_managed is unchanged.
func TestFwmCfwManaged_Update_Unchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	called := false
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwManaged", func(request *fwm.ModifyCfwManagedRequest) (*fwm.ModifyCfwManagedResponse, error) {
		called = true
		resp := fwm.NewModifyCfwManagedResponse()
		resp.Response = &fwm.ModifyCfwManagedResponseParams{
			RequestId: ptrStrFCM("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeOrganMembers", func(request *fwm.DescribeOrganMembersRequest) (*fwm.DescribeOrganMembersResponse, error) {
		resp := fwm.NewDescribeOrganMembersResponse()
		resp.Response = &fwm.DescribeOrganMembersResponseParams{
			Members: []*fwm.OrganMemberItem{
				{
					MemberId:   ptrStrFCM("100012345678"),
					AppId:      ptrStrFCM("1234567890"),
					CfwManaged: func() *int64 { v := int64(1); return &v }(),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"app_id":      "1234567890",
		"member_id":   "100012345678",
		"cfw_managed": 1,
	})
	d.SetId("1234567890#100012345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, called, "ModifyCfwManaged should not be called when cfw_managed is unchanged")
}

// TestFwmCfwManaged_Delete verifies that Delete returns nil without calling any cloud api.
func TestFwmCfwManaged_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwm.Client{}
	patches.ApplyMethodReturn(newMockMetaForFwmCfwManaged().client, "UseFwmClient", fwmClient)

	called := false
	patches.ApplyMethodFunc(fwmClient, "ModifyCfwManaged", func(request *fwm.ModifyCfwManagedRequest) (*fwm.ModifyCfwManagedResponse, error) {
		called = true
		return fwm.NewModifyCfwManagedResponse(), nil
	})

	meta := newMockMetaForFwmCfwManaged()
	res := svcfwm.ResourceTencentCloudFwmCfwManaged()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("1234567890#100012345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.False(t, called, "Delete should not call ModifyCfwManaged")
}
