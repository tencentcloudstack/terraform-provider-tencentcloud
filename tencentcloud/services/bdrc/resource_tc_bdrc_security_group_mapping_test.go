package bdrc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcbdrc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"
)

type mockMetaBdrcSecurityGroupMapping struct {
	client *connectivity.TencentCloudClient
}

var _ tccommon.ProviderMeta = &mockMetaBdrcSecurityGroupMapping{}

func (m *mockMetaBdrcSecurityGroupMapping) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

func newMockMetaBdrcSecurityGroupMapping() *mockMetaBdrcSecurityGroupMapping {
	return &mockMetaBdrcSecurityGroupMapping{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBdrc(s string) *string {
	return &s
}

func mockDescribeSecurityGroupMappingsResponse(sitePairId, mappingId, srcSgId, targetSgId, lifeState string) *bdrcv20260330.DescribeSecurityGroupMappingsResponse {
	resp := bdrcv20260330.NewDescribeSecurityGroupMappingsResponse()
	resp.Response = &bdrcv20260330.DescribeSecurityGroupMappingsResponseParams{
		SecurityGroupMappingSet: []*bdrcv20260330.SecurityGroupMapping{
			{
				SecurityGroupMappingId: ptrStringBdrc(mappingId),
				SitePairId:             ptrStringBdrc(sitePairId),
				SourceSecurityGroupId:  ptrStringBdrc(srcSgId),
				TargetSecurityGroupId:  ptrStringBdrc(targetSgId),
				LifeState:              ptrStringBdrc(lifeState),
			},
		},
		RequestId: ptrStringBdrc("fake-request-id"),
	}
	return resp
}

func mockEmptyDescribeSecurityGroupMappingsResponse() *bdrcv20260330.DescribeSecurityGroupMappingsResponse {
	resp := bdrcv20260330.NewDescribeSecurityGroupMappingsResponse()
	resp.Response = &bdrcv20260330.DescribeSecurityGroupMappingsResponseParams{
		SecurityGroupMappingSet: []*bdrcv20260330.SecurityGroupMapping{},
		RequestId:               ptrStringBdrc("fake-request-id"),
	}
	return resp
}

// TestBdrcSecurityGroupMapping_Create verifies Create passes correct params and sets the combined id.
func TestBdrcSecurityGroupMapping_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSecurityGroupMapping().client, "UseBdrcV20260330Client", bdrcClient)

	var capturedRequest *bdrcv20260330.CreateSecurityGroupMappingRequest
	patches.ApplyMethodFunc(bdrcClient, "CreateSecurityGroupMappingWithContext", func(ctx context.Context, request *bdrcv20260330.CreateSecurityGroupMappingRequest) (*bdrcv20260330.CreateSecurityGroupMappingResponse, error) {
		capturedRequest = request
		resp := bdrcv20260330.NewCreateSecurityGroupMappingResponse()
		resp.Response = &bdrcv20260330.CreateSecurityGroupMappingResponseParams{
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bdrcClient, "DescribeSecurityGroupMappingsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeSecurityGroupMappingsRequest) (*bdrcv20260330.DescribeSecurityGroupMappingsResponse, error) {
		return mockDescribeSecurityGroupMappingsResponse("sp-test", "sgm-001", "sg-src", "sg-target", "NORMAL"), nil
	})

	meta := newMockMetaBdrcSecurityGroupMapping()
	res := svcbdrc.ResourceTencentCloudBdrcSecurityGroupMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sp-test",
		"src_security_group_id":    "sg-src",
		"target_security_group_id": "sg-target",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sp-test#sgm-001", d.Id())

	assert.NotNil(t, capturedRequest.SitePairId)
	assert.Equal(t, "sp-test", *capturedRequest.SitePairId)
	assert.NotNil(t, capturedRequest.SrcSecurityGroupId)
	assert.Equal(t, "sg-src", *capturedRequest.SrcSecurityGroupId)
	assert.NotNil(t, capturedRequest.TargetSecurityGroupId)
	assert.Equal(t, "sg-target", *capturedRequest.TargetSecurityGroupId)

	assert.Equal(t, "sgm-001", d.Get("security_group_mapping_id").(string))
	assert.Equal(t, "sg-src", d.Get("source_security_group_id").(string))
	assert.Equal(t, "sg-target", d.Get("target_security_group_id").(string))
	assert.Equal(t, "NORMAL", d.Get("life_state").(string))
}

// TestBdrcSecurityGroupMapping_Read verifies Read populates fields from Describe response.
func TestBdrcSecurityGroupMapping_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSecurityGroupMapping().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeSecurityGroupMappingsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeSecurityGroupMappingsRequest) (*bdrcv20260330.DescribeSecurityGroupMappingsResponse, error) {
		return mockDescribeSecurityGroupMappingsResponse("sp-read", "sgm-read", "sg-src-read", "sg-target-read", "NORMAL"), nil
	})

	meta := newMockMetaBdrcSecurityGroupMapping()
	res := svcbdrc.ResourceTencentCloudBdrcSecurityGroupMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sp-read",
		"src_security_group_id":    "sg-src-read",
		"target_security_group_id": "sg-target-read",
	})
	d.SetId("sp-read#sgm-read")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sp-read#sgm-read", d.Id())
	assert.Equal(t, "sgm-read", d.Get("security_group_mapping_id").(string))
	assert.Equal(t, "sp-read", d.Get("site_pair_id").(string))
	assert.Equal(t, "sg-src-read", d.Get("source_security_group_id").(string))
	assert.Equal(t, "sg-target-read", d.Get("target_security_group_id").(string))
	assert.Equal(t, "NORMAL", d.Get("life_state").(string))
}

// TestBdrcSecurityGroupMapping_Read_NotFound verifies Read clears id when mapping not found.
func TestBdrcSecurityGroupMapping_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSecurityGroupMapping().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeSecurityGroupMappingsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeSecurityGroupMappingsRequest) (*bdrcv20260330.DescribeSecurityGroupMappingsResponse, error) {
		return mockEmptyDescribeSecurityGroupMappingsResponse(), nil
	})

	meta := newMockMetaBdrcSecurityGroupMapping()
	res := svcbdrc.ResourceTencentCloudBdrcSecurityGroupMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sp-missing",
		"src_security_group_id":    "sg-src-missing",
		"target_security_group_id": "sg-target-missing",
	})
	d.SetId("sp-missing#sgm-missing")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestBdrcSecurityGroupMapping_Update_Immutable verifies Update returns error when immutable field changes.
func TestBdrcSecurityGroupMapping_Update_Immutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSecurityGroupMapping().client, "UseBdrcV20260330Client", bdrcClient)

	meta := newMockMetaBdrcSecurityGroupMapping()
	res := svcbdrc.ResourceTencentCloudBdrcSecurityGroupMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sp-upd",
		"src_security_group_id":    "sg-src-upd",
		"target_security_group_id": "sg-target-upd",
	})
	d.SetId("sp-upd#sgm-upd")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "src_security_group_id"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "src_security_group_id")
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestBdrcSecurityGroupMapping_Delete verifies Delete passes correct params.
func TestBdrcSecurityGroupMapping_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSecurityGroupMapping().client, "UseBdrcV20260330Client", bdrcClient)

	var capturedRequest *bdrcv20260330.DeleteSecurityGroupMappingRequest
	patches.ApplyMethodFunc(bdrcClient, "DeleteSecurityGroupMappingWithContext", func(ctx context.Context, request *bdrcv20260330.DeleteSecurityGroupMappingRequest) (*bdrcv20260330.DeleteSecurityGroupMappingResponse, error) {
		capturedRequest = request
		resp := bdrcv20260330.NewDeleteSecurityGroupMappingResponse()
		resp.Response = &bdrcv20260330.DeleteSecurityGroupMappingResponseParams{
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSecurityGroupMapping()
	res := svcbdrc.ResourceTencentCloudBdrcSecurityGroupMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sp-del",
		"src_security_group_id":    "sg-src-del",
		"target_security_group_id": "sg-target-del",
	})
	d.SetId("sp-del#sgm-del")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest.SitePairId)
	assert.Equal(t, "sp-del", *capturedRequest.SitePairId)
	assert.NotNil(t, capturedRequest.SecurityGroupMappingIds)
	assert.Equal(t, 1, len(capturedRequest.SecurityGroupMappingIds))
	assert.Equal(t, "sgm-del", *capturedRequest.SecurityGroupMappingIds[0])
}
