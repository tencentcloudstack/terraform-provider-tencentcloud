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

// mpgoaMockMeta implements tccommon.ProviderMeta
type mpgoaMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mpgoaMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mpgoaMockMeta{}

func ptrMPGOAString(s string) *string {
	return &s
}

func ptrMPGOAInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoMultiPathGatewayOriginAclDS" -v -count=1 -gcflags="all=-l"

// TestTeoMultiPathGatewayOriginAclDS_ReadSuccess tests Read with full response
func TestTeoMultiPathGatewayOriginAclDS_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: &teov20220901.MultiPathGatewayCurrentOriginACL{
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrMPGOAString("10.0.0.0/8")},
						IPv6: []*string{ptrMPGOAString("::1/128")},
					},
					Version:  ptrMPGOAInt64(1),
					IsPlaned: ptrMPGOAString("true"),
				},
				MultiPathGatewayNextOriginACL: &teov20220901.MultiPathGatewayNextOriginACL{
					Version: ptrMPGOAInt64(2),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrMPGOAString("10.0.0.0/8"), ptrMPGOAString("172.16.0.0/12")},
						IPv6: []*string{ptrMPGOAString("::1/128"), ptrMPGOAString("::2/128")},
					},
					AddedAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrMPGOAString("172.16.0.0/12")},
						IPv6: []*string{ptrMPGOAString("::2/128")},
					},
					RemovedAddresses: &teov20220901.Addresses{
						IPv4: []*string{},
						IPv6: []*string{},
					},
					NoChangeAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrMPGOAString("10.0.0.0/8")},
						IPv6: []*string{ptrMPGOAString("::1/128")},
					},
				},
			},
			RequestId: ptrMPGOAString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mpgoaMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2noqxz9b6klw",
		"gateway_id": "gw-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noqxz9b6klw#gw-12345678", d.Id())

	originAclInfo := d.Get("multi_path_gateway_origin_acl_info").([]interface{})
	assert.Len(t, originAclInfo, 1)
	infoMap := originAclInfo[0].(map[string]interface{})

	currentAclList := infoMap["multi_path_gateway_current_origin_acl"].([]interface{})
	assert.Len(t, currentAclList, 1)
	currentAcl := currentAclList[0].(map[string]interface{})
	assert.Equal(t, "1", currentAcl["version"])
	assert.Equal(t, "true", currentAcl["is_planed"])

	entireAddresses := currentAcl["entire_addresses"].([]interface{})
	assert.Len(t, entireAddresses, 1)
	addrMap := entireAddresses[0].(map[string]interface{})
	ipv4Set := addrMap["ipv4"].(*schema.Set)
	assert.Contains(t, ipv4Set.List(), "10.0.0.0/8")

	nextAclList := infoMap["multi_path_gateway_next_origin_acl"].([]interface{})
	assert.Len(t, nextAclList, 1)
	nextAcl := nextAclList[0].(map[string]interface{})
	assert.Equal(t, "2", nextAcl["version"])
}

// TestTeoMultiPathGatewayOriginAclDS_ReadWithNilCurrentACL tests Read when current ACL is nil
func TestTeoMultiPathGatewayOriginAclDS_ReadWithNilCurrentACL(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: nil,
				MultiPathGatewayNextOriginACL:    nil,
			},
			RequestId: ptrMPGOAString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mpgoaMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2noqxz9b6klw",
		"gateway_id": "gw-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noqxz9b6klw#gw-12345678", d.Id())

	originAclInfo := d.Get("multi_path_gateway_origin_acl_info").([]interface{})
	assert.Len(t, originAclInfo, 1)
	if originAclInfo[0] != nil {
		infoMap := originAclInfo[0].(map[string]interface{})
		currentAclRaw, ok := infoMap["multi_path_gateway_current_origin_acl"]
		if ok && currentAclRaw != nil {
			currentAclList := currentAclRaw.([]interface{})
			assert.Len(t, currentAclList, 0)
		}
		nextAclRaw, ok := infoMap["multi_path_gateway_next_origin_acl"]
		if ok && nextAclRaw != nil {
			nextAclList := nextAclRaw.([]interface{})
			assert.Len(t, nextAclList, 0)
		}
	}
}

// TestTeoMultiPathGatewayOriginAclDS_ReadWithNilEntireAddresses tests Read when EntireAddresses is nil
func TestTeoMultiPathGatewayOriginAclDS_ReadWithNilEntireAddresses(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: &teov20220901.MultiPathGatewayCurrentOriginACL{
					EntireAddresses: nil,
					Version:         ptrMPGOAInt64(1),
					IsPlaned:        ptrMPGOAString("true"),
				},
			},
			RequestId: ptrMPGOAString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mpgoaMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2noqxz9b6klw",
		"gateway_id": "gw-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2noqxz9b6klw#gw-12345678", d.Id())

	originAclInfo := d.Get("multi_path_gateway_origin_acl_info").([]interface{})
	assert.Len(t, originAclInfo, 1)
	infoMap := originAclInfo[0].(map[string]interface{})

	currentAclList := infoMap["multi_path_gateway_current_origin_acl"].([]interface{})
	assert.Len(t, currentAclList, 1)
	currentAcl := currentAclList[0].(map[string]interface{})
	assert.Equal(t, "1", currentAcl["version"])
	assert.Equal(t, "true", currentAcl["is_planed"])

	entireAddresses := currentAcl["entire_addresses"].([]interface{})
	assert.Len(t, entireAddresses, 0)
}

// TestTeoMultiPathGatewayOriginAclDS_APIError tests Read handles API error
func TestTeoMultiPathGatewayOriginAclDS_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid ZoneId")
	})

	meta := &mpgoaMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "invalid-zone-id",
		"gateway_id": "gw-12345678",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoMultiPathGatewayOriginAclDS_Schema validates schema definition
func TestTeoMultiPathGatewayOriginAclDS_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "gateway_id")
	assert.Contains(t, res.Schema, "multi_path_gateway_origin_acl_info")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	gatewayId := res.Schema["gateway_id"]
	assert.Equal(t, schema.TypeString, gatewayId.Type)
	assert.True(t, gatewayId.Required)

	originAclInfo := res.Schema["multi_path_gateway_origin_acl_info"]
	assert.Equal(t, schema.TypeList, originAclInfo.Type)
	assert.True(t, originAclInfo.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}
