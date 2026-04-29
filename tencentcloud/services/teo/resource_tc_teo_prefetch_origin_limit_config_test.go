package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaPrefetchOriginLimit struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPrefetchOriginLimit) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPrefetchOriginLimit{}

func newMockMetaPrefetchOriginLimit() *mockMetaPrefetchOriginLimit {
	return &mockMetaPrefetchOriginLimit{client: &connectivity.TencentCloudClient{}}
}

func ptrStringPrefetchOriginLimit(s string) *string {
	return &s
}

func ptrInt64PrefetchOriginLimit(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/teo/ -run "TestPrefetchOriginLimit" -v -count=1 -gcflags="all=-l"

// TestPrefetchOriginLimit_Create_Success tests Create calls API and sets composite ID
func TestPrefetchOriginLimit_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPrefetchOriginLimitWithContext", func(_ context.Context, request *teov20220901.ModifyPrefetchOriginLimitRequest) (*teov20220901.ModifyPrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewModifyPrefetchOriginLimitResponse()
		resp.Response = &teov20220901.ModifyPrefetchOriginLimitResponseParams{
			RequestId: ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribePrefetchOriginLimit", func(request *teov20220901.DescribePrefetchOriginLimitRequest) (*teov20220901.DescribePrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewDescribePrefetchOriginLimitResponse()
		resp.Response = &teov20220901.DescribePrefetchOriginLimitResponseParams{
			TotalCount: ptrInt64PrefetchOriginLimit(1),
			Limits: []*teov20220901.PrefetchOriginLimit{
				{
					ZoneId:     ptrStringPrefetchOriginLimit("zone-1234567890"),
					DomainName: ptrStringPrefetchOriginLimit("example.com"),
					Area:       ptrStringPrefetchOriginLimit("Overseas"),
					Bandwidth:  ptrInt64PrefetchOriginLimit(200),
				},
			},
			RequestId: ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   200,
		"enabled":     "on",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#example.com#Overseas", d.Id())
}

// TestPrefetchOriginLimit_Create_APIError tests Create handles API error
func TestPrefetchOriginLimit_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPrefetchOriginLimitWithContext", func(_ context.Context, request *teov20220901.ModifyPrefetchOriginLimitRequest) (*teov20220901.ModifyPrefetchOriginLimitResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-invalid",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   200,
		"enabled":     "on",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestPrefetchOriginLimit_Read_Success tests Read retrieves configuration data
func TestPrefetchOriginLimit_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribePrefetchOriginLimit", func(request *teov20220901.DescribePrefetchOriginLimitRequest) (*teov20220901.DescribePrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewDescribePrefetchOriginLimitResponse()
		resp.Response = &teov20220901.DescribePrefetchOriginLimitResponseParams{
			TotalCount: ptrInt64PrefetchOriginLimit(1),
			Limits: []*teov20220901.PrefetchOriginLimit{
				{
					ZoneId:     ptrStringPrefetchOriginLimit("zone-1234567890"),
					DomainName: ptrStringPrefetchOriginLimit("example.com"),
					Area:       ptrStringPrefetchOriginLimit("Overseas"),
					Bandwidth:  ptrInt64PrefetchOriginLimit(200),
				},
			},
			RequestId: ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   200,
		"enabled":     "on",
	})
	d.SetId("zone-1234567890#example.com#Overseas")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", d.Get("domain_name"))
	assert.Equal(t, "Overseas", d.Get("area"))
	assert.Equal(t, 200, d.Get("bandwidth"))
}

// TestPrefetchOriginLimit_Read_NotFound tests Read handles resource not found
func TestPrefetchOriginLimit_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribePrefetchOriginLimit", func(request *teov20220901.DescribePrefetchOriginLimitRequest) (*teov20220901.DescribePrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewDescribePrefetchOriginLimitResponse()
		resp.Response = &teov20220901.DescribePrefetchOriginLimitResponseParams{
			TotalCount: ptrInt64PrefetchOriginLimit(0),
			Limits:     []*teov20220901.PrefetchOriginLimit{},
			RequestId:  ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   200,
		"enabled":     "on",
	})
	d.SetId("zone-1234567890#example.com#Overseas")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestPrefetchOriginLimit_Update_Success tests Update calls API and then Read
func TestPrefetchOriginLimit_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPrefetchOriginLimitWithContext", func(_ context.Context, request *teov20220901.ModifyPrefetchOriginLimitRequest) (*teov20220901.ModifyPrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewModifyPrefetchOriginLimitResponse()
		resp.Response = &teov20220901.ModifyPrefetchOriginLimitResponseParams{
			RequestId: ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribePrefetchOriginLimit", func(request *teov20220901.DescribePrefetchOriginLimitRequest) (*teov20220901.DescribePrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewDescribePrefetchOriginLimitResponse()
		resp.Response = &teov20220901.DescribePrefetchOriginLimitResponseParams{
			TotalCount: ptrInt64PrefetchOriginLimit(1),
			Limits: []*teov20220901.PrefetchOriginLimit{
				{
					ZoneId:     ptrStringPrefetchOriginLimit("zone-1234567890"),
					DomainName: ptrStringPrefetchOriginLimit("example.com"),
					Area:       ptrStringPrefetchOriginLimit("Overseas"),
					Bandwidth:  ptrInt64PrefetchOriginLimit(500),
				},
			},
			RequestId: ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   500,
		"enabled":     "on",
	})
	d.SetId("zone-1234567890#example.com#Overseas")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 500, d.Get("bandwidth"))
}

// TestPrefetchOriginLimit_Delete_Success tests Delete calls API with Enabled=off
func TestPrefetchOriginLimit_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPrefetchOriginLimitWithContext", func(_ context.Context, request *teov20220901.ModifyPrefetchOriginLimitRequest) (*teov20220901.ModifyPrefetchOriginLimitResponse, error) {
		resp := teov20220901.NewModifyPrefetchOriginLimitResponse()
		resp.Response = &teov20220901.ModifyPrefetchOriginLimitResponseParams{
			RequestId: ptrStringPrefetchOriginLimit("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   200,
		"enabled":     "on",
	})
	d.SetId("zone-1234567890#example.com#Overseas")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestPrefetchOriginLimit_Delete_APIError tests Delete handles API error
func TestPrefetchOriginLimit_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchOriginLimit().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPrefetchOriginLimitWithContext", func(_ context.Context, request *teov20220901.ModifyPrefetchOriginLimitRequest) (*teov20220901.ModifyPrefetchOriginLimitResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Config not found")
	})

	meta := newMockMetaPrefetchOriginLimit()
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"domain_name": "example.com",
		"area":        "Overseas",
		"bandwidth":   200,
		"enabled":     "on",
	})
	d.SetId("zone-1234567890#example.com#Overseas")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestPrefetchOriginLimit_Schema validates schema definition
func TestPrefetchOriginLimit_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPrefetchOriginLimit()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "domain_name")
	domainName := res.Schema["domain_name"]
	assert.Equal(t, schema.TypeString, domainName.Type)
	assert.True(t, domainName.Required)
	assert.True(t, domainName.ForceNew)

	assert.Contains(t, res.Schema, "area")
	area := res.Schema["area"]
	assert.Equal(t, schema.TypeString, area.Type)
	assert.True(t, area.Required)
	assert.True(t, area.ForceNew)

	// Check required fields without ForceNew
	assert.Contains(t, res.Schema, "bandwidth")
	bandwidth := res.Schema["bandwidth"]
	assert.Equal(t, schema.TypeInt, bandwidth.Type)
	assert.True(t, bandwidth.Required)
	assert.False(t, bandwidth.ForceNew)

	assert.Contains(t, res.Schema, "enabled")
	enabled := res.Schema["enabled"]
	assert.Equal(t, schema.TypeString, enabled.Type)
	assert.True(t, enabled.Required)
	assert.False(t, enabled.ForceNew)
}
