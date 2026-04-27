package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// mockMetaL4Proxy implements tccommon.ProviderMeta
type mockMetaL4Proxy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaL4Proxy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaL4Proxy{}

func newMockMetaL4Proxy() *mockMetaL4Proxy {
	return &mockMetaL4Proxy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringL4Proxy(s string) *string {
	return &s
}

// go test -test.run TestAccTencentCloudTeoL4ProxyResource_basic -v -timeout=0
func TestAccTencentCloudTeoL4ProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckL4ProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL4Proxy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckL4ProxyExists("tencentcloud_teo_l4_proxy.teo_l4_proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy.teo_l4_proxy", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "accelerate_mainland", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "ipv6", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "proxy_name", "proxy-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "static_ip", "off"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l4_proxy.teo_l4_proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL4ProxyUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckL4ProxyExists("tencentcloud_teo_l4_proxy.teo_l4_proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy.teo_l4_proxy", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "accelerate_mainland", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "ipv6", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "proxy_name", "proxy-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "static_ip", "off"),
				),
			},
		},
	})
}

func testAccCheckL4ProxyDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_l4_proxy" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		proxy, err := service.DescribeTeoL4ProxyById(ctx, zoneId, proxyId)
		if proxy != nil {
			return fmt.Errorf("zone l4 proxy %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckL4ProxyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		proxy, err := service.DescribeTeoL4ProxyById(ctx, zoneId, proxyId)
		if proxy == nil {
			return fmt.Errorf("zone l4 proxy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoL4Proxy = `

resource "tencentcloud_teo_l4_proxy" "teo_l4_proxy" {
  accelerate_mainland = "off"
  area                = "overseas"
  ipv6                = "on"
  proxy_name          = "proxy-test"
  static_ip           = "off"
  zone_id             = "zone-2qtuhspy7cr6"
}
`

const testAccTeoL4ProxyUp = `

resource "tencentcloud_teo_l4_proxy" "teo_l4_proxy" {
  accelerate_mainland = "off"
  area                = "overseas"
  ipv6                = "off"
  proxy_name          = "proxy-test"
  static_ip           = "off"
  zone_id             = "zone-2qtuhspy7cr6"
}
`

// go test ./tencentcloud/services/teo/ -run "TestAccTeoL4ProxyProxyId_Create" -v -count=1 -gcflags="all=-l"
// TestAccTeoL4ProxyProxyId_Create tests that proxy_id is set correctly after Create
func TestAccTeoL4ProxyProxyId_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL4Proxy().client, "UseTeoClient", teoClient)

	// Mock CreateL4ProxyWithContext to return a response with ProxyId
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ context.Context, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		resp := teov20220901.NewCreateL4ProxyResponse()
		resp.Response = &teov20220901.CreateL4ProxyResponseParams{
			ProxyId:   ptrStringL4Proxy("proxy-12345678"),
			RequestId: ptrStringL4Proxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock TeoService.DescribeTeoL4ProxyById for both the state refresh in CreatePostHandleResponse and the Read call
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy("zone-test1234"),
			ProxyId:            ptrStringL4Proxy("proxy-12345678"),
			ProxyName:          ptrStringL4Proxy("proxy-test"),
			Area:               ptrStringL4Proxy("overseas"),
			Ipv6:               ptrStringL4Proxy("on"),
			StaticIp:           ptrStringL4Proxy("off"),
			AccelerateMainland: ptrStringL4Proxy("off"),
			Status:             ptrStringL4Proxy("online"),
		}, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-test1234",
		"proxy_name":          "proxy-test",
		"area":                "overseas",
		"ipv6":                "on",
		"static_ip":           "off",
		"accelerate_mainland": "off",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify proxy_id is set correctly
	proxyId := d.Get("proxy_id").(string)
	assert.Equal(t, "proxy-12345678", proxyId)

	// Verify composite ID
	assert.Equal(t, "zone-test1234#proxy-12345678", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoL4ProxyProxyId_Read" -v -count=1 -gcflags="all=-l"
// TestAccTeoL4ProxyProxyId_Read tests that proxy_id is set correctly after Read
func TestAccTeoL4ProxyProxyId_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock TeoService.DescribeTeoL4ProxyById for the Read flow
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		assert.Equal(t, "zone-test1234", zoneId)
		assert.Equal(t, "proxy-87654321", proxyId)
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy("zone-test1234"),
			ProxyId:            ptrStringL4Proxy("proxy-87654321"),
			ProxyName:          ptrStringL4Proxy("proxy-test"),
			Area:               ptrStringL4Proxy("overseas"),
			Ipv6:               ptrStringL4Proxy("off"),
			StaticIp:           ptrStringL4Proxy("off"),
			AccelerateMainland: ptrStringL4Proxy("off"),
		}, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
	})
	d.SetId("zone-test1234#proxy-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify proxy_id is set correctly from Read
	proxyId := d.Get("proxy_id").(string)
	assert.Equal(t, "proxy-87654321", proxyId)
}
