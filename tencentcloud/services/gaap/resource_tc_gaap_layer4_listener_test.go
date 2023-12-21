package gaap_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"

	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_gaap_layer4_listener
	resource.AddTestSweepers("tencentcloud_gaap_layer4_listener", &resource.Sweeper{
		Name: "tencentcloud_gaap_layer4_listener",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)
			service := svcgaap.NewGaapService(client.GetAPIV3Conn())
			proxyIds := []string{tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapProxyId2}
			for _, proxyId := range proxyIds {
				proxyIdTmp := proxyId
				tcpListeners, err := service.DescribeTCPListeners(ctx, &proxyIdTmp, nil, nil, nil)
				if err != nil {
					return err
				}
				for _, tcpListener := range tcpListeners {
					instanceName := *tcpListener.ListenerName

					now := time.Now()
					createTime := time.Unix(int64(*tcpListener.CreateTime), 0)
					interval := now.Sub(createTime).Minutes()

					if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
						continue
					}

					if tccommon.NeedProtect == 1 && int64(interval) < 30 {
						continue
					}

					ee := service.DeleteLayer4Listener(ctx, *tcpListener.ListenerId, proxyId, *tcpListener.Protocol)
					if ee != nil {
						continue
					}

				}
				udpListeners, err := service.DescribeUDPListeners(ctx, &proxyIdTmp, nil, nil, nil)
				if err != nil {
					return err
				}

				for _, udpListener := range udpListeners {
					instanceName := *udpListener.ListenerName

					now := time.Now()
					createTime := time.Unix(int64(*udpListener.CreateTime), 0)
					interval := now.Sub(createTime).Minutes()

					if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
						continue
					}

					if tccommon.NeedProtect == 1 && int64(interval) < 30 {
						continue
					}

					ee := service.DeleteLayer4Listener(ctx, *udpListener.ListenerId, proxyId, *udpListener.Protocol)
					if ee != nil {
						continue
					}

				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer4ListenerDestroy(id, "TCP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9090"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer4_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_TcpDomain(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer4ListenerDestroy(id, "TCP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerTcpDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9091"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_update(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer4ListenerDestroy(id, "TCP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9092"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "true"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "interval", "5"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerUpdateNameAndHealthConfigAndScheduler,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-listener-new"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "scheduler", "wrr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "interval", "11"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "connect_timeout", "10"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerUpdateNoHealthCheck,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "false"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerTcpUpdateRealserverSet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_udp_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer4ListenerDestroy(id, "UDP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUdp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9093"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer4_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_udpDomain(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer4ListenerDestroy(id, "UDP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUdpDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9095"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "DOMAIN"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer4Listener_udpUpdate(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer4ListenerDestroy(id, "UDP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer4ListenerUdpUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udp-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "port", "9094"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "scheduler", "rr"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_type", "IP"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "health_check", "false"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "realserver_bind_set.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer4_listener.foo", "proxy_id"),
				),
			},
			{
				Config: testAccGaapLayer4ListenerUdpUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer4ListenerExists("tencentcloud_gaap_layer4_listener.foo", id, "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer4_listener.foo", "name", "ci-test-gaap-4-udpListener-new"),
				),
			},
		},
	})
}

func testAccCheckGaapLayer4ListenerExists(n string, id *string, protocol string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no listener ID is set")
		}

		service := svcgaap.NewGaapService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		switch protocol {
		case "TCP":
			listeners, err := service.DescribeTCPListeners(context.TODO(), nil, &rs.Primary.ID, nil, nil)
			if err != nil {
				return err
			}

			for _, l := range listeners {
				if l.ListenerId == nil {
					return errors.New("listener id is nil")
				}
				if rs.Primary.ID == *l.ListenerId {
					*id = *l.ListenerId
					break
				}
			}

		case "UDP":
			listeners, err := service.DescribeUDPListeners(context.TODO(), nil, &rs.Primary.ID, nil, nil)
			if err != nil {
				return err
			}

			for _, l := range listeners {
				if l.ListenerId == nil {
					return errors.New("listener id is nil")
				}
				if rs.Primary.ID == *l.ListenerId {
					*id = *l.ListenerId
					break
				}
			}
		}

		if id == nil {
			return fmt.Errorf("listener not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckGaapLayer4ListenerDestroy(id *string, protocol string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcgaap.NewGaapService(client)

		switch protocol {
		case "TCP":
			listeners, err := service.DescribeTCPListeners(context.TODO(), nil, id, nil, nil)
			if err != nil {
				if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == svcgaap.GAAPResourceNotFound || (sdkError.Code == "InvalidParameter" && sdkError.Message == "ListenerId") {
						return nil
					}
				}

				return err
			}
			if len(listeners) > 0 {
				return errors.New("listener still exists")
			}

		case "UDP":
			listeners, err := service.DescribeUDPListeners(context.TODO(), nil, id, nil, nil)
			if err != nil {
				return err
			}
			if len(listeners) > 0 {
				return errors.New("listener still exists")
			}
		}

		return nil
	}
}

var testAccGaapLayer4ListenerBasic = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 9090
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1, tcacctest.DefaultGaapRealserverIpId2, tcacctest.DefaultGaapRealserverIp2)

var testAccGaapLayer4ListenerTcpDomain = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 9091
  realserver_type = "DOMAIN"
  proxy_id        = "%s"
  health_check    = true

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverDomainId1, tcacctest.DefaultGaapRealserverDomain1, tcacctest.DefaultGaapRealserverDomainId2, tcacctest.DefaultGaapRealserverDomain2)

var testAccGaapLayer4ListenerUpdateBasic = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 9092
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1, tcacctest.DefaultGaapRealserverIpId2, tcacctest.DefaultGaapRealserverIp2)

var testAccGaapLayer4ListenerUpdateNameAndHealthConfigAndScheduler = fmt.Sprintf(`
resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener-new"
  port            = 9092
  scheduler       = "wrr"
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = true
  interval      = 11
  connect_timeout = 10

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1, tcacctest.DefaultGaapRealserverIpId2, tcacctest.DefaultGaapRealserverIp2)

var testAccGaapLayer4ListenerUpdateNoHealthCheck = fmt.Sprintf(`
resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener-new"
  port            = 9092
  scheduler       = "wrr"
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = false

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1)

var testAccGaapLayer4ListenerTcpUpdateRealserverSet = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "TCP"
  name            = "ci-test-gaap-4-listener"
  port            = 9092
  scheduler       = "wrr"
  realserver_type = "IP"
  proxy_id        = "%s"
  health_check    = false

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1)

var testAccGaapLayer4ListenerUdp = fmt.Sprintf(`
resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "UDP"
  name            = "ci-test-gaap-4-udp-listener"
  port            = 9093
  realserver_type = "IP"
  proxy_id        = "%s"

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1, tcacctest.DefaultGaapRealserverIpId2, tcacctest.DefaultGaapRealserverIp2)

var testAccGaapLayer4ListenerUdpDomain = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "UDP"
  name            = "ci-test-gaap-4-udp-listener"
  port            = 9095
  realserver_type = "DOMAIN"
  proxy_id        = "%s"

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverDomainId1, tcacctest.DefaultGaapRealserverDomain1, tcacctest.DefaultGaapRealserverDomainId2, tcacctest.DefaultGaapRealserverDomain2)

var testAccGaapLayer4ListenerUdpUpdateBasic = fmt.Sprintf(`
resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "UDP"
  name            = "ci-test-gaap-4-udp-listener"
  port            = 9094
  realserver_type = "IP"
  proxy_id        = "%s"

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1, tcacctest.DefaultGaapRealserverIpId2, tcacctest.DefaultGaapRealserverIp2)

var testAccGaapLayer4ListenerUdpUpdate = fmt.Sprintf(`

resource tencentcloud_gaap_layer4_listener "foo" {
  protocol        = "UDP"
  name            = "ci-test-gaap-4-udpListener-new"
  port            = 9094
  realserver_type = "IP"
  proxy_id        = "%s"

  realserver_bind_set {
    id   = "%s"
    ip   = "%s"
    port = 80
  }

  realserver_bind_set {
    id     = "%s"
    ip     = "%s"
    port   = 80
  }
}
`, tcacctest.DefaultGaapProxyId, tcacctest.DefaultGaapRealserverIpId1, tcacctest.DefaultGaapRealserverIp1, tcacctest.DefaultGaapRealserverIpId2, tcacctest.DefaultGaapRealserverIp2)
