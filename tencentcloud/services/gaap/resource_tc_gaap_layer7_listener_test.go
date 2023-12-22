package gaap_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_gaap_layer7_listener
	resource.AddTestSweepers("tencentcloud_gaap_layer7_listener", &resource.Sweeper{
		Name: "tencentcloud_gaap_layer7_listener",
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
				httpListeners, err := service.DescribeHTTPListeners(ctx, &proxyIdTmp, nil, nil, nil)
				if err != nil {
					return err
				}
				for _, httpListener := range httpListeners {
					instanceName := *httpListener.ListenerName

					now := time.Now()
					createTime := time.Unix(int64(*httpListener.CreateTime), 0)
					interval := now.Sub(createTime).Minutes()

					if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
						continue
					}

					if tccommon.NeedProtect == 1 && int64(interval) < 30 {
						continue
					}

					ee := service.DeleteLayer7Listener(ctx, *httpListener.ListenerId, proxyId, *httpListener.Protocol)
					if ee != nil {
						continue
					}

				}
				httpsListeners, err := service.DescribeHTTPSListeners(ctx, &proxyIdTmp, nil, nil, nil)
				if err != nil {
					return err
				}

				for _, httpsListener := range httpsListeners {
					instanceName := *httpsListener.ListenerName

					now := time.Now()
					createTime := time.Unix(int64(*httpsListener.CreateTime), 0)
					interval := now.Sub(createTime).Minutes()

					if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
						continue
					}

					if tccommon.NeedProtect == 1 && int64(interval) < 30 {
						continue
					}

					ee := service.DeleteLayer7Listener(ctx, *httpsListener.ListenerId, proxyId, *httpsListener.Protocol)
					if ee != nil {
						continue
					}

				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_basic(t *testing.T) {
	t.Parallel()
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, "HTTP"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8080"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener-new"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer7_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_https_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8081"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener-new"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer7_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsTwoWayAuthentication(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsTwoWayAuthentication,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8082"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsForwardHttps(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsForwardHttps,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8083"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsPolyClientCertificateIds(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsPolyClientCertificateIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8084"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpsPolyClientCertificateIdsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_layer7_listener.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapLayer7Listener_httpsCcToPoly(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapLayer7ListenerDestroy(id, "HTTPS"),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapLayer7ListenerHttpsCcToPolyOld,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "name", "ci-test-gaap-l7-listener"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "port", "8085"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "forward_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "auth_type", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "proxy_id"),
				),
			},
			{
				Config: testAccGaapLayer7ListenerHttpsCcToPolyNew,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapLayer7ListenerExists("tencentcloud_gaap_layer7_listener.foo", id, "HTTPS"),
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_layer7_listener.foo", "client_certificate_id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_layer7_listener.foo", "client_certificate_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGaapLayer7ListenerExists(n string, id *string, protocol string) resource.TestCheckFunc {
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
		case "HTTP":
			listeners, err := service.DescribeHTTPListeners(context.TODO(), nil, &rs.Primary.ID, nil, nil)
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

		case "HTTPS":
			listeners, err := service.DescribeHTTPSListeners(context.TODO(), nil, &rs.Primary.ID, nil, nil)
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

func testAccCheckGaapLayer7ListenerDestroy(id *string, protocol string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcgaap.NewGaapService(client)

		switch protocol {
		case "HTTP":
			listeners, err := service.DescribeHTTPListeners(context.TODO(), nil, id, nil, nil)
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

		case "HTTPS":
			listeners, err := service.DescribeHTTPSListeners(context.TODO(), nil, id, nil, nil)
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
		}

		return nil
	}
}

var testAccGaapLayer7ListenerBasic = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 8080
  proxy_id = "%s"
}
`, tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpUpdateName = fmt.Sprintf(`
resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener-new"
  port     = 8080
  proxy_id = "%s"
}
`, tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttps = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 8081
  proxy_id         = "%s"
  certificate_id   = tencentcloud_gaap_certificate.foo.id
  forward_protocol = "HTTP"
  auth_type        = 0
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsUpdate = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener-new"
  port             = 8081
  proxy_id         = "%s"
  certificate_id   = tencentcloud_gaap_certificate.bar.id
  forward_protocol = "HTTP"
  auth_type        = 0
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsTwoWayAuthentication = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 8082
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  forward_protocol      = "HTTP"
  auth_type             = 1
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF", tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsForwardHttps = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol         = "HTTPS"
  name             = "ci-test-gaap-l7-listener"
  port             = 8083
  proxy_id         = "%s"
  certificate_id   = tencentcloud_gaap_certificate.foo.id
  forward_protocol = "HTTPS"
  auth_type        = 0
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF", tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsPolyClientCertificateIds = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 8084
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  forward_protocol            = "HTTP"
  auth_type                   = 1
  client_certificate_ids = [tencentcloud_gaap_certificate.bar.id]
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF", tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsPolyClientCertificateIdsUpdate = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "client1" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "client2" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 8084
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  forward_protocol            = "HTTP"
  auth_type                   = 1
  client_certificate_ids = [tencentcloud_gaap_certificate.client1.id, tencentcloud_gaap_certificate.client2.id]
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF",
	tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsCcToPolyOld = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol              = "HTTPS"
  name                  = "ci-test-gaap-l7-listener"
  port                  = 8085
  proxy_id              = "%s"
  certificate_id        = tencentcloud_gaap_certificate.foo.id
  forward_protocol      = "HTTP"
  auth_type             = 1
  client_certificate_id = tencentcloud_gaap_certificate.bar.id
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF", tcacctest.DefaultGaapProxyId)

var testAccGaapLayer7ListenerHttpsCcToPolyNew = fmt.Sprintf(`
resource tencentcloud_gaap_certificate "foo" {
  type    = "SERVER"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_certificate "bar" {
  type    = "CLIENT"
  content = %s
  key     = %s
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol                    = "HTTPS"
  name                        = "ci-test-gaap-l7-listener"
  port                        = 8085
  proxy_id                    = "%s"
  certificate_id              = tencentcloud_gaap_certificate.foo.id
  forward_protocol            = "HTTP"
  auth_type                   = 1
  client_certificate_ids = [tencentcloud_gaap_certificate.bar.id]
}

`, "<<EOF"+testAccGaapCertificateServerCert+"EOF", "<<EOF"+testAccGaapCertificateServerKey+"EOF",
	"<<EOF"+testAccGaapCertificateClientCA+"EOF", "<<EOF"+testAccGaapCertificateClientCAKey+"EOF", tcacctest.DefaultGaapProxyId)
