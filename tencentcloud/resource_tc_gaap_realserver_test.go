package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapRealserver_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_realserver.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.1.1.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapRealserver_domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_realserver.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "domain", "www.qq.com"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "ip"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapRealserver_updateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapRealserverBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_realserver.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "ip", "1.1.1.1"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_realserver.foo", "domain"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "project_id", "0"),
				),
			},
			{
				Config: testAccGaapRealserverNewName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_realserver.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_realserver.foo", "name", "ci-test-gaap-realserver-new"),
				),
			},
		},
	})
}

const testAccGaapRealserverBasic = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverDomain = `
resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver"
}
`

const testAccGaapRealserverNewName = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver-new"
}
`
