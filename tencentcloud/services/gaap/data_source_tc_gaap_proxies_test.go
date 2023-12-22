package gaap_test

import (
	"fmt"
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudGaapProxies_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.foo", "proxies.0.ip"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.bandwidth", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.concurrent", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.access_region", "Guangzhou"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.foo", "proxies.0.realserver_region", "Beijing"),
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
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapProxiesProjectId,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.projectId"),
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
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.access"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.access", "proxies.0.concurrent"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.access", "proxies.0.access_region", "Guangzhou"),
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
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies.realserver"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.domain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.concurrent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies.realserver", "proxies.0.access_region"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_proxies.realserver", "proxies.0.realserver_region", "Beijing"),
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

var TestAccDataSourceTencentCloudGaapProxiesBasic = fmt.Sprintf(`
data tencentcloud_gaap_proxies "foo" {
  ids = ["%s"]
}
`, tcacctest.DefaultGaapProxyId)

const TestAccDataSourceTencentCloudGaapProxiesProjectId = `
data tencentcloud_gaap_proxies "projectId" {
  project_id = 0
}
`

const TestAccDataSourceTencentCloudGaapProxiesAccessRegion = `

data tencentcloud_gaap_proxies "access" {
  access_region = "Guangzhou"
}
`

const TestAccDataSourceTencentCloudGaapProxiesRealserverRegion = `

data tencentcloud_gaap_proxies "realserver" {
  realserver_region = "Beijing"
}
`
