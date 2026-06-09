package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaOriginAcl implements tccommon.ProviderMeta
type mockMetaOriginAcl struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaOriginAcl) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaOriginAcl{}

func newMockMetaOriginAcl() *mockMetaOriginAcl {
	return &mockMetaOriginAcl{client: &connectivity.TencentCloudClient{}}
}

func ptrStringOriginAcl(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginAcl" -v -count=1 -gcflags="all=-l"

// TestTeoOriginAcl_CreateWithOriginACLFamily tests Create with origin_acl_family specified
func TestTeoOriginAcl_CreateWithOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "EnableOriginACLWithContext", func(_ context.Context, request *teov20220901.EnableOriginACLRequest) (*teov20220901.EnableOriginACLResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.OriginACLFamily)
		assert.Equal(t, "mlc", *request.OriginACLFamily)

		resp := teov20220901.NewEnableOriginACLResponse()
		resp.Response = &teov20220901.EnableOriginACLResponseParams{
			RequestId: ptrStringOriginAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:          ptrStringOriginAcl("online"),
				L7Hosts:         []*string{ptrStringOriginAcl("example1.com")},
				L4ProxyIds:      []*string{ptrStringOriginAcl("sid-abc123")},
				OriginACLFamily: ptrStringOriginAcl("mlc"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAcl()
	res := teo.ResourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":           "zone-12345678",
		"origin_acl_family": "mlc",
		"l7_hosts":          []interface{}{"example1.com"},
		"l4_proxy_ids":      []interface{}{"sid-abc123"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
	assert.Equal(t, "mlc", d.Get("origin_acl_family"))
}

// TestTeoOriginAcl_CreateWithoutOriginACLFamily tests Create without origin_acl_family specified
func TestTeoOriginAcl_CreateWithoutOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "EnableOriginACLWithContext", func(_ context.Context, request *teov20220901.EnableOriginACLRequest) (*teov20220901.EnableOriginACLResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Nil(t, request.OriginACLFamily)

		resp := teov20220901.NewEnableOriginACLResponse()
		resp.Response = &teov20220901.EnableOriginACLResponseParams{
			RequestId: ptrStringOriginAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:          ptrStringOriginAcl("online"),
				L7Hosts:         []*string{ptrStringOriginAcl("example1.com")},
				OriginACLFamily: ptrStringOriginAcl("gaz"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAcl()
	res := teo.ResourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"l7_hosts":     []interface{}{"example1.com"},
		"l4_proxy_ids": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
	assert.Equal(t, "gaz", d.Get("origin_acl_family"))
}

// TestTeoOriginAcl_ReadWithOriginACLFamily tests Read with OriginACLFamily in response
func TestTeoOriginAcl_ReadWithOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:          ptrStringOriginAcl("online"),
				L7Hosts:         []*string{ptrStringOriginAcl("example1.com")},
				L4ProxyIds:      []*string{ptrStringOriginAcl("sid-abc123")},
				OriginACLFamily: ptrStringOriginAcl("mlc"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAcl()
	res := teo.ResourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mlc", d.Get("origin_acl_family"))
}

// TestTeoOriginAcl_ReadWithNilOriginACLFamily tests Read when OriginACLFamily is nil
func TestTeoOriginAcl_ReadWithNilOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAcl().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:     ptrStringOriginAcl("online"),
				L7Hosts:    []*string{ptrStringOriginAcl("example1.com")},
				L4ProxyIds: []*string{ptrStringOriginAcl("sid-abc123")},
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAcl()
	res := teo.ResourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("origin_acl_family"))
}

// TestTeoOriginAcl_UpdateOriginACLFamily tests Update with origin_acl_family change
func TestTeoOriginAcl_UpdateOriginACLFamily(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAcl().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyOriginACLWithContext", func(_ context.Context, request *teov20220901.ModifyOriginACLRequest) (*teov20220901.ModifyOriginACLResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.OriginACLFamily)
		assert.Equal(t, "emc", *request.OriginACLFamily)
		modifyCalled = true

		resp := teov20220901.NewModifyOriginACLResponse()
		resp.Response = &teov20220901.ModifyOriginACLResponseParams{
			RequestId: ptrStringOriginAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:          ptrStringOriginAcl("online"),
				L7Hosts:         []*string{ptrStringOriginAcl("example1.com")},
				OriginACLFamily: ptrStringOriginAcl("emc"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAcl()
	res := teo.ResourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":           "zone-12345678",
		"origin_acl_family": "emc",
		"l7_hosts":          []interface{}{"example1.com"},
		"l4_proxy_ids":      []interface{}{},
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "emc", d.Get("origin_acl_family"))
}

// TestTeoOriginAcl_UpdateOriginACLFamilyOnly tests Update when only origin_acl_family changes (no entity changes)
func TestTeoOriginAcl_UpdateOriginACLFamilyOnly(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaOriginAcl().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyOriginACLWithContext", func(_ context.Context, request *teov20220901.ModifyOriginACLRequest) (*teov20220901.ModifyOriginACLResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.OriginACLFamily)
		assert.Equal(t, "mlc", *request.OriginACLFamily)
		modifyCalled = true

		resp := teov20220901.NewModifyOriginACLResponse()
		resp.Response = &teov20220901.ModifyOriginACLResponseParams{
			RequestId: ptrStringOriginAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeOriginACLWithContext", func(_ context.Context, request *teov20220901.DescribeOriginACLRequest) (*teov20220901.DescribeOriginACLResponse, error) {
		resp := teov20220901.NewDescribeOriginACLResponse()
		resp.Response = &teov20220901.DescribeOriginACLResponseParams{
			OriginACLInfo: &teov20220901.OriginACLInfo{
				Status:          ptrStringOriginAcl("online"),
				L7Hosts:         []*string{ptrStringOriginAcl("example1.com")},
				OriginACLFamily: ptrStringOriginAcl("mlc"),
			},
		}
		return resp, nil
	})

	meta := newMockMetaOriginAcl()
	res := teo.ResourceTencentCloudTeoOriginAcl()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":           "zone-12345678",
		"origin_acl_family": "mlc",
		"l7_hosts":          []interface{}{"example1.com"},
		"l4_proxy_ids":      []interface{}{},
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "mlc", d.Get("origin_acl_family"))
}

// TestTeoOriginAcl_Schema tests the schema definition
func TestTeoOriginAcl_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoOriginAcl()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "origin_acl_family")

	originACLFamily := res.Schema["origin_acl_family"]
	assert.Equal(t, schema.TypeString, originACLFamily.Type)
	assert.True(t, originACLFamily.Optional)
	assert.True(t, originACLFamily.Computed)
}

// TestAccTencentCloudTeoOriginAclResource_basic is the acceptance test
func TestAccTencentCloudTeoOriginAclResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "zone_id"),
				),
			},
			{
				Config: testAccTeoOriginAclUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_acl.example", "zone_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_acl.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoOriginAcl = `
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id        = "zone-3edjdliiw3he"
  l7_hosts = [
    "1.makn.cn",
    "2.makn.cn",
    "3.makn.cn",
    "4.makn.cn",
    "5.makn.cn",
    "6.makn.cn",
  ]

  l4_proxy_ids = [
    "sid-3edjfy5n10wh",
    "sid-3edjg2kml7pg",
    "sid-3edjg6t8dw78",
    "sid-3edjgc30nbgx",
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
`

const testAccTeoOriginAclUpdate = `
resource "tencentcloud_teo_origin_acl" "example" {
  zone_id        = "zone-3edjdliiw3he"
  l7_hosts = [
    "1.makn.cn",
    "2.makn.cn",
    "3.makn.cn",
  ]

  l4_proxy_ids = [
    "sid-3edjg6t8dw78",
    "sid-3edjgc30nbgx",
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
`
