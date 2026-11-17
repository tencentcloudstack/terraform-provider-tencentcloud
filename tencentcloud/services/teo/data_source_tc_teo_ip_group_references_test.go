package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoIPGroupReferencesDataSource" -v -count=1 -gcflags="all=-l"

// TestTeoIPGroupReferencesDataSource_ReadSuccess tests successful read with references data (single page)
func TestTeoIPGroupReferencesDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeIPGroupReferences", func(request *teov20220901.DescribeIPGroupReferencesRequest) (*teov20220901.DescribeIPGroupReferencesResponse, error) {
		resp := teov20220901.NewDescribeIPGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeIPGroupReferencesResponseParams{
			References: []*teov20220901.IPGroupReference{
				{
					ZoneId:        ptrString("zone-2qtuhspy7cr6"),
					EntityType:    ptrString("WebSec.ZonePolicy"),
					EntityId:      ptrString("zone-2qtuhspy7cr6"),
					EntityName:    ptrString(""),
					SubEntityType: ptrString("WebSec.ExceptionRule"),
					SubEntityId:   ptrString("rule-001"),
					SubEntityName: ptrString("exception rule"),
				},
				{
					ZoneId:        ptrString("zone-2qtuhspy7cr6"),
					EntityType:    ptrString("DDoS.L4Proxy"),
					EntityId:      ptrString("proxy-xxx"),
					EntityName:    ptrString(""),
					SubEntityType: ptrString("DDoS.L4Proxy.IpAccessControl"),
					SubEntityId:   ptrString(""),
					SubEntityName: ptrString("block"),
				},
			},
			TotalCount: ptrInt64(2),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoIPGroupReferences()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2qtuhspy7cr6",
		"group_id": 123,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	references := d.Get("references").([]interface{})
	assert.Len(t, references, 2)

	first := references[0].(map[string]interface{})
	assert.Equal(t, "zone-2qtuhspy7cr6", first["zone_id"])
	assert.Equal(t, "WebSec.ZonePolicy", first["entity_type"])
	assert.Equal(t, "zone-2qtuhspy7cr6", first["entity_id"])
	assert.Equal(t, "WebSec.ExceptionRule", first["sub_entity_type"])
	assert.Equal(t, "rule-001", first["sub_entity_id"])
	assert.Equal(t, "exception rule", first["sub_entity_name"])

	second := references[1].(map[string]interface{})
	assert.Equal(t, "DDoS.L4Proxy", second["entity_type"])
	assert.Equal(t, "proxy-xxx", second["entity_id"])
	assert.Equal(t, "DDoS.L4Proxy.IpAccessControl", second["sub_entity_type"])
	assert.Equal(t, "block", second["sub_entity_name"])
}

// TestTeoIPGroupReferencesDataSource_ReadPaginated tests pagination across multiple pages
func TestTeoIPGroupReferencesDataSource_ReadPaginated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeIPGroupReferences", func(request *teov20220901.DescribeIPGroupReferencesRequest) (*teov20220901.DescribeIPGroupReferencesResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeIPGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeIPGroupReferencesResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		if callCount == 1 {
			// first page: full page (limit=200), so loop continues
			resp.Response.References = make([]*teov20220901.IPGroupReference, 200)
			for i := range resp.Response.References {
				resp.Response.References[i] = &teov20220901.IPGroupReference{
					ZoneId:     ptrString("zone-paginated"),
					EntityType: ptrString("WebSec.HostPolicy"),
					EntityId:   ptrString("example.com"),
				}
			}
			resp.Response.TotalCount = ptrInt64(201)
		} else {
			// second page: less than limit, loop ends
			resp.Response.References = []*teov20220901.IPGroupReference{
				{
					ZoneId:     ptrString("zone-paginated"),
					EntityType: ptrString("WebSec.Template"),
					EntityId:   ptrString("template-001"),
				},
			}
			resp.Response.TotalCount = ptrInt64(201)
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoIPGroupReferences()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-paginated",
		"group_id": 456,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 2, callCount)

	references := d.Get("references").([]interface{})
	assert.Len(t, references, 201)
}

// TestTeoIPGroupReferencesDataSource_ReadEmpty tests read with empty references returns error (NonRetryableError)
func TestTeoIPGroupReferencesDataSource_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeIPGroupReferences", func(request *teov20220901.DescribeIPGroupReferencesRequest) (*teov20220901.DescribeIPGroupReferencesResponse, error) {
		resp := teov20220901.NewDescribeIPGroupReferencesResponse()
		resp.Response = &teov20220901.DescribeIPGroupReferencesResponseParams{
			References: []*teov20220901.IPGroupReference{},
			TotalCount: ptrInt64(0),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoIPGroupReferences()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-empty",
		"group_id": 789,
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	// NonRetryableError should not clear the id; id remains empty (never set)
	assert.Empty(t, d.Id())
}

// TestTeoIPGroupReferencesDataSource_Schema validates schema definition
func TestTeoIPGroupReferencesDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoIPGroupReferences()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "group_id")
	assert.Contains(t, res.Schema, "references")
	assert.Contains(t, res.Schema, "result_output_file")

	// zone_id is Required TypeString
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	// group_id is Required TypeInt
	groupId := res.Schema["group_id"]
	assert.Equal(t, schema.TypeInt, groupId.Type)
	assert.True(t, groupId.Required)

	// references is Computed TypeList
	references := res.Schema["references"]
	assert.Equal(t, schema.TypeList, references.Type)
	assert.True(t, references.Computed)

	// references element is a Resource with expected sub-fields
	elemResource, ok := references.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.Contains(t, elemResource.Schema, "zone_id")
	assert.Contains(t, elemResource.Schema, "entity_type")
	assert.Contains(t, elemResource.Schema, "entity_id")
	assert.Contains(t, elemResource.Schema, "entity_name")
	assert.Contains(t, elemResource.Schema, "sub_entity_type")
	assert.Contains(t, elemResource.Schema, "sub_entity_id")
	assert.Contains(t, elemResource.Schema, "sub_entity_name")

	// all reference sub-fields are Computed TypeString
	for _, field := range []string{"zone_id", "entity_type", "entity_id", "entity_name", "sub_entity_type", "sub_entity_id", "sub_entity_name"} {
		f := elemResource.Schema[field]
		assert.Equal(t, schema.TypeString, f.Type, "field %s should be TypeString", field)
		assert.True(t, f.Computed, "field %s should be Computed", field)
	}

	// result_output_file is Optional TypeString
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}
