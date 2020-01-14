package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudGaapProxies_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.bandwidth", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.concurrent", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.access_region", "SouthChina"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.realserver_region", "NorthChina"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.scalable"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.version"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapProxies_filter(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesProjectId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.projectId"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.concurrent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.access_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.realserver_region"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.projectId", "proxies.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.scalable"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_proxies.projectId", "proxies.0.support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.projectId", "proxies.0.version"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesAccessRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.access"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.concurrent"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.access", "proxies.0.access_region", "SouthChina"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.realserver_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.scalable"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_proxies.access", "proxies.0.support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.version"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesRealserverRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.realserver"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.concurrent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.access_region"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.realserver", "proxies.0.realserver_region", "NorthChina"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.scalable"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_proxies.realserver", "proxies.0.support_protocols.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.version"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapProxiesBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

data tencentcloud_gaap_proxies "foo" {
  ids = [tencentcloud_gaap_proxy.foo.id]
}
`

const TestAccDataSourceTencentCloudGaapProxiesProjectId = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

data tencentcloud_gaap_proxies "projectId" {
  project_id = 0
}
`

const TestAccDataSourceTencentCloudGaapProxiesAccessRegion = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

data tencentcloud_gaap_proxies "access" {
  access_region = tencentcloud_gaap_proxy.foo.access_region
}
`

const TestAccDataSourceTencentCloudGaapProxiesRealserverRegion = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

data tencentcloud_gaap_proxies "realserver" {
  realserver_region = tencentcloud_gaap_proxy.foo.realserver_region
}
`
