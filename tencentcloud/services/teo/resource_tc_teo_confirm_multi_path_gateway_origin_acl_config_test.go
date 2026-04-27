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

// mockMetaConfirmMpgwAcl implements tccommon.ProviderMeta
type mockMetaConfirmMpgwAcl struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaConfirmMpgwAcl) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaConfirmMpgwAcl{}

func newMockMetaConfirmMpgwAcl() *mockMetaConfirmMpgwAcl {
	return &mockMetaConfirmMpgwAcl{client: &connectivity.TencentCloudClient{}}
}

func ptrStringConfirmMpgwAcl(s string) *string {
	return &s
}

func ptrInt64ConfirmMpgwAcl(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoConfirmMultiPathGatewayOriginAcl" -v -count=1 -gcflags="all=-l"

// TestAccTeoConfirmMultiPathGatewayOriginAcl_Create tests Create sets ID and delegates to Update
func TestAccTeoConfirmMultiPathGatewayOriginAcl_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmMpgwAcl().client, "UseTeoV20220901Client", teoClient)

	// Patch DescribeMultiPathGatewayOriginACL for Read after Update
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: &teov20220901.MultiPathGatewayCurrentOriginACL{
					Version:  ptrInt64ConfirmMpgwAcl(1),
					IsPlaned: ptrStringConfirmMpgwAcl("true"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.0.0/24")},
						IPv6: []*string{ptrStringConfirmMpgwAcl("::1/128")},
					},
				},
			},
			RequestId: ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfirmMpgwAcl()
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"gateway_id": "gw-87654321",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#gw-87654321", d.Id())
	assert.Equal(t, "zone-12345678", d.Get("zone_id"))
	assert.Equal(t, "gw-87654321", d.Get("gateway_id"))
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_CreateWithVersion tests Create with origin_acl_version specified
func TestAccTeoConfirmMultiPathGatewayOriginAcl_CreateWithVersion(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmMpgwAcl().client, "UseTeoV20220901Client", teoClient)

	// Patch ConfirmMultiPathGatewayOriginACL for Update
	patches.ApplyMethodFunc(teoClient, "ConfirmMultiPathGatewayOriginACLWithContext", func(_ context.Context, request *teov20220901.ConfirmMultiPathGatewayOriginACLRequest) (*teov20220901.ConfirmMultiPathGatewayOriginACLResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-87654321", *request.GatewayId)
		assert.Equal(t, int64(2), *request.OriginACLVersion)

		resp := teov20220901.NewConfirmMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.ConfirmMultiPathGatewayOriginACLResponseParams{
			RequestId: ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeMultiPathGatewayOriginACL for Read after Update
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: &teov20220901.MultiPathGatewayCurrentOriginACL{
					Version:  ptrInt64ConfirmMpgwAcl(2),
					IsPlaned: ptrStringConfirmMpgwAcl("true"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.0.0/24")},
						IPv6: []*string{ptrStringConfirmMpgwAcl("::1/128")},
					},
				},
			},
			RequestId: ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfirmMpgwAcl()
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-12345678",
		"gateway_id":         "gw-87654321",
		"origin_acl_version": 2,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#gw-87654321", d.Id())
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_Read tests Read operation with full ACL info
func TestAccTeoConfirmMultiPathGatewayOriginAcl_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmMpgwAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-87654321", *request.GatewayId)

		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: &teov20220901.MultiPathGatewayCurrentOriginACL{
					Version:  ptrInt64ConfirmMpgwAcl(1),
					IsPlaned: ptrStringConfirmMpgwAcl("false"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.0.0/24"), ptrStringConfirmMpgwAcl("10.0.1.0/24")},
						IPv6: []*string{ptrStringConfirmMpgwAcl("::1/128")},
					},
				},
				MultiPathGatewayNextOriginACL: &teov20220901.MultiPathGatewayNextOriginACL{
					Version: ptrInt64ConfirmMpgwAcl(2),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.0.0/24"), ptrStringConfirmMpgwAcl("10.0.2.0/24")},
						IPv6: []*string{ptrStringConfirmMpgwAcl("::1/128"), ptrStringConfirmMpgwAcl("::2/128")},
					},
					AddedAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.2.0/24")},
						IPv6: []*string{ptrStringConfirmMpgwAcl("::2/128")},
					},
					RemovedAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.1.0/24")},
						IPv6: []*string{},
					},
					NoChangeAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.0.0/24")},
						IPv6: []*string{ptrStringConfirmMpgwAcl("::1/128")},
					},
				},
			},
			RequestId: ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfirmMpgwAcl()
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"gateway_id": "gw-87654321",
	})
	d.SetId("zone-12345678#gw-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	aclInfo := d.Get("multi_path_gateway_origin_acl_info").([]interface{})
	assert.Equal(t, 1, len(aclInfo))

	aclInfoMap := aclInfo[0].(map[string]interface{})

	// Check current ACL
	currentAclList := aclInfoMap["multi_path_gateway_current_origin_acl"].([]interface{})
	assert.Equal(t, 1, len(currentAclList))
	currentAclMap := currentAclList[0].(map[string]interface{})
	assert.Equal(t, 1, currentAclMap["version"])
	assert.Equal(t, "false", currentAclMap["is_planed"])

	entireAddressesList := currentAclMap["entire_addresses"].([]interface{})
	assert.Equal(t, 1, len(entireAddressesList))
	entireAddressesMap := entireAddressesList[0].(map[string]interface{})
	ipv4List := entireAddressesMap["ipv4"].([]interface{})
	assert.Equal(t, 2, len(ipv4List))

	// Check next ACL
	nextAclList := aclInfoMap["multi_path_gateway_next_origin_acl"].([]interface{})
	assert.Equal(t, 1, len(nextAclList))
	nextAclMap := nextAclList[0].(map[string]interface{})
	assert.Equal(t, 2, nextAclMap["version"])
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_ReadNotFound tests Read when resource is not found
func TestAccTeoConfirmMultiPathGatewayOriginAcl_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmMpgwAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: nil,
			RequestId:                     ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfirmMpgwAcl()
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"gateway_id": "gw-87654321",
	})
	d.SetId("zone-12345678#gw-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_Update tests Update with origin_acl_version change
func TestAccTeoConfirmMultiPathGatewayOriginAcl_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmMpgwAcl().client, "UseTeoV20220901Client", teoClient)

	confirmCalled := false
	patches.ApplyMethodFunc(teoClient, "ConfirmMultiPathGatewayOriginACLWithContext", func(_ context.Context, request *teov20220901.ConfirmMultiPathGatewayOriginACLRequest) (*teov20220901.ConfirmMultiPathGatewayOriginACLResponse, error) {
		confirmCalled = true
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-87654321", *request.GatewayId)
		assert.Equal(t, int64(2), *request.OriginACLVersion)

		resp := teov20220901.NewConfirmMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.ConfirmMultiPathGatewayOriginACLResponseParams{
			RequestId: ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayOriginACL", func(request *teov20220901.DescribeMultiPathGatewayOriginACLRequest) (*teov20220901.DescribeMultiPathGatewayOriginACLResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayOriginACLResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayOriginACLResponseParams{
			MultiPathGatewayOriginACLInfo: &teov20220901.MultiPathGatewayOriginACLInfo{
				MultiPathGatewayCurrentOriginACL: &teov20220901.MultiPathGatewayCurrentOriginACL{
					Version:  ptrInt64ConfirmMpgwAcl(2),
					IsPlaned: ptrStringConfirmMpgwAcl("true"),
					EntireAddresses: &teov20220901.Addresses{
						IPv4: []*string{ptrStringConfirmMpgwAcl("10.0.0.0/24")},
						IPv6: []*string{},
					},
				},
			},
			RequestId: ptrStringConfirmMpgwAcl("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfirmMpgwAcl()
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-12345678",
		"gateway_id":         "gw-87654321",
		"origin_acl_version": 2,
	})
	d.SetId("zone-12345678#gw-87654321")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, confirmCalled)
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_UpdateAPIError tests Update handles API error
func TestAccTeoConfirmMultiPathGatewayOriginAcl_UpdateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmMpgwAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ConfirmMultiPathGatewayOriginACLWithContext", func(_ context.Context, request *teov20220901.ConfirmMultiPathGatewayOriginACLRequest) (*teov20220901.ConfirmMultiPathGatewayOriginACLResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaConfirmMpgwAcl()
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":            "zone-invalid",
		"gateway_id":         "gw-87654321",
		"origin_acl_version": 2,
	})
	d.SetId("zone-invalid#gw-87654321")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_Delete tests Delete is no-op
func TestAccTeoConfirmMultiPathGatewayOriginAcl_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"gateway_id": "gw-87654321",
	})
	d.SetId("zone-12345678#gw-87654321")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestAccTeoConfirmMultiPathGatewayOriginAcl_Schema validates schema definition
func TestAccTeoConfirmMultiPathGatewayOriginAcl_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoConfirmMultiPathGatewayOriginAcl()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "gateway_id")
	assert.Contains(t, res.Schema, "origin_acl_version")
	assert.Contains(t, res.Schema, "multi_path_gateway_origin_acl_info")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	gatewayId := res.Schema["gateway_id"]
	assert.Equal(t, schema.TypeString, gatewayId.Type)
	assert.True(t, gatewayId.Required)
	assert.True(t, gatewayId.ForceNew)

	originAclVersion := res.Schema["origin_acl_version"]
	assert.Equal(t, schema.TypeInt, originAclVersion.Type)
	assert.True(t, originAclVersion.Optional)

	aclInfo := res.Schema["multi_path_gateway_origin_acl_info"]
	assert.Equal(t, schema.TypeList, aclInfo.Type)
	assert.True(t, aclInfo.Computed)
}
