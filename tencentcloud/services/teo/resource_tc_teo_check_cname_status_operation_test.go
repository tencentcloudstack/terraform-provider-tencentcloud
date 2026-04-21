package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaCheckCnameStatus is an alias for mockMeta, reusing the shared definition.

// go test ./tencentcloud/services/teo/ -run "TestAccTeoCheckCnameStatus" -v -count=1 -gcflags="all=-l"
// TestAccTeoCheckCnameStatus_MultipleDomains tests checking CNAME status for multiple domains
func TestAccTeoCheckCnameStatus_MultipleDomains(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch CheckCnameStatus to return success with multiple domains
	patches.ApplyMethodFunc(teoClient, "CheckCnameStatus", func(request *teov20220901.CheckCnameStatusRequest) (*teov20220901.CheckCnameStatusResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, 2, len(request.RecordNames))

		resp := teov20220901.NewCheckCnameStatusResponse()
		resp.Response = &teov20220901.CheckCnameStatusResponseParams{
			CnameStatus: []*teov20220901.CnameStatus{
				{
					RecordName: ptrString("example.com"),
					Cname:      ptrString("example.com.cdn.dnsv1.com"),
					Status:     ptrString("active"),
				},
				{
					RecordName: ptrString("test.example.com"),
					Cname:      ptrString("test.example.com.cdn.dnsv1.com"),
					Status:     ptrString("active"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"record_names": []interface{}{
			"example.com",
			"test.example.com",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	cnameStatus := d.Get("cname_status").([]interface{})
	assert.Equal(t, 2, len(cnameStatus))

	firstStatus := cnameStatus[0].(map[string]interface{})
	assert.Equal(t, "example.com", firstStatus["record_name"])
	assert.Equal(t, "example.com.cdn.dnsv1.com", firstStatus["cname"])
	assert.Equal(t, "active", firstStatus["status"])

	secondStatus := cnameStatus[1].(map[string]interface{})
	assert.Equal(t, "test.example.com", secondStatus["record_name"])
	assert.Equal(t, "test.example.com.cdn.dnsv1.com", secondStatus["cname"])
	assert.Equal(t, "active", secondStatus["status"])
}

// TestAccTeoCheckCnameStatus_SingleDomain tests checking CNAME status for single domain
func TestAccTeoCheckCnameStatus_SingleDomain(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CheckCnameStatus", func(request *teov20220901.CheckCnameStatusRequest) (*teov20220901.CheckCnameStatusResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, 1, len(request.RecordNames))

		resp := teov20220901.NewCheckCnameStatusResponse()
		resp.Response = &teov20220901.CheckCnameStatusResponseParams{
			CnameStatus: []*teov20220901.CnameStatus{
				{
					RecordName: ptrString("example.com"),
					Cname:      ptrString("example.com.cdn.dnsv1.com"),
					Status:     ptrString("active"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"record_names": []interface{}{
			"example.com",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	cnameStatus := d.Get("cname_status").([]interface{})
	assert.Equal(t, 1, len(cnameStatus))

	firstStatus := cnameStatus[0].(map[string]interface{})
	assert.Equal(t, "example.com", firstStatus["record_name"])
	assert.Equal(t, "example.com.cdn.dnsv1.com", firstStatus["cname"])
	assert.Equal(t, "active", firstStatus["status"])
}

// TestAccTeoCheckCnameStatus_EmptyRecordNames tests handling empty record_names list
func TestAccTeoCheckCnameStatus_EmptyRecordNames(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"record_names": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record_names is required")
}

// TestAccTeoCheckCnameStatus_APIError tests handling API errors
func TestAccTeoCheckCnameStatus_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CheckCnameStatus", func(request *teov20220901.CheckCnameStatusRequest) (*teov20220901.CheckCnameStatusResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"record_names": []interface{}{
			"example.com",
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestAccTeoCheckCnameStatus_MissingZoneId tests error handling when zone_id is missing
func TestAccTeoCheckCnameStatus_MissingZoneId(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"record_names": []interface{}{
			"example.com",
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone_id is required")
}

// TestAccTeoCheckCnameStatus_MissingRecordNames tests error handling when record_names is missing
func TestAccTeoCheckCnameStatus_MissingRecordNames(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record_names is required")
}

// TestAccTeoCheckCnameStatus_Read tests Read is no-op
func TestAccTeoCheckCnameStatus_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"record_names": []interface{}{
			"example.com",
		},
	})
	d.SetId("zone-12345678")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoCheckCnameStatus_Delete tests Delete is no-op
func TestAccTeoCheckCnameStatus_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"record_names": []interface{}{
			"example.com",
		},
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoCheckCnameStatus_Schema validates schema definition
func TestAccTeoCheckCnameStatus_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCheckCnameStatusOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "record_names")
	assert.Contains(t, res.Schema, "cname_status")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	recordNames := res.Schema["record_names"]
	assert.Equal(t, schema.TypeList, recordNames.Type)
	assert.True(t, recordNames.Required)
	assert.True(t, recordNames.ForceNew)

	cnameStatus := res.Schema["cname_status"]
	assert.Equal(t, schema.TypeList, cnameStatus.Type)
	assert.True(t, cnameStatus.Computed)
}
