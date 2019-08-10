package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapLayer7Listeners_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_listenerId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersListenerId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_listenerName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersListenerName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_port(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
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
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_httpsListenerId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_httpsListenerName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_httpsPort(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.certificate_name"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_id"),
					resource.TestCheckNoResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.client_certificate_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.auth_type", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapLayer7ListenersBasic = `
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

data tencentcloud_gaap_layer7_listeners "foo" {
  protocol = "HTTP"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersListenerId = `
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

resource tencentcloud_gaap_layer7_listener "bar" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 8080
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "foo" {
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
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_layer7_listener "bar" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener-2"
  port     = 8080
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "foo" {
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
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_layer7_listener "bar" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener-2"
  port     = 8080
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "foo" {
  protocol = "HTTP"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  port     = 80
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersHttps = `
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

data tencentcloud_gaap_layer7_listeners "foo" {
  protocol = "HTTPS"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerId = `
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
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_layer7_listener "bar" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 8080
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "foo" {
  protocol    = "HTTPS"
  proxy_id    = "${tencentcloud_gaap_proxy.foo.id}"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerName = `
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
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_layer7_listener "bar" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 8080
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "foo" {
  protocol    = "HTTPS"
  proxy_id    = "${tencentcloud_gaap_proxy.foo.id}"
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.name}"
}
`

const TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsPort = `
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
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

resource tencentcloud_gaap_layer7_listener "bar" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 8080
  certificate_id   = "ssl-id" // TODO
  forward_protocol = "HTTP"
  proxy_id         = "${tencentcloud_gaap_proxy.foo.id}"
}

data tencentcloud_gaap_layer7_listeners "foo" {
  protocol = "HTTPS"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  port     = "80"
}
`
