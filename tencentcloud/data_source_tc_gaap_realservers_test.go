package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapRealservers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapRealserversBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.0.ip", "1.1.1.1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.project_id"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapRealservers_domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapRealserversDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.0.domain", "www.qq.com"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.project_id"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapRealserversBasic = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

data tencentcloud_gaap_realservers "foo" {
  ip = "1.1.1.1"
}
`

const TestAccDataSourceTencentCloudGaapRealserversDomain = `
resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver"
}

data tencentcloud_gaap_realservers "foo" {
  domain = "www.qq.com"
}
`
