package teo_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoAvailableOriginAclFamilyDS" -v -count=1 -gcflags="all=-l"

// TestTeoAvailableOriginAclFamilyDS_ReadSuccess tests a successful single-page read with full field mapping
func TestTeoAvailableOriginAclFamilyDS_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, uint64(100), *request.Limit)
		assert.Equal(t, uint64(0), *request.Offset)
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount: ptrInt64AclFamily(1),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{
				{
					Version:         ptrString("gaz-20240101"),
					ActiveTime:      ptrString("2024-01-01T00:00:00+08:00"),
					OriginACLFamily: ptrString("gaz"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrString("192.168.1.0/24"), ptrString("192.168.2.0/24")},
						IPv6: []*string{ptrString("2400:3200::/32")},
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	infos := d.Get("origin_acl_family_infos").([]interface{})
	assert.Len(t, infos, 1)

	infoMap := infos[0].(map[string]interface{})
	assert.Equal(t, "gaz-20240101", infoMap["version"].(string))
	assert.Equal(t, "2024-01-01T00:00:00+08:00", infoMap["active_time"].(string))
	assert.Equal(t, "gaz", infoMap["origin_acl_family"].(string))

	entireAddresses := infoMap["entire_addresses"].([]interface{})
	assert.Len(t, entireAddresses, 1)
	entireMap := entireAddresses[0].(map[string]interface{})
	ipv4 := entireMap["i_pv4"].(*schema.Set).List()
	assert.Len(t, ipv4, 2)
	assert.Equal(t, "192.168.1.0/24", ipv4[0].(string))
	assert.Equal(t, "192.168.2.0/24", ipv4[1].(string))
	ipv6 := entireMap["i_pv6"].(*schema.Set).List()
	assert.Len(t, ipv6, 1)
	assert.Equal(t, "2400:3200::/32", ipv6[0].(string))
}

// TestTeoAvailableOriginAclFamilyDS_Paginated tests that results from multiple pages are accumulated
func TestTeoAvailableOriginAclFamilyDS_Paginated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	var callCount int32
	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			// full page (100 items == limit) so the service loop continues to page 2
			fullPage := make([]*teov20220901.OriginACLFamilyInfo, 0, 100)
			for i := 0; i < 100; i++ {
				fullPage = append(fullPage, &teov20220901.OriginACLFamilyInfo{
					Version:         ptrString(fmt.Sprintf("gaz-2024010%02d", i)),
					ActiveTime:      ptrString("2024-01-01T00:00:00+08:00"),
					OriginACLFamily: ptrString("gaz"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrString("192.168.1.0/24")},
					},
				})
			}
			resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
			resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
				TotalCount:           ptrInt64AclFamily(101),
				OriginACLFamilyInfos: fullPage,
				RequestId:            ptrString("fake-request-id-1"),
			}
			return resp, nil
		}
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount: ptrInt64AclFamily(101),
			OriginACLFamilyInfos: buildOriginACLFamilyInfos(
				"mlc-20240101", "mlc",
				[]string{"10.0.0.0/24"}, []string{"2400:3200::/32"},
			),
			RequestId: ptrString("fake-request-id-2"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	infos := d.Get("origin_acl_family_infos").([]interface{})
	assert.Len(t, infos, 101)
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount))
	assert.Equal(t, "gaz", infos[0].(map[string]interface{})["origin_acl_family"].(string))
	assert.Equal(t, "mlc", infos[100].(map[string]interface{})["origin_acl_family"].(string))
}

// TestTeoAvailableOriginAclFamilyDS_Filters tests that filters are passed through to the request
func TestTeoAvailableOriginAclFamilyDS_Filters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "OriginACLFamily", *request.Filters[0].Name)
		assert.Len(t, request.Filters[0].Values, 1)
		assert.Equal(t, "gaz", *request.Filters[0].Values[0])
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount:           ptrInt64AclFamily(0),
			OriginACLFamilyInfos: buildOriginACLFamilyInfos("gaz-20240101", "gaz", []string{"192.168.1.0/24"}, []string{}),
			RequestId:            ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "OriginACLFamily",
				"values": []interface{}{"gaz"},
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	infos := d.Get("origin_acl_family_infos").([]interface{})
	assert.Len(t, infos, 1)
	assert.Equal(t, "gaz", infos[0].(map[string]interface{})["origin_acl_family"].(string))
}

// TestTeoAvailableOriginAclFamilyDS_EmptyResult tests that an empty result list (normal empty) does not error
func TestTeoAvailableOriginAclFamilyDS_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount:           ptrInt64AclFamily(0),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{},
			RequestId:            ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	infos := d.Get("origin_acl_family_infos").([]interface{})
	assert.Len(t, infos, 0)
}

// TestTeoAvailableOriginAclFamilyDS_NilResponse tests that a nil response returns an error without clearing state id
func TestTeoAvailableOriginAclFamilyDS_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		return nil, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})
	d.SetId("existing-state-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "response is nil")
	// state id should NOT be cleared; the NonRetryableError preserves it
	assert.Equal(t, "existing-state-id", d.Id())
}

// TestTeoAvailableOriginAclFamilyDS_APIError tests that an API error is propagated
func TestTeoAvailableOriginAclFamilyDS_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceUnavailable.ZoneNotFound, Message=zone not found")
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-missing",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ZoneNotFound")
}

// TestTeoAvailableOriginAclFamilyDS_Schema validates schema definition
func TestTeoAvailableOriginAclFamilyDS_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "result_output_file")
	assert.Contains(t, res.Schema, "origin_acl_family_infos")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)

	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)

	infos := res.Schema["origin_acl_family_infos"]
	assert.Equal(t, schema.TypeList, infos.Type)
	assert.True(t, infos.Computed)
}

// buildOriginACLFamilyInfos constructs a single-element OriginACLFamilyInfo slice for test brevity
func buildOriginACLFamilyInfos(version, family string, ipv4, ipv6 []string) []*teov20220901.OriginACLFamilyInfo {
	ipv4List := make([]*string, 0, len(ipv4))
	for _, v := range ipv4 {
		ipv4List = append(ipv4List, ptrString(v))
	}
	ipv6List := make([]*string, 0, len(ipv6))
	for _, v := range ipv6 {
		ipv6List = append(ipv6List, ptrString(v))
	}
	return []*teov20220901.OriginACLFamilyInfo{
		{
			Version:         ptrString(version),
			ActiveTime:      ptrString("2024-01-01T00:00:00+08:00"),
			OriginACLFamily: ptrString(family),
			EntireAddresses: &teov20220901.Addresses{
				IPv4: ipv4List,
				IPv6: ipv6List,
			},
		},
	}
}

func ptrInt64AclFamily(v int64) *int64 {
	return &v
}
