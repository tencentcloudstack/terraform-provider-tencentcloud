package teo_test

import (
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

// aoafMockMeta implements tccommon.ProviderMeta
type aoafMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *aoafMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &aoafMockMeta{}

func ptrAOAFString(s string) *string {
	return &s
}

func ptrAOAFInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoAvailableOriginAclFamilyDS" -v -count=1 -gcflags="all=-l"

// TestTeoAvailableOriginAclFamilyDS_ReadSuccess tests Read with full response
func TestTeoAvailableOriginAclFamilyDS_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount: ptrAOAFInt64(1),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{
				{
					Version:         ptrAOAFString("gaz-12345"),
					ActiveTime:      ptrAOAFString("2024-09-01T00:00:00+08:00"),
					OriginACLFamily: ptrAOAFString("gaz"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrAOAFString("10.0.0.0/8")},
						IPv6: []*string{ptrAOAFString("::1/128")},
					},
				},
			},
			RequestId: ptrAOAFString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &aoafMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	assert.Equal(t, 1, d.Get("total_count").(int))

	infoSet := d.Get("origin_acl_family_info_set").([]interface{})
	assert.Len(t, infoSet, 1)
	infoMap := infoSet[0].(map[string]interface{})
	assert.Equal(t, "gaz-12345", infoMap["version"])
	assert.Equal(t, "2024-09-01T00:00:00+08:00", infoMap["active_time"])
	assert.Equal(t, "gaz", infoMap["origin_acl_family"])

	entireAddresses := infoMap["entire_addresses"].([]interface{})
	assert.Len(t, entireAddresses, 1)
	addrMap := entireAddresses[0].(map[string]interface{})
	ipv4Set := addrMap["ipv4"].(*schema.Set)
	assert.Contains(t, ipv4Set.List(), "10.0.0.0/8")
	ipv6Set := addrMap["ipv6"].(*schema.Set)
	assert.Contains(t, ipv6Set.List(), "::1/128")
}

// TestTeoAvailableOriginAclFamilyDS_ReadMultipleItems tests Read with multiple items
func TestTeoAvailableOriginAclFamilyDS_ReadMultipleItems(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount: ptrAOAFInt64(2),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{
				{
					Version:         ptrAOAFString("gaz-12345"),
					ActiveTime:      ptrAOAFString("2024-09-01T00:00:00+08:00"),
					OriginACLFamily: ptrAOAFString("gaz"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrAOAFString("10.0.0.0/8")},
						IPv6: []*string{ptrAOAFString("::1/128")},
					},
				},
				{
					Version:         ptrAOAFString("mlc-67890"),
					ActiveTime:      ptrAOAFString("2024-09-01T00:00:00+08:00"),
					OriginACLFamily: ptrAOAFString("mlc"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrAOAFString("172.16.0.0/12")},
						IPv6: []*string{ptrAOAFString("::2/128")},
					},
				},
			},
			RequestId: ptrAOAFString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &aoafMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, 2, d.Get("total_count").(int))

	infoSet := d.Get("origin_acl_family_info_set").([]interface{})
	assert.Len(t, infoSet, 2)
	assert.Equal(t, 1, callCount)
}

// TestTeoAvailableOriginAclFamilyDS_ReadWithNilEntireAddresses tests Read when EntireAddresses is nil
func TestTeoAvailableOriginAclFamilyDS_ReadWithNilEntireAddresses(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount: ptrAOAFInt64(1),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{
				{
					Version:         ptrAOAFString("gaz-12345"),
					ActiveTime:      ptrAOAFString("2024-09-01T00:00:00+08:00"),
					OriginACLFamily: ptrAOAFString("gaz"),
					EntireAddresses: nil,
				},
			},
			RequestId: ptrAOAFString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &aoafMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	infoSet := d.Get("origin_acl_family_info_set").([]interface{})
	assert.Len(t, infoSet, 1)
	infoMap := infoSet[0].(map[string]interface{})
	assert.Equal(t, "gaz-12345", infoMap["version"])
	assert.Equal(t, "gaz", infoMap["origin_acl_family"])

	entireAddresses := infoMap["entire_addresses"].([]interface{})
	assert.Len(t, entireAddresses, 0)
}

// TestTeoAvailableOriginAclFamilyDS_ReadEmpty tests Read handles empty response (NonRetryableError)
func TestTeoAvailableOriginAclFamilyDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount:           ptrAOAFInt64(0),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{},
			RequestId:            ptrAOAFString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &aoafMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

// TestTeoAvailableOriginAclFamilyDS_ReadWithFilters tests Read with filters set
func TestTeoAvailableOriginAclFamilyDS_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	var capturedFilters []*teov20220901.Filter
	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		capturedFilters = request.Filters
		resp := teov20220901.NewDescribeAvailableOriginACLFamilyResponse()
		resp.Response = &teov20220901.DescribeAvailableOriginACLFamilyResponseParams{
			TotalCount: ptrAOAFInt64(1),
			OriginACLFamilyInfos: []*teov20220901.OriginACLFamilyInfo{
				{
					Version:         ptrAOAFString("mlc-67890"),
					ActiveTime:      ptrAOAFString("2024-09-01T00:00:00+08:00"),
					OriginACLFamily: ptrAOAFString("mlc"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrAOAFString("172.16.0.0/12")},
						IPv6: nil,
					},
				},
			},
			RequestId: ptrAOAFString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &aoafMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "OriginACLFamily",
				"values": []interface{}{"mlc"},
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	assert.Len(t, capturedFilters, 1)
	assert.Equal(t, "OriginACLFamily", *capturedFilters[0].Name)
	assert.Contains(t, capturedFilters[0].Values, ptrAOAFString("mlc"))

	infoSet := d.Get("origin_acl_family_info_set").([]interface{})
	assert.Len(t, infoSet, 1)
	infoMap := infoSet[0].(map[string]interface{})
	assert.Equal(t, "mlc", infoMap["origin_acl_family"])

	entireAddresses := infoMap["entire_addresses"].([]interface{})
	assert.Len(t, entireAddresses, 1)
	addrMap := entireAddresses[0].(map[string]interface{})
	ipv4Set := addrMap["ipv4"].(*schema.Set)
	assert.Contains(t, ipv4Set.List(), "172.16.0.0/12")
}

// TestTeoAvailableOriginAclFamilyDS_APIError tests Read handles API error
func TestTeoAvailableOriginAclFamilyDS_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAvailableOriginACLFamily", func(request *teov20220901.DescribeAvailableOriginACLFamilyRequest) (*teov20220901.DescribeAvailableOriginACLFamilyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid ZoneId")
	})

	meta := &aoafMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "invalid-zone-id",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoAvailableOriginAclFamilyDS_Schema validates schema definition
func TestTeoAvailableOriginAclFamilyDS_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoAvailableOriginAclFamily()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "origin_acl_family_info_set")
	assert.Contains(t, res.Schema, "total_count")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)

	infoSet := res.Schema["origin_acl_family_info_set"]
	assert.Equal(t, schema.TypeList, infoSet.Type)
	assert.True(t, infoSet.Computed)

	totalCount := res.Schema["total_count"]
	assert.Equal(t, schema.TypeInt, totalCount.Type)
	assert.True(t, totalCount.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}
