package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingGaapProxyResource_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaapProxyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccTestingGaapProxyUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapProxyExists("tencentcloud_gaap_proxy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy-basic"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "access_region", "Beijing"),
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
				Config: testAccTestingGaapProxyNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy-new"),
				),
			},
			{
				Config: testAccTestingGaapProxyNewBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
				),
			},
		},
	})
}

const testAccTestingGaapProxyUpdateBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-basic"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "Beijing"
  realserver_region = "Beijing"
}
`

const testAccTestingGaapProxyNewName = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "Beijing"
  realserver_region = "Beijing"
}
`

const testAccTestingGaapProxyNewBandwidth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 20
  concurrent        = 2
  access_region     = "Beijing"
  realserver_region = "Beijing"
}
`
