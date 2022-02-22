package tencentcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudGaapLayer4Listeners_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.connect_timeout", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.interval", "5"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer4Listeners_tcp(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersListenerName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.name"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.id"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.name", regexp.MustCompile("ci-test-gaap-4-listener")),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.create_time"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer4Listeners_UDP(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.protocol", "UDP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time")),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersUDPName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.name"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.protocol", "UDP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.id"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.name", regexp.MustCompile("ci-test-gaap-4-listener")),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.scheduler"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.health_check", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.create_time"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersUDPPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.port"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.protocol", "UDP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.port", "80"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.scheduler"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.health_check", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.create_time"),
				),
			},
		},
	})
}

var gaapLayer4Listener = fmt.Sprintf(`
resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_layer4_listener" "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 80
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true
  interval        = 5
  connect_timeout = 2

  realserver_bind_set {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }
}
`, defaultGaapProxyId)

var gaapLayer4Listener2 = fmt.Sprintf(`
resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "bar" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener-bar"
  port            = 443
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true
  interval        = 5
  connect_timeout = 2

  realserver_bind_set {
    id     = tencentcloud_gaap_realserver.bar.id
    ip     = tencentcloud_gaap_realserver.bar.ip
    port   = 80
  }
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer4ListenersBasic = gaapLayer4Listener + `

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol    = "TCP"
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}
`

var TestAccDataSourceTencentCloudGaapLayer4ListenersListenerName = gaapLayer4Listener + gaapLayer4Listener2 + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "name" {
  protocol      = "TCP"
  proxy_id      = "%s"
  listener_name = tencentcloud_gaap_layer4_listener.foo.name
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer4ListenersPort = gaapLayer4Listener + gaapLayer4Listener2 + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "port" {
  protocol = "TCP"
  proxy_id = "%s"
  port     = tencentcloud_gaap_layer4_listener.foo.port
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer4ListenersUDP = fmt.Sprintf(`
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "UDP"
  name            = "ci-test-gaap-4-listener"
  port            = 80
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = false

  realserver_bind_set {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }
}

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol    = "UDP"
  proxy_id    = "%s"
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}
`, defaultGaapProxyId, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer4ListenersUDPName = TestAccDataSourceTencentCloudGaapLayer4ListenersUDP + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "name" {
  protocol      = "UDP"
  proxy_id      = "%s"
  listener_name = tencentcloud_gaap_layer4_listener.foo.name
}
`, defaultGaapProxyId)

var TestAccDataSourceTencentCloudGaapLayer4ListenersUDPPort = TestAccDataSourceTencentCloudGaapLayer4ListenersUDP + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "port" {
  protocol = "UDP"
  proxy_id = "%s"
  port     = tencentcloud_gaap_layer4_listener.foo.port
}
`, defaultGaapProxyId)
