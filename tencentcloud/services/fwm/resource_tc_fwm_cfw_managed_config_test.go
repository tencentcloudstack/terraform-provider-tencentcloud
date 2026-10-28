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

func ptrStrCfwManaged(s string) *string { return &s }

func ptrInt64CfwManaged(v int64) *int64 { return &v }

// mockCfwManagedClient wires both client accessors used by this resource to
// the same fake SDK client:
//   - UseFwmV20250611Client -> DescribeOrganMembers (Read, via DescribeFwmMemberGroup)
//   - UseFwmClient          -> ModifyCfwManaged   (Update)
func mockCfwManagedClient(patches *gomonkey.Patches, meta *mockMetaForFwmCfwManaged, client *fwmv20250611.Client) {
	patches.ApplyMethodReturn(meta.client, "UseFwmV20250611Client", client)
	patches.ApplyMethodReturn(meta.client, "UseFwmClient", client)
}

type modifyCfwManagedCapture struct {
	Called     bool
	MemberId   string
	CfwManaged int64
	// CfwManagedSet records whether request.CfwManaged was non-nil, which
	// matters for the "disable managed" (0) case.
	CfwManagedSet bool
}

func mockModifyCfwManaged(patches *gomonkey.Patches, client *fwmv20250611.Client, capture *modifyCfwManagedCapture) {
	patches.ApplyMethodFunc(client, "ModifyCfwManagedWithContext", func(_ context.Context, request *fwmv20250611.ModifyCfwManagedRequest) (*fwmv20250611.ModifyCfwManagedResponse, error) {
		capture.Called = true
		if request.MemberId != nil {
			capture.MemberId = *request.MemberId
		}
		if request.CfwManaged != nil {
			capture.CfwManagedSet = true
			capture.CfwManaged = *request.CfwManaged
		}

		resp := fwmv20250611.NewModifyCfwManagedResponse()
		resp.Response = &fwmv20250611.ModifyCfwManagedResponseParams{
			RequestId: ptrStrCfwManaged("fake-request-id"),
		}
		return resp, nil
	})
}

// buildCfwManagedMember builds the API item returned by DescribeOrganMembers.
func buildCfwManagedMember(memberId string, cfwManaged int64) *fwmv20250611.OrganMemberItem {
	return &fwmv20250611.OrganMemberItem{
		MemberId:   ptrStrCfwManaged(memberId),
		AppId:      ptrStrCfwManaged("1234567890"),
		CfwManaged: ptrInt64CfwManaged(cfwManaged),
	}
}

// mockDescribeOrganMembersForCfwManaged hooks the Context variant on purpose:
// the non-Context wrapper is inlined by the compiler, so patching it has no
// effect on darwin/arm64.
func mockDescribeOrganMembersForCfwManaged(patches *gomonkey.Patches, client *fwmv20250611.Client, members []*fwmv20250611.OrganMemberItem) {
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId: ptrStrCfwManaged("fake-request-id"),
			Members:   members,
		}
		return resp, nil
	})
}

// mockDescribeOrganMembersForCfwManagedPaged returns the page selected by
// request.Offset so the service pagination loop can be exercised.
func mockDescribeOrganMembersForCfwManagedPaged(patches *gomonkey.Patches, client *fwmv20250611.Client, pages [][]*fwmv20250611.OrganMemberItem) {
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
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId: ptrStrCfwManaged("fake-request-id"),
			Members:   page,
		}
		return resp, nil
	})
}

func TestFwmCfwManagedConfig_Create_Enable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)

	capture := &modifyCfwManagedCapture{}
	mockModifyCfwManaged(patches, client, capture)
	mockDescribeOrganMembersForCfwManaged(patches, client, []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-xxxx", 1),
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":   "mem-xxxx",
		"cfw_managed": 1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
	assert.True(t, capture.Called)
	assert.Equal(t, "mem-xxxx", capture.MemberId)
	assert.True(t, capture.CfwManagedSet)
	assert.Equal(t, int64(1), capture.CfwManaged)
	assert.Equal(t, 1, d.Get("cfw_managed"))
}

// TestFwmCfwManagedConfig_Create_Disable covers the zero value: cfw_managed=0
// must still be sent to the API (GetOkExists semantics), otherwise the switch
// can never be turned off.
func TestFwmCfwManagedConfig_Create_Disable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)

	capture := &modifyCfwManagedCapture{}
	mockModifyCfwManaged(patches, client, capture)
	mockDescribeOrganMembersForCfwManaged(patches, client, []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-xxxx", 0),
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":   "mem-xxxx",
		"cfw_managed": 0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
	assert.True(t, capture.Called)
	assert.True(t, capture.CfwManagedSet)
	assert.Equal(t, int64(0), capture.CfwManaged)
	assert.Equal(t, 0, d.Get("cfw_managed"))
}

func TestFwmCfwManagedConfig_Read_Exists(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)
	mockDescribeOrganMembersForCfwManaged(patches, client, []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-xxxx", 1),
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("mem-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
	assert.Equal(t, "mem-xxxx", d.Get("member_id"))
	assert.Equal(t, 1, d.Get("cfw_managed"))
}

func TestFwmCfwManagedConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)
	mockDescribeOrganMembersForCfwManaged(patches, client, []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-other", 1),
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("mem-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmCfwManagedConfig_Read_EmptyResponse verifies a nil API response is
// returned as an error instead of silently clearing the resource.
func TestFwmCfwManagedConfig_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return &fwmv20250611.DescribeOrganMembersResponse{}, nil
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("mem-xxxx")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "mem-xxxx", d.Id())
}

// TestFwmCfwManagedConfig_Read_Pagination verifies Read finds a member that
// lives on the second page of DescribeOrganMembers.
func TestFwmCfwManagedConfig_Read_Pagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	firstPage := make([]*fwmv20250611.OrganMemberItem, 0, 100)
	for i := 0; i < 100; i++ {
		firstPage = append(firstPage, buildCfwManagedMember("mem-other", 0))
	}
	secondPage := []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-target", 1),
	}

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)
	mockDescribeOrganMembersForCfwManagedPaged(patches, client, [][]*fwmv20250611.OrganMemberItem{firstPage, secondPage})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("mem-target")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mem-target", d.Id())
	assert.Equal(t, 1, d.Get("cfw_managed"))
}

func TestFwmCfwManagedConfig_Update_Change(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)

	capture := &modifyCfwManagedCapture{}
	mockModifyCfwManaged(patches, client, capture)
	mockDescribeOrganMembersForCfwManaged(patches, client, []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-xxxx", 1),
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":   "mem-xxxx",
		"cfw_managed": 1,
	})
	d.SetId("mem-xxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.Equal(t, "mem-xxxx", capture.MemberId)
	assert.Equal(t, int64(1), capture.CfwManaged)
	assert.Equal(t, 1, d.Get("cfw_managed"))
}

// TestFwmCfwManagedConfig_Update_AlwaysCallsModify locks in the current
// behavior: Update does not diff cfw_managed, it always pushes the value.
func TestFwmCfwManagedConfig_Update_AlwaysCallsModify(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)

	capture := &modifyCfwManagedCapture{}
	mockModifyCfwManaged(patches, client, capture)
	mockDescribeOrganMembersForCfwManaged(patches, client, []*fwmv20250611.OrganMemberItem{
		buildCfwManagedMember("mem-xxxx", 0),
	})

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":   "mem-xxxx",
		"cfw_managed": 0,
	})
	d.SetId("mem-xxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, capture.Called)
	assert.True(t, capture.CfwManagedSet)
	assert.Equal(t, int64(0), capture.CfwManaged)
}

// TestFwmCfwManagedConfig_Delete verifies delete only drops the resource from
// state: the member itself must not be removed from the organization.
func TestFwmCfwManagedConfig_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaForFwmCfwManaged()
	client := &fwmv20250611.Client{}
	mockCfwManagedClient(patches, meta, client)

	capture := &modifyCfwManagedCapture{}
	mockModifyCfwManaged(patches, client, capture)

	res := svcfwm.ResourceTencentCloudFwmCfwManagedConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"member_id":   "mem-xxxx",
		"cfw_managed": 1,
	})
	d.SetId("mem-xxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.False(t, capture.Called)
	assert.Equal(t, "mem-xxxx", d.Id())
}
