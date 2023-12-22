package gaap_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalGaapResource_proxy(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "realserver_region", "Beijing"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalable"),
					resource.TestMatchResourceAttr("tencentcloud_gaap_proxy.foo", "support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_proxy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccInternationalGaapProxyBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "Guangzhou"
  realserver_region = "Beijing"
  network_type = "normal"
}
`
