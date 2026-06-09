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

// TestIdentifyZoneOperation_Success tests successful zone identification
func TestIdentifyZoneOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch IdentifyZone to return success
	patches.ApplyMethodFunc(teoClient, "IdentifyZone", func(request *teov20220901.IdentifyZoneRequest) (*teov20220901.IdentifyZoneResponse, error) {
		assert.Equal(t, "example.com", *request.ZoneName)
		assert.Nil(t, request.Domain) // Domain is not set in this test

		resp := teov20220901.NewIdentifyZoneResponse()
		resp.Response = &teov20220901.IdentifyZoneResponseParams{
			Ascription: &teov20220901.AscriptionInfo{
				Subdomain:   ptrString("_teo-verification"),
				RecordType:  ptrString("TXT"),
				RecordValue: ptrString("teo-verify-code-123456"),
			},
			FileAscription: &teov20220901.FileAscriptionInfo{
				IdentifyPath:    ptrString("/.well-known/teo-verification"),
				IdentifyContent: ptrString("verify-content-789012"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify ascription was set
	ascription := d.Get("ascription").([]interface{})
	assert.Len(t, ascription, 1)
	ascriptionMap := ascription[0].(map[string]interface{})
	assert.Equal(t, "_teo-verification", ascriptionMap["subdomain"])
	assert.Equal(t, "TXT", ascriptionMap["record_type"])
	assert.Equal(t, "teo-verify-code-123456", ascriptionMap["record_value"])

	// Verify file_ascription was set
	fileAscription := d.Get("file_ascription").([]interface{})
	assert.Len(t, fileAscription, 1)
	fileAscriptionMap := fileAscription[0].(map[string]interface{})
	assert.Equal(t, "/.well-known/teo-verification", fileAscriptionMap["identify_path"])
	assert.Equal(t, "verify-content-789012", fileAscriptionMap["identify_content"])
}

// TestIdentifyZoneOperation_WithDomain tests successful subdomain identification
func TestIdentifyZoneOperation_WithDomain(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IdentifyZone", func(request *teov20220901.IdentifyZoneRequest) (*teov20220901.IdentifyZoneResponse, error) {
		assert.Equal(t, "example.com", *request.ZoneName)
		assert.Equal(t, "www.example.com", *request.Domain)

		resp := teov20220901.NewIdentifyZoneResponse()
		resp.Response = &teov20220901.IdentifyZoneResponseParams{
			Ascription: &teov20220901.AscriptionInfo{
				Subdomain:   ptrString("_teo-verification.www"),
				RecordType:  ptrString("TXT"),
				RecordValue: ptrString("teo-verify-code-sub-789"),
			},
			FileAscription: &teov20220901.FileAscriptionInfo{
				IdentifyPath:    ptrString("/.well-known/teo-verification/www"),
				IdentifyContent: ptrString("verify-content-sub-456"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
		"domain":    "www.example.com",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	// Verify ascription contains subdomain-specific values
	ascription := d.Get("ascription").([]interface{})
	assert.Len(t, ascription, 1)
	ascriptionMap := ascription[0].(map[string]interface{})
	assert.Equal(t, "_teo-verification.www", ascriptionMap["subdomain"])
	assert.Equal(t, "teo-verify-code-sub-789", ascriptionMap["record_value"])

	// Verify file_ascription contains subdomain-specific values
	fileAscription := d.Get("file_ascription").([]interface{})
	assert.Len(t, fileAscription, 1)
	fileAscriptionMap := fileAscription[0].(map[string]interface{})
	assert.Equal(t, "/.well-known/teo-verification/www", fileAscriptionMap["identify_path"])
	assert.Equal(t, "verify-content-sub-456", fileAscriptionMap["identify_content"])
}

// TestIdentifyZoneOperation_MissingZoneName tests missing required parameter
func TestIdentifyZoneOperation_MissingZoneName(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()

	// Create resource data without zone_name (empty string)
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zone_name is required")
}

// TestIdentifyZoneOperation_APIError tests API error handling
func TestIdentifyZoneOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "IdentifyZone", func(request *teov20220901.IdentifyZoneRequest) (*teov20220901.IdentifyZoneResponse, error) {
		assert.Equal(t, "invalid-zone.com", *request.ZoneName)
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "invalid-zone.com",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestIdentifyZoneOperation_Read tests Read is no-op
func TestIdentifyZoneOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
	})
	d.SetId("example.com")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestIdentifyZoneOperation_Delete tests Delete is no-op
func TestIdentifyZoneOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_name": "example.com",
	})
	d.SetId("example.com")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestIdentifyZoneOperation_Schema validates schema definition
func TestIdentifyZoneOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoIdentifyZoneOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_name")
	assert.Contains(t, res.Schema, "domain")
	assert.Contains(t, res.Schema, "ascription")
	assert.Contains(t, res.Schema, "file_ascription")

	// Check zone_name
	zoneName := res.Schema["zone_name"]
	assert.Equal(t, schema.TypeString, zoneName.Type)
	assert.True(t, zoneName.Required)
	assert.True(t, zoneName.ForceNew)

	// Check domain
	domain := res.Schema["domain"]
	assert.Equal(t, schema.TypeString, domain.Type)
	assert.True(t, domain.Optional)
	assert.True(t, domain.ForceNew)

	// Check ascription
	ascription := res.Schema["ascription"]
	assert.Equal(t, schema.TypeList, ascription.Type)
	assert.True(t, ascription.Computed)

	// Check file_ascription
	fileAscription := res.Schema["file_ascription"]
	assert.Equal(t, schema.TypeList, fileAscription.Type)
	assert.True(t, fileAscription.Computed)
}
