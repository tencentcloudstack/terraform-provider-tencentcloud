package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudGaapLayer4Listener_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_TcpDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerTcpDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerUpdateNameAndHealthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "10"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerUpdateNoHealthCheck,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_updateRealserverSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerTcpUpdateRealserverSet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_udp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUdp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_udpDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUdpDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_udpUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUdp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerUdpUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_gaap_layer4_listener.foo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "80", "80"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "delay_loop"),
					resource.TestCheckNoResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
				),
			},
		},
	})
}

const testAccGaapLayer4ListenerBasic = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
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
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapLayer4ListenerTcpDomain = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  domain = "qq.com"
  name   = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "TCP"
  name                = "ci-test-gaap-4-listener"
  port                = 80
  realserver_type     = "DOMAIN"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  health_check        = true
  delay_loop          = 3
  connect_timeout     = 3
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.domain}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.domain}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapLayer4ListenerUpdateNameAndHealthConfig = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "TCP"
  name                = "ci-test-gaap-4-listener-new"
  port                = 80
  realserver_type     = "IP"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  health_check        = true
  delay_loop          = 10
  connect_timeout     = 10
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapLayer4ListenerUpdateNoHealthCheck = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "TCP"
  name                = "ci-test-gaap-4-listener-new"
  port                = 80
  realserver_type     = "IP"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  health_check        = false
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapLayer4ListenerTcpUpdateRealserverSet = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
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
    },
  ]
}
`

const testAccGaapLayer4ListenerUdp = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "UDP"
  name                = "ci-test-gaap-4-udp-listener"
  port                = 80
  realserver_type     = "IP"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapLayer4ListenerUdpDomain = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq.com"
  name   = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  domain = "qq.com"
  name   = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "UDP"
  name                = "ci-test-gaap-4-udp-listener"
  port                = 80
  realserver_type     = "DOMAIN"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.domain}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.domain}"
      port   = 80
      weight = 2
    }
  ]
}
`

const testAccGaapLayer4ListenerUdpUpdate = `
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "unknown" // TODO
  realserver_region = "unknown" // TODO
}

resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource tencentcloud_gaap_realserver "bar" {
  ip   = "119.29.29.29"
  name = "ci-test-gaap-realserver2"
}

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol            = "UDP"
  name                = "ci-test-gaap-4-udp-listener-new"
  port                = 80
  realserver_type     = "IP"
  proxy_id            = "${tencentcloud_gaap_proxy.foo.id}"
  realserver_bind_set = [
    {
      id   = "${tencentcloud_gaap_realserver.foo}"
      ip   = "${tencentcloud_gaap_realserver.foo.ip}"
      port = 80
    },
    {
      id     = "${tencentcloud_gaap_realserver.bar}"
      ip     = "${tencentcloud_gaap_realserver.bar.ip}"
      port   = 80
      weight = 2
    }
  ]
}
`
