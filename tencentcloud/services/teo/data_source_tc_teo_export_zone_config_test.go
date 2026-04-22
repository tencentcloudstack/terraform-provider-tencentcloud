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

// TestExportZoneConfigDataSource_ExportAll tests exporting all types of zone config
func TestExportZoneConfigDataSource_ExportAll(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ExportZoneConfig", func(request *teov20220901.ExportZoneConfigRequest) (*teov20220901.ExportZoneConfigResponse, error) {
		assert.Equal(t, "zone-3fkff38fyw8s", *request.ZoneId)
		assert.Nil(t, request.Types)

		resp := teov20220901.NewExportZoneConfigResponse()
		resp.Response = &teov20220901.ExportZoneConfigResponseParams{
			Content:   ptrString(`{"L7AccelerationConfig":{},"WebSecurity":{}}`),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoExportZoneConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3fkff38fyw8s", d.Id())
	assert.Equal(t, `{"L7AccelerationConfig":{},"WebSecurity":{}}`, d.Get("content"))
}

// TestExportZoneConfigDataSource_ExportWithTypes tests exporting specific types of zone config
func TestExportZoneConfigDataSource_ExportWithTypes(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ExportZoneConfig", func(request *teov20220901.ExportZoneConfigRequest) (*teov20220901.ExportZoneConfigResponse, error) {
		assert.Equal(t, "zone-3fkff38fyw8s", *request.ZoneId)
		assert.NotNil(t, request.Types)
		assert.Len(t, request.Types, 1)
		assert.Equal(t, "L7AccelerationConfig", *request.Types[0])

		resp := teov20220901.NewExportZoneConfigResponse()
		resp.Response = &teov20220901.ExportZoneConfigResponseParams{
			Content:   ptrString(`{"L7AccelerationConfig":{}}`),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoExportZoneConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
		"types":   []interface{}{"L7AccelerationConfig"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3fkff38fyw8s", d.Id())
	assert.Equal(t, `{"L7AccelerationConfig":{}}`, d.Get("content"))
}

// TestExportZoneConfigDataSource_ExportWithMultipleTypes tests exporting multiple types of zone config
func TestExportZoneConfigDataSource_ExportWithMultipleTypes(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ExportZoneConfig", func(request *teov20220901.ExportZoneConfigRequest) (*teov20220901.ExportZoneConfigResponse, error) {
		assert.Equal(t, "zone-3fkff38fyw8s", *request.ZoneId)
		assert.NotNil(t, request.Types)
		assert.Len(t, request.Types, 2)

		resp := teov20220901.NewExportZoneConfigResponse()
		resp.Response = &teov20220901.ExportZoneConfigResponseParams{
			Content:   ptrString(`{"L7AccelerationConfig":{},"WebSecurity":{}}`),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoExportZoneConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
		"types":   []interface{}{"L7AccelerationConfig", "WebSecurity"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3fkff38fyw8s", d.Id())
	assert.Equal(t, `{"L7AccelerationConfig":{},"WebSecurity":{}}`, d.Get("content"))
}

// TestExportZoneConfigDataSource_APIError tests API error handling
func TestExportZoneConfigDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ExportZoneConfig", func(request *teov20220901.ExportZoneConfigRequest) (*teov20220901.ExportZoneConfigResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoExportZoneConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestExportZoneConfigDataSource_Schema validates schema definition
func TestExportZoneConfigDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoExportZoneConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "types")
	assert.Contains(t, res.Schema, "content")
	assert.Contains(t, res.Schema, "result_output_file")

	// Check zone_id
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	// Check types
	types := res.Schema["types"]
	assert.Equal(t, schema.TypeList, types.Type)
	assert.True(t, types.Optional)

	// Check content
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Computed)

	// Check result_output_file
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}
