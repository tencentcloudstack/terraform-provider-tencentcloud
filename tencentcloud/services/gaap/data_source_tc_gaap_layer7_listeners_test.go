package gaap_test

import (
	"fmt"
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudGaapLayer7Listeners_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersListenerId,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.listenerId"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.proxy_id"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersListenerName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.listenerName"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerName", "listeners.0.proxy_id"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersPort,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.protocol", "HTTP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.proxy_id"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer7Listeners_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerId,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.listenerId"),
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
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.listenerId", "listeners.0.proxy_id"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.name"),
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
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.name", "listeners.0.proxy_id"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsPort,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer7_listeners.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.port", "81"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.certificate_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.forward_protocol"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer7_listeners.port", "listeners.0.proxy_id"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudGaapLayer7ListenersListenerId = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}

data tencentcloud_gaap_layer7_listeners "listenerId" {
  protocol    = "HTTP"
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
}
`, tcacctest.DefaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer7ListenersListenerName = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}

data tencentcloud_gaap_layer7_listeners "listenerName" {
  protocol      = "HTTP"
  proxy_id      = "%s"
  listener_name = tencentcloud_gaap_layer7_listener.foo.name
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer7ListenersPort = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "%s"
}

data tencentcloud_gaap_layer7_listeners "port" {
  protocol = "HTTP"
  proxy_id = "%s"
  port     = tencentcloud_gaap_layer7_listener.foo.port
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerId = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 81
  certificate_id   = tencentcloud_gaap_certificate.foo.id
  auth_type        = 0
  forward_protocol = "HTTP"
  proxy_id         = "%s"
}

data tencentcloud_gaap_layer7_listeners "listenerId" {
  protocol    = "HTTPS"
  proxy_id    = "%s"
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsListenerName = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 81
  certificate_id   = tencentcloud_gaap_certificate.foo.id
  auth_type        = 0
  forward_protocol = "HTTP"
  proxy_id         = "%s"
}

data tencentcloud_gaap_layer7_listeners "name" {
  protocol      = "HTTPS"
  proxy_id      = "%s"
  listener_name = tencentcloud_gaap_layer7_listener.foo.name
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer7ListenersHttpsPort = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 81
  certificate_id   = tencentcloud_gaap_certificate.foo.id
  auth_type        = 0
  forward_protocol = "HTTP"
  proxy_id         = "%s"
}

data tencentcloud_gaap_layer7_listeners "port" {
  protocol = "HTTPS"
  proxy_id = "%s"
  port     = tencentcloud_gaap_layer7_listener.foo.port
}
`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapProxyId)
