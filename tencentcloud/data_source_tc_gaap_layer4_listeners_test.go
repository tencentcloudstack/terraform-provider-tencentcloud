package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudGaapLayer4Listeners_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.connect_timeout", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.delay_loop", "3"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time")),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer4Listeners_listenerId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersListenerId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.connect_timeout", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.delay_loop", "3"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer4Listeners_listenerName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersListenerName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.connect_timeout", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.delay_loop", "3"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapLayer4Listeners_port(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapLayer4ListenersPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_layer4_listeners.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.protocol", "TCP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.port", "80"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.scheduler", "rr"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.health_check", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.connect_timeout", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.delay_loop", "3"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_layer4_listeners.foo", "listeners.0.create_time"),
				),
			},
		},
	})
}

const proxy = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}
`

const listener = `

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "TCP"
  name                = "ci-test-gaap-4-listener"
  port                = 80
  realserver_type     = "IP"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  health_check        = true
  delay_loop          = 3
  connect_timeout     = 3
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    }
  ]
}
`

const listener2 = `

resource tencentcloud_gaap_realserver "bar" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_layer4_listener "bar" {
  protocol            = "TCP"
  name                = "ci-test-gaap-4-listener-bar"
  port                = 443
  realserver_type     = "IP"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  health_check        = true
  delay_loop          = 3
  connect_timeout     = 3
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.bar}"
      ip   = "${tencentcloud_gaap_realserver.bar.ip}"
      port = 80
    }
  ]
}
`

const TestAccDataSourceTencentCloudGaapLayer4ListenersBasic = proxy + listener + `

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol = "TCP"
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapLayer4ListenersListenerId = proxy + listener + `

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol = "TCP"
  listener_id = "${tencentcloud_gaap_layer4_listener.foo.id}"
}
`

const TestAccDataSourceTencentCloudGaapLayer4ListenersListenerName = proxy + listener + listener2 + `

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol = "TCP"
  listener_name = "${tencentcloud_gaap_layer4_listener.foo.name}"
}
`

const TestAccDataSourceTencentCloudGaapLayer4ListenersPort = proxy + listener + listener2 + `

data tencentcloud_gaap_layer4_listeners "foo" {
  protocol = "TCP"
  port = "${tencentcloud_gaap_layer4_listener.foo.port}"
}
`
