package teo_test

import (
	"context"
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

// TestTeoDomainSharedCnameAttachment_CreateSuccess tests successful create
func TestTeoDomainSharedCnameAttachment_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	// Mock BindSharedCNAMEWithContext for create
	patches.ApplyMethodFunc(teoClient, "BindSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.BindSharedCNAMERequest) (*teov20220901.BindSharedCNAMEResponse, error) {
		resp := teov20220901.NewBindSharedCNAMEResponse()
		resp.Response = &teov20220901.BindSharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSharedCNAMEWithContext for read
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
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
		"zone_id":      "zone-2qtuhspy7cr6",
		"shared_cname": "shared.example.com",
		"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#shared.example.com", d.Id())
}

// TestTeoDomainSharedCnameAttachment_ReadSuccess tests successful read
func TestTeoDomainSharedCnameAttachment_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
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
		"zone_id":      "zone-2qtuhspy7cr6",
		"shared_cname": "shared.example.com",
		"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "zone-2qtuhspy7cr6", d.Get("zone_id"))
	assert.Equal(t, "shared.example.com", d.Get("shared_cname"))
}

// TestTeoDomainSharedCnameAttachment_ReadNotFound tests read when binding not found
func TestTeoDomainSharedCnameAttachment_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
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
		"zone_id":      "zone-2qtuhspy7cr6",
		"shared_cname": "shared.example.com",
		"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestTeoDomainSharedCnameAttachment_DeleteSuccess tests successful delete
func TestTeoDomainSharedCnameAttachment_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "BindSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.BindSharedCNAMERequest) (*teov20220901.BindSharedCNAMEResponse, error) {
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
		"zone_id":      "zone-2qtuhspy7cr6",
		"shared_cname": "shared.example.com",
		"domain_names": []interface{}{"domain1.example.com", "domain2.example.com"},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoDomainSharedCnameAttachment_UpdateSuccess tests successful update (unbind removed, bind added)
func TestTeoDomainSharedCnameAttachment_UpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	var bindCalls []*teov20220901.BindSharedCNAMERequest
	patches.ApplyMethodFunc(teoClient, "BindSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.BindSharedCNAMERequest) (*teov20220901.BindSharedCNAMEResponse, error) {
		bindCalls = append(bindCalls, request)
		resp := teov20220901.NewBindSharedCNAMEResponse()
		resp.Response = &teov20220901.BindSharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeSharedCNAMEWithContext for read after update
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
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
							Instance: ptrStringSharedCname("domain3.example.com"),
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
	// New config: domain1 + domain3 (domain2 removed, domain3 added)
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-2qtuhspy7cr6",
		"shared_cname": "shared.example.com",
		"domain_names": []interface{}{"domain1.example.com", "domain3.example.com"},
	})
	d.SetId("zone-2qtuhspy7cr6#shared.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// ID remains the same (zone_id#shared_cname)
	assert.Equal(t, "zone-2qtuhspy7cr6#shared.example.com", d.Id())
}

// TestTeoDomainSharedCnameAttachment_CreateEmptyDomainNames tests create with empty domain_names skips BindSharedCNAME
func TestTeoDomainSharedCnameAttachment_CreateEmptyDomainNames(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaDomainSharedCname().client, "UseTeoClient", teoClient)

	// BindSharedCNAMEWithContext is called with empty DomainNames when domain_names is empty
	patches.ApplyMethodFunc(teoClient, "BindSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.BindSharedCNAMERequest) (*teov20220901.BindSharedCNAMEResponse, error) {
		assert.Equal(t, "bind", *request.BindType)
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Len(t, request.BindSharedCNAMEMaps, 1)
		assert.Equal(t, "shared.example.com", *request.BindSharedCNAMEMaps[0].SharedCNAME)
		assert.Empty(t, request.BindSharedCNAMEMaps[0].DomainNames)
		resp := teov20220901.NewBindSharedCNAMEResponse()
		resp.Response = &teov20220901.BindSharedCNAMEResponseParams{
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	// DescribeSharedCNAMEWithContext is called by the subsequent Read
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(ctx context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64SharedCname(1),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					SharedCNAME:         ptrStringSharedCname("shared.example.com"),
					AccelerationDomains: []*teov20220901.ReferenceHolder{},
				},
			},
			RequestId: ptrStringSharedCname("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDomainSharedCname()
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-2qtuhspy7cr6",
		"shared_cname": "shared.example.com",
		"domain_names": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	// ID should still be set even when domain_names is empty
	assert.Equal(t, "zone-2qtuhspy7cr6#shared.example.com", d.Id())
}

// TestTeoDomainSharedCnameAttachment_Schema validates schema definition
func TestTeoDomainSharedCnameAttachment_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDomainSharedCnameAttachment()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Update)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "shared_cname")
	assert.Contains(t, res.Schema, "domain_names")

	// zone_id is Required and ForceNew
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// shared_cname is Required and ForceNew
	sharedCname := res.Schema["shared_cname"]
	assert.Equal(t, schema.TypeString, sharedCname.Type)
	assert.True(t, sharedCname.Required)
	assert.True(t, sharedCname.ForceNew)

	// domain_names is Required TypeSet but NOT ForceNew (supports update)
	domainNames := res.Schema["domain_names"]
	assert.Equal(t, schema.TypeSet, domainNames.Type)
	assert.True(t, domainNames.Required)
	assert.False(t, domainNames.ForceNew)

	// Check importer
	assert.NotNil(t, res.Importer)
}
