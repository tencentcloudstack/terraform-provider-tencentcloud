package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapLayer7Listeners_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersListenerId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.listenerId"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.create_time"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersListenerName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.listenerName"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.create_time"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.listenerId"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.create_time"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.name"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.name", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.create_time"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.create_time"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapLayer7ListenersListenerId = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "listenerId" {
  protocol    = "HTTP"
  proxy_id    = "${tencentcloud_gaap_proxy.foo.id}"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersListenerName = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "listenerName" {
  protocol      = "HTTP"
  proxy_id      = "${tencentcloud_gaap_proxy.foo.id}"
  listener_name = "${tencentcloud_gaap_layer7_listener.foo.name}"
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersPort = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "port" {
  protocol = "HTTP"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  port     = 80
}
`

var TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerId = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  certificate_id   = "${tencentcloud_gaap_certificate.foo.id}"
  auth_type        = 0
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "listenerId" {
  protocol    = "HTTPS"
  proxy_id    = "${tencentcloud_gaap_proxy.foo.id}"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
}
` + testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF")

var TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerName = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  certificate_id   = "${tencentcloud_gaap_certificate.foo.id}"
  auth_type        = 0
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "name" {
  protocol      = "HTTPS"
  proxy_id      = "${tencentcloud_gaap_proxy.foo.id}"
  listener_name = "${tencentcloud_gaap_layer7_listener.foo.name}"
}
` + testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF")

var TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsPort = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 80
  certificate_id   = "${tencentcloud_gaap_certificate.foo.id}"
  auth_type        = 0
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "port" {
  protocol = "HTTPS"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  port     = "80"
}
` + testAccGaapCertificate(2, "<<EOF\n"+testAccGaapCertificateServerCert+"EOF", "", "<<EOF\n"+testAccGaapCertificateServerKey+"EOF")
