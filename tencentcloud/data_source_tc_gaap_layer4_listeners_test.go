package tencentcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "8101"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.connect_timeout", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.interval", "5"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.proxy_id"),
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
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.proxy_id"),
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
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.port", "8104"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.scheduler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.health_check"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.proxy_id"),
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
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "8106"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.proxy_id")),
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
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.name", "listeners.0.proxy_id"),
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
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.port", "8106"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.realserver_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.scheduler"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.health_check", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.port", "listeners.0.proxy_id"),
				),
			},
		},
	})
}

func gaapLayer4Listener(port int) string {
	return fmt.Sprintf(`

resource "tencentcloud_gaap_layer4_listener" "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = %d
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true
  interval        = 5
  connect_timeout = 2

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, port, defaultGaapProxyId2, defaultGaapRealserverIpId1, defaultGaapRealserverIp1)
}

func gaapLayer4Listener2(port int) string {
	return fmt.Sprintf(`
resource tencentcloud_gaap_layer4_listener "bar" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener-bar"
  port            = %d
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true
  interval        = 5
  connect_timeout = 2

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, port, defaultGaapProxyId2, defaultGaapRealserverIpId2, defaultGaapRealserverIp2)
}

var TestAccDataSourceTencentCloudGaapLayer4ListenersBasic = gaapLayer4Listener(8101) + `

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol    = "TCP"
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}
`

var TestAccDataSourceTencentCloudGaapLayer4ListenersListenerName = gaapLayer4Listener(8102) + gaapLayer4Listener2(8103) + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "name" {
  protocol      = "TCP"
  proxy_id      = "%s"
  listener_name = tencentcloud_gaap_layer4_listener.foo.name
}
`, defaultGaapProxyId2)

var TestAccDataSourceTencentCloudGaapLayer4ListenersPort = gaapLayer4Listener(8104) + gaapLayer4Listener2(8105) + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "port" {
  protocol = "TCP"
  proxy_id = "%s"
  port     = tencentcloud_gaap_layer4_listener.foo.port
}
`, defaultGaapProxyId2)

var TestAccDataSourceTencentCloudGaapLayer4ListenersUDP = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "UDP"
  name            = "ci-test-gaap-4-listener"
  port            = 8106
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = false

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol    = "UDP"
  proxy_id    = "%s"
  listener_id = tencentcloud_gaap_layer4_listener.foo.id
}
`, defaultGaapProxyId2, defaultGaapRealserverIpId2, defaultGaapRealserverIp2, defaultGaapProxyId2)

var TestAccDataSourceTencentCloudGaapLayer4ListenersUDPName = TestAccDataSourceTencentCloudGaapLayer4ListenersUDP + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "name" {
  protocol      = "UDP"
  proxy_id      = "%s"
  listener_name = tencentcloud_gaap_layer4_listener.foo.name
}
`, defaultGaapProxyId2)

var TestAccDataSourceTencentCloudGaapLayer4ListenersUDPPort = TestAccDataSourceTencentCloudGaapLayer4ListenersUDP + fmt.Sprintf(`

data tencentcloud_gaap_layer4_listeners "port" {
  protocol = "UDP"
  proxy_id = "%s"
  port     = tencentcloud_gaap_layer4_listener.foo.port
}
`, defaultGaapProxyId2)
