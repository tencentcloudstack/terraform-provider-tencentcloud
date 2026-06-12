package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoDomainSharedCnameAttachment" -v -count=1 -gcflags="all=-l"

type mockMetaDomainSharedCname struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDomainSharedCname) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDomainSharedCname{}

func newMockMetaDomainSharedCname() *mockMetaDomainSharedCname {
	return &mockMetaDomainSharedCname{client: &connectivity.TencentCloudClient{}}
}

func ptrStringSharedCname(s string) *string {
	return &s
}

func ptrInt64SharedCname(n int64) *int64 {
	return &n
}

// TestTeoDomainSharedCnameAttachment_CreateSuccess tests successful create
func TestTeoDomainSharedCnameAttachment_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	// Mock BindSharedCNAME for create
	patches.ApplyMethodFunc(teoClient, "BindSharedCNAME", func(request *teov20220901.BindSharedCNAMERequest) (*teov20220901.BindSharedCNAMEResponse, error) {
		resp := teov20220901.NewBindSharedCNAMEResponse()
		resp.Response = &teov20220901.BindSharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSharedCNAME for read
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAME", func(request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					SharedCNAME:     ptrStringSharedCname("shared.example.com"),
					Type:            ptrStringSharedCname("custom"),
					BindDomainCount: ptrInt64SharedCname(2),
					AccelerationDomains: []*teov20220901.ReferenceHolder{
						{
							ZoneId:   ptrStringSharedCname("zone-2qtuhspy7cr6"),
							Type:     ptrStringSharedCname("acceleration-domain"),
							Instance: ptrStringSharedCname("domain1.example.com"),
						},
						{
							ZoneId:   ptrStringSharedCname("zone-2qtuhspy7cr6"),
							Type:     ptrStringSharedCname("acceleration-domain"),
							Instance: ptrStringSharedCname("domain2.example.com"),
						},
					},
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDomainSharedCname()
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"bind_shared_cname_maps": []interface{}{
			map[string]interface{}{
				"shared_cname": "shared.example.com",
				"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com", d.Id())
}

// TestTeoDomainSharedCnameAttachment_ReadSuccess tests successful read
func TestTeoDomainSharedCnameAttachment_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAME", func(request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					SharedCNAME:     ptrStringSharedCname("shared.example.com"),
					Type:            ptrStringSharedCname("custom"),
					BindDomainCount: ptrInt64SharedCname(2),
					AccelerationDomains: []*teov20220901.ReferenceHolder{
						{
							ZoneId:   ptrStringSharedCname("zone-2qtuhspy7cr6"),
							Type:     ptrStringSharedCname("acceleration-domain"),
							Instance: ptrStringSharedCname("domain1.example.com"),
						},
						{
							ZoneId:   ptrStringSharedCname("zone-2qtuhspy7cr6"),
							Type:     ptrStringSharedCname("acceleration-domain"),
							Instance: ptrStringSharedCname("domain2.example.com"),
						},
					},
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDomainSharedCname()
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"bind_shared_cname_maps": []interface{}{
			map[string]interface{}{
				"shared_cname": "shared.example.com",
				"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
			},
		},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "zone-2qtuhspy7cr6", d.Get("zone_id"))
}

// TestTeoDomainSharedCnameAttachment_ReadNotFound tests read when binding not found
func TestTeoDomainSharedCnameAttachment_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAME", func(request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount:      ptrInt64SharedCname(0),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{},
			RequestId:       ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDomainSharedCname()
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"bind_shared_cname_maps": []interface{}{
			map[string]interface{}{
				"shared_cname": "shared.example.com",
				"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
			},
		},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestTeoDomainSharedCnameAttachment_ReadDomainNotBound tests read when domain is not bound
func TestTeoDomainSharedCnameAttachment_ReadDomainNotBound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAME", func(request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					SharedCNAME:     ptrStringSharedCname("shared.example.com"),
					Type:            ptrStringSharedCname("custom"),
					BindDomainCount: ptrInt64SharedCname(1),
					AccelerationDomains: []*teov20220901.ReferenceHolder{
						{
							ZoneId:   ptrStringSharedCname("zone-2qtuhspy7cr6"),
							Type:     ptrStringSharedCname("acceleration-domain"),
							Instance: ptrStringSharedCname("domain1.example.com"),
						},
					},
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDomainSharedCname()
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"bind_shared_cname_maps": []interface{}{
			map[string]interface{}{
				"shared_cname": "shared.example.com",
				"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
			},
		},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// domain2.example.com is not bound, so resource should be removed from state
	assert.Empty(t, d.Id())
}

// TestTeoDomainSharedCnameAttachment_DeleteSuccess tests successful delete
func TestTeoDomainSharedCnameAttachment_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "BindSharedCNAME", func(request *teov20220901.BindSharedCNAMERequest) (*teov20220901.BindSharedCNAMEResponse, error) {
		assert.Equal(t, "unbind", *request.BindType)
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Len(t, request.BindSharedCNAMEMaps, 1)
		assert.Equal(t, "shared.example.com", *request.BindSharedCNAMEMaps[0].SharedCNAME)
		assert.Len(t, request.BindSharedCNAMEMaps[0].DomainNames, 2)

		resp := teov20220901.NewBindSharedCNAMEResponse()
		resp.Response = &teov20220901.BindSharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDomainSharedCname()
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"bind_shared_cname_maps": []interface{}{
			map[string]interface{}{
				"shared_cname": "shared.example.com",
				"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
			},
		},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDomainSharedCnameAttachment_Schema validates schema definition
func TestTeoDomainSharedCnameAttachment_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Delete)
	assert.Nil(t, res.Update)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "bind_shared_cname_maps")

	// zone_id is Required and ForceNew
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// bind_shared_cname_maps is Required and ForceNew
	bindMaps := res.Schema["bind_shared_cname_maps"]
	assert.Equal(t, schema.TypeList, bindMaps.Type)
	assert.True(t, bindMaps.Required)
	assert.True(t, bindMaps.ForceNew)

	// Check importer
	assert.NotNil(t, res.Importer)
}
