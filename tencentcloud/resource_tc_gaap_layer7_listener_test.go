package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapLayer7Listener_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpUpdateName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id", "ssl-id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id", "ssl-id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id", "ssl-id-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsTwoWayAuthentication(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsTwoWayAuthentication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id", "ssl-id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id", "ccid"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsForwardHttps(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsForwardHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer7_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "certificate_id", "ssl-id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
				),
			},
		},
	})
}

const testAccGaapLayer7ListenerBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}
`

const testAccGaapLayer7ListenerHttpUpdateName = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener-new"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}
`

const testAccGaapLayer7ListenerHttps = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  auth_type        = 0
}
`

const testAccGaapLayer7ListenerHttpsUpdate = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener-new"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "ssl-id-new" // TODO
  forward_protocol = "HTTPS"
  auth_type        = 0
}
`

const testAccGaapLayer7ListenerHttpsTwoWayAuthentication = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 80
  proxy_id              = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id        = "ssl-id" // TODO
  forward_protocol      = "HTTP"
  auth_type             = 1
  client_certificate_id = "ccid" // TODO
}
`

const testAccGaapLayer7ListenerHttpsForwardHttps = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTPS"
  auth_type        = 0
}
`
