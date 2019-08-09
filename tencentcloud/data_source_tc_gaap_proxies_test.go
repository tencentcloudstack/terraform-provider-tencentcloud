package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapProxies_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.access_region", "unknown"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.realserver_region", "unknown"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.scalarable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.support_protocols"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.version"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapProxies_projectId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesProjectId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.concurrent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.access_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.realserver_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.scalarable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.support_protocols"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.version"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapProxies_accessRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesAccessRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.concurrent"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.access_region", "unknown"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.realserver_region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.scalarable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.support_protocols"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.version"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapProxies_realserverRegion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesRealserverRegion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.concurrent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.access_region"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.realserver_region", "unknown"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.scalarable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.support_protocols"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.forward_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.version"),
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
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

data tencentcloud_gaap_proxies "foo" {
  ids = ["${tencentcloud_gaap_proxy.foo.id}"]
}
`

const TestAccDataSourceTencentCloudGaapProxiesProjectId = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

data tencentcloud_gaap_proxies "foo" {
  project_id = 0
}
`

const TestAccDataSourceTencentCloudGaapProxiesAccessRegion = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

data tencentcloud_gaap_proxies "foo" {
  access_region = "${tencentcloud_gaap_proxy.foo.access_region}"
}
`

const TestAccDataSourceTencentCloudGaapProxiesRealserverRegion = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

data tencentcloud_gaap_proxies "foo" {
  realserver_region = "${tencentcloud_gaap_proxy.foo.realserver_region}"
}
`
