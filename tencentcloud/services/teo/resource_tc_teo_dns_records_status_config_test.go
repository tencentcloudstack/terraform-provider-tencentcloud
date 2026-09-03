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

type mockMetaTeoDnsRecordsStatus struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoDnsRecordsStatus) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoDnsRecordsStatus{}

func newMockMetaTeoDnsRecordsStatus() *mockMetaTeoDnsRecordsStatus {
	return &mockMetaTeoDnsRecordsStatus{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDnsRecordsStatus(s string) *string {
	return &s
}

func ptrInt64DnsRecordsStatus(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/teo/ -run "TestTeoDnsRecordsStatus" -v -count=1 -gcflags="all=-l"

// TestTeoDnsRecordsStatus_Create_Success_Enable tests Create with records_to_enable
func TestTeoDnsRecordsStatus_Create_Success_Enable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsRecordsStatus(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringDnsRecordsStatus("zone-1234567890"),
					RecordId:   ptrStringDnsRecordsStatus("record-1234567890"),
					Name:       ptrStringDnsRecordsStatus("www.example.com"),
					Type:       ptrStringDnsRecordsStatus("A"),
					Location:   ptrStringDnsRecordsStatus("Default"),
					Content:    ptrStringDnsRecordsStatus("1.2.3.4"),
					TTL:        ptrInt64DnsRecordsStatus(300),
					Weight:     ptrInt64DnsRecordsStatus(-1),
					Priority:   ptrInt64DnsRecordsStatus(0),
					Status:     ptrStringDnsRecordsStatus("enable"),
					CreatedOn:  ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-1234567890",
		"records_to_enable":  []interface{}{"record-1234567890"},
		"records_to_disable": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
}

// TestTeoDnsRecordsStatus_Create_Success_Disable tests Create with records_to_disable
func TestTeoDnsRecordsStatus_Create_Success_Disable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsRecordsStatus(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringDnsRecordsStatus("zone-1234567890"),
					RecordId:   ptrStringDnsRecordsStatus("record-1234567890"),
					Name:       ptrStringDnsRecordsStatus("www.example.com"),
					Type:       ptrStringDnsRecordsStatus("A"),
					Location:   ptrStringDnsRecordsStatus("Default"),
					Content:    ptrStringDnsRecordsStatus("1.2.3.4"),
					TTL:        ptrInt64DnsRecordsStatus(300),
					Weight:     ptrInt64DnsRecordsStatus(-1),
					Priority:   ptrInt64DnsRecordsStatus(0),
					Status:     ptrStringDnsRecordsStatus("disable"),
					CreatedOn:  ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-1234567890",
		"records_to_enable":  []interface{}{},
		"records_to_disable": []interface{}{"record-1234567890"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
}

// TestTeoDnsRecordsStatus_Create_NoRecords tests Create with empty records_to_enable and records_to_disable
func TestTeoDnsRecordsStatus_Create_NoRecords(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		modifyCalled = true
		return nil, fmt.Errorf("should not be called")
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsRecordsStatus(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringDnsRecordsStatus("zone-1234567890"),
					RecordId:   ptrStringDnsRecordsStatus("record-1234567890"),
					Name:       ptrStringDnsRecordsStatus("www.example.com"),
					Type:       ptrStringDnsRecordsStatus("A"),
					Location:   ptrStringDnsRecordsStatus("Default"),
					Content:    ptrStringDnsRecordsStatus("1.2.3.4"),
					TTL:        ptrInt64DnsRecordsStatus(300),
					Weight:     ptrInt64DnsRecordsStatus(-1),
					Priority:   ptrInt64DnsRecordsStatus(0),
					Status:     ptrStringDnsRecordsStatus("enable"),
					CreatedOn:  ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-1234567890",
		"records_to_enable":  []interface{}{},
		"records_to_disable": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
	assert.False(t, modifyCalled)
}

// TestTeoDnsRecordsStatus_Create_APIError tests Create handles API error
func TestTeoDnsRecordsStatus_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-invalid",
		"records_to_enable":  []interface{}{"record-1234567890"},
		"records_to_disable": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoDnsRecordsStatus_Read_Success tests Read retrieves DNS records
func TestTeoDnsRecordsStatus_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsRecordsStatus(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringDnsRecordsStatus("zone-1234567890"),
					RecordId:   ptrStringDnsRecordsStatus("record-1234567890"),
					Name:       ptrStringDnsRecordsStatus("www.example.com"),
					Type:       ptrStringDnsRecordsStatus("A"),
					Location:   ptrStringDnsRecordsStatus("Default"),
					Content:    ptrStringDnsRecordsStatus("1.2.3.4"),
					TTL:        ptrInt64DnsRecordsStatus(300),
					Weight:     ptrInt64DnsRecordsStatus(-1),
					Priority:   ptrInt64DnsRecordsStatus(0),
					Status:     ptrStringDnsRecordsStatus("enable"),
					CreatedOn:  ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
}

// TestTeoDnsRecordsStatus_Read_NotFound tests Read handles resource not found
func TestTeoDnsRecordsStatus_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsRecordsStatus(0),
			DnsRecords: []*teov20220901.DnsRecord{},
			RequestId:  ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecordsStatus_Update_Success tests Update with records_to_enable change
func TestTeoDnsRecordsStatus_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		resp := teov20220901.NewDescribeDnsRecordsResponse()
		resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
			TotalCount: ptrInt64DnsRecordsStatus(1),
			DnsRecords: []*teov20220901.DnsRecord{
				{
					ZoneId:     ptrStringDnsRecordsStatus("zone-1234567890"),
					RecordId:   ptrStringDnsRecordsStatus("record-1234567890"),
					Name:       ptrStringDnsRecordsStatus("www.example.com"),
					Type:       ptrStringDnsRecordsStatus("A"),
					Location:   ptrStringDnsRecordsStatus("Default"),
					Content:    ptrStringDnsRecordsStatus("1.2.3.4"),
					TTL:        ptrInt64DnsRecordsStatus(300),
					Weight:     ptrInt64DnsRecordsStatus(-1),
					Priority:   ptrInt64DnsRecordsStatus(0),
					Status:     ptrStringDnsRecordsStatus("enable"),
					CreatedOn:  ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
					ModifiedOn: ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-1234567890",
		"records_to_enable":  []interface{}{"record-1234567890"},
		"records_to_disable": []interface{}{},
	})
	d.SetId("zone-1234567890")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoDnsRecordsStatus_Update_APIError tests Update handles API error during Read
func TestTeoDnsRecordsStatus_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-invalid",
		"records_to_enable":  []interface{}{"record-1234567890"},
		"records_to_disable": []interface{}{},
	})
	d.SetId("zone-invalid")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecordsStatus_Delete_Success tests Delete is a no-op
func TestTeoDnsRecordsStatus_Delete_Success(t *testing.T) {
	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
}

// TestTeoDnsRecordsStatus_Schema validates schema definition
func TestTeoDnsRecordsStatus_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check zone_id field
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check records_to_enable and records_to_disable
	assert.Contains(t, res.Schema, "records_to_enable")
	recordsToEnable := res.Schema["records_to_enable"]
	assert.Equal(t, schema.TypeList, recordsToEnable.Type)
	assert.True(t, recordsToEnable.Optional)

	assert.Contains(t, res.Schema, "records_to_disable")
	recordsToDisable := res.Schema["records_to_disable"]
	assert.Equal(t, schema.TypeList, recordsToDisable.Type)
	assert.True(t, recordsToDisable.Optional)

	// Verify pagination parameters and query-only fields are NOT exposed
	assert.NotContains(t, res.Schema, "filters")
	assert.NotContains(t, res.Schema, "sort_by")
	assert.NotContains(t, res.Schema, "sort_order")
	assert.NotContains(t, res.Schema, "match")
	assert.NotContains(t, res.Schema, "dns_records")
	assert.NotContains(t, res.Schema, "limit")
	assert.NotContains(t, res.Schema, "offset")
}
