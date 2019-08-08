package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapProxy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_updateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewname,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_updateBandwidth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewBandwidth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_updateConcurrent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyNewConcurrent,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapProxy_enable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },

		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "false"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
			{
				Config: testAccGaapProxyEnable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_proxy.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "name", "ci-test-gaap-proxy"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "concurrent", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","access_region","unknown"),
					// resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo","realserver_region","unknown"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_proxy.foo", "enable", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_proxy.foo", "tags"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "access_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "scalarable"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "support_protocols"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_proxy.foo", "forward_ip"),
				),
			},
		},
	})
}

const testAccGaapProxyBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}
`

const testAccGaapProxyNewname = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy-new"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}
`

const testAccGaapProxyNewBandwidth = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 20
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}
`

const testAccGaapProxyNewConcurrent = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 10
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}
`

const testAccGaapProxyDisable = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
  enable            = false
}
`

const testAccGaapProxyEnable = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
  enable            = true
}
`
