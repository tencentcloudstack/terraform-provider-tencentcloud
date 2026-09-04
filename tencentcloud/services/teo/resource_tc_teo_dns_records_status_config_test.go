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

// buildDescribeDnsRecordsResponse builds a DescribeDnsRecords response with one record.
func buildDescribeDnsRecordsResponse(recordId, status string) *teov20220901.DescribeDnsRecordsResponse {
	resp := teov20220901.NewDescribeDnsRecordsResponse()
	resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
		TotalCount: ptrInt64DnsRecordsStatus(1),
		DnsRecords: []*teov20220901.DnsRecord{
			{
				ZoneId:     ptrStringDnsRecordsStatus("zone-1234567890"),
				RecordId:   ptrStringDnsRecordsStatus(recordId),
				Name:       ptrStringDnsRecordsStatus("www.example.com"),
				Type:       ptrStringDnsRecordsStatus("A"),
				Location:   ptrStringDnsRecordsStatus("Default"),
				Content:    ptrStringDnsRecordsStatus("1.2.3.4"),
				TTL:        ptrInt64DnsRecordsStatus(300),
				Weight:     ptrInt64DnsRecordsStatus(-1),
				Priority:   ptrInt64DnsRecordsStatus(0),
				Status:     ptrStringDnsRecordsStatus(status),
				CreatedOn:  ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
				ModifiedOn: ptrStringDnsRecordsStatus("2024-01-01T00:00:00Z"),
			},
		},
		RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
	}
	return resp
}

// buildEmptyDescribeDnsRecordsResponse builds a DescribeDnsRecords response with no records.
func buildEmptyDescribeDnsRecordsResponse() *teov20220901.DescribeDnsRecordsResponse {
	resp := teov20220901.NewDescribeDnsRecordsResponse()
	resp.Response = &teov20220901.DescribeDnsRecordsResponseParams{
		TotalCount: ptrInt64DnsRecordsStatus(0),
		DnsRecords: []*teov20220901.DnsRecord{},
		RequestId:  ptrStringDnsRecordsStatus("fake-request-id"),
	}
	return resp
}

// patchTeoClient patches the mock meta's UseTeoV20220901Client and returns the teoClient.
func patchTeoClient(patches *gomonkey.Patches) *teov20220901.Client {
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoDnsRecordsStatus().client, "UseTeoV20220901Client", teoClient)
	return teoClient
}

// TestTeoDnsRecordsStatus_Create_Success_Enable tests Create with status=enable.
// Create sets the id to zone_id#records_id then reuses Update, which calls ModifyDnsRecordsStatus
// with RecordsToEnable populated, polls the status, and finally calls Read.
func TestTeoDnsRecordsStatus_Create_Success_Enable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		modifyCalled = true
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.NotNil(t, request.RecordsToEnable)
		assert.Equal(t, "record-1234567890", *request.RecordsToEnable[0])
		assert.Nil(t, request.RecordsToDisable)
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return buildDescribeDnsRecordsResponse("record-1234567890", "enable"), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"records_id": "record-1234567890",
		"status":     "enable",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890"+tccommon.FILED_SP+"record-1234567890", d.Id())
	assert.True(t, modifyCalled)
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecordsStatus_Create_Success_Disable tests Create with status=disable.
func TestTeoDnsRecordsStatus_Create_Success_Disable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		modifyCalled = true
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.Nil(t, request.RecordsToEnable)
		assert.NotNil(t, request.RecordsToDisable)
		assert.Equal(t, "record-1234567890", *request.RecordsToDisable[0])
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return buildDescribeDnsRecordsResponse("record-1234567890", "disable"), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"records_id": "record-1234567890",
		"status":     "disable",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890"+tccommon.FILED_SP+"record-1234567890", d.Id())
	assert.True(t, modifyCalled)
	assert.Equal(t, "disable", d.Get("status"))
}

// TestTeoDnsRecordsStatus_Create_APIError tests Create handles ModifyDnsRecordsStatus error.
func TestTeoDnsRecordsStatus_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return buildDescribeDnsRecordsResponse("record-1234567890", "enable"), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-invalid",
		"records_id": "record-1234567890",
		"status":     "enable",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoDnsRecordsStatus_Read_Success tests Read retrieves the DNS record and keeps resource in state.
func TestTeoDnsRecordsStatus_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, "id", *request.Filters[0].Name)
		assert.Equal(t, "record-1234567890", *request.Filters[0].Values[0])
		return buildDescribeDnsRecordsResponse("record-1234567890", "enable"), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-1234567890" + tccommon.FILED_SP + "record-1234567890")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890"+tccommon.FILED_SP+"record-1234567890", d.Id())
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecordsStatus_Read_NotFound tests Read clears id when DescribeDnsRecords returns empty.
func TestTeoDnsRecordsStatus_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return buildEmptyDescribeDnsRecordsResponse(), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-1234567890" + tccommon.FILED_SP + "record-1234567890")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoDnsRecordsStatus_Update_Success tests Update calls ModifyDnsRecordsStatus on status change.
func TestTeoDnsRecordsStatus_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		modifyCalled = true
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.NotNil(t, request.RecordsToEnable)
		assert.Equal(t, "record-1234567890", *request.RecordsToEnable[0])
		resp := teov20220901.NewModifyDnsRecordsStatusResponse()
		resp.Response = &teov20220901.ModifyDnsRecordsStatusResponseParams{
			RequestId: ptrStringDnsRecordsStatus("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return buildDescribeDnsRecordsResponse("record-1234567890", "enable"), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"records_id": "record-1234567890",
		"status":     "enable",
	})
	d.SetId("zone-1234567890" + tccommon.FILED_SP + "record-1234567890")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "enable", d.Get("status"))
}

// TestTeoDnsRecordsStatus_Update_APIError tests Update handles ModifyDnsRecordsStatus error.
func TestTeoDnsRecordsStatus_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := patchTeoClient(patches)

	patches.ApplyMethodFunc(teoClient, "ModifyDnsRecordsStatusWithContext", func(_ context.Context, request *teov20220901.ModifyDnsRecordsStatusRequest) (*teov20220901.ModifyDnsRecordsStatusResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	patches.ApplyMethodFunc(teoClient, "DescribeDnsRecordsWithContext", func(_ context.Context, request *teov20220901.DescribeDnsRecordsRequest) (*teov20220901.DescribeDnsRecordsResponse, error) {
		return buildDescribeDnsRecordsResponse("record-1234567890", "enable"), nil
	})

	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-invalid",
		"records_id": "record-1234567890",
		"status":     "enable",
	})
	d.SetId("zone-invalid" + tccommon.FILED_SP + "record-1234567890")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoDnsRecordsStatus_Delete_Success tests Delete is a no-op.
func TestTeoDnsRecordsStatus_Delete_Success(t *testing.T) {
	meta := newMockMetaTeoDnsRecordsStatus()
	res := teo.ResourceTencentCloudTeoDnsRecordsStatus()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"records_id": "record-1234567890",
		"status":     "enable",
	})
	d.SetId("zone-1234567890" + tccommon.FILED_SP + "record-1234567890")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890"+tccommon.FILED_SP+"record-1234567890", d.Id())
}

// TestTeoDnsRecordsStatus_Schema validates schema definition.
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

	// Check records_id field
	assert.Contains(t, res.Schema, "records_id")
	recordsId := res.Schema["records_id"]
	assert.Equal(t, schema.TypeString, recordsId.Type)
	assert.True(t, recordsId.Required)
	assert.True(t, recordsId.ForceNew)

	// Check status field
	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Required)
	assert.False(t, status.ForceNew)

	// Verify old fields are NOT exposed
	assert.NotContains(t, res.Schema, "records_to_enable")
	assert.NotContains(t, res.Schema, "records_to_disable")

	// Verify pagination parameters and query-only fields are NOT exposed
	assert.NotContains(t, res.Schema, "filters")
	assert.NotContains(t, res.Schema, "sort_by")
	assert.NotContains(t, res.Schema, "sort_order")
	assert.NotContains(t, res.Schema, "match")
	assert.NotContains(t, res.Schema, "dns_records")
	assert.NotContains(t, res.Schema, "limit")
	assert.NotContains(t, res.Schema, "offset")
}
