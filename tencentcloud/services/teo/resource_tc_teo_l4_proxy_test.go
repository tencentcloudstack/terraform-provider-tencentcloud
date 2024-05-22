package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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
