package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudClbListenerRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListenerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_basic", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_basic", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "domain", "abc.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "url", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "target_type", "TARGETGROUP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "forward_type", "HTTPS"),
				),
			},
			{
				Config: testAccClbListenerRule__basic_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_basic", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_basic", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "domain", "abc.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "url", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_basic", "forward_type", "HTTP"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener_rule.rule_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListenerRuleResource_full(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbListenerRule_full, tcacctest.DefaultSshCertificate, tcacctest.DefaultSshCertificateB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule_full"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_full", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_full", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "domain", "abc.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "url", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_unhealth_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_code", "31"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_domain", "abc.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "certificate_id", tcacctest.DefaultSshCertificateB),
				),
			}, {
				Config: fmt.Sprintf(testAccClbListenerRule_update, tcacctest.DefaultSshCertificate, tcacctest.DefaultSshCertificateB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule_full"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_full", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener_rule.rule_full", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "domain", "abcd.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "url", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_interval_time", "300"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_health_num", "6"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_unhealth_num", "6"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_method", "HEAD"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_code", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_domain", "abcd.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "health_check_http_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_full", "certificate_id", tcacctest.DefaultSshCertificateB),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener_rule.rule_full",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListenerRuleResource_oauth(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
			tcacctest.AccStepSetRegion(t, "ap-jakarta")
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListenerRule_oauth,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule_oauth"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_oauth", "oauth.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_oauth", "oauth.0.oauth_enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_oauth", "oauth.0.oauth_failure_status", "REJECT"),
				),
			},
			{
				Config: testAccClbListenerRule_oauthUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerRuleExists("tencentcloud_clb_listener_rule.rule_oauth"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_oauth", "oauth.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_oauth", "oauth.0.oauth_enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener_rule.rule_oauth", "oauth.0.oauth_failure_status", "BYPASS"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener_rule.rule_oauth",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClbListenerRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_listener_rule" {
			continue
		}
		resourceId := rs.Primary.ID
		items := strings.Split(resourceId, tccommon.FILED_SP)
		itemLength := len(items)
		locationId := items[itemLength-1]
		listenerId := rs.Primary.Attributes["listener_id"]
		clbId := rs.Primary.Attributes["clb_id"]
		//this function is not supported by api, need to be travelled
		filter := map[string]string{"rule_id": locationId, "listener_id": listenerId, "clb_id": clbId}
		rules, err := clbService.DescribeRulesByFilter(ctx, filter)
		if len(rules) > 0 && err == nil {
			return fmt.Errorf("[CHECK][CLB listener rule][Destroy] check: CLB listener rule still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbListenerRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB listener rule][Exists] check: CLB listener rule %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB listener rule][Exists] check: CLB listener rule id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceId := rs.Primary.ID
		items := strings.Split(resourceId, tccommon.FILED_SP)
		itemLength := len(items)
		locationId := items[itemLength-1]
		listenerId := rs.Primary.Attributes["listener_id"]
		clbId := rs.Primary.Attributes["clb_id"]
		filter := map[string]string{"rule_id": locationId, "listener_id": listenerId, "clb_id": clbId}
		rules, err := clbService.DescribeRulesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(rules) == 0 {
			return fmt.Errorf("[CHECK][CLB listener rule][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbListenerRule_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
  forward_type        = "HTTPS"
}
`

const testAccClbListenerRule__basic_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  forward_type        = "HTTP"
}
`

const testAccClbListenerRule_full = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-full"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id               = tencentcloud_clb_instance.clb_basic.id
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
  sni_switch = true
}

resource "tencentcloud_clb_listener_rule" "rule_full" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_id                = tencentcloud_clb_listener.listener_basic.listener_id
  domain                     = "abc.com"
  url                        = "/"
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_switch        = true
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_http_path     = "/"
  health_check_http_domain   = "abc.com"
  health_check_http_code     = "31"
  health_check_http_method   = "GET"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
}
`

const testAccClbListenerRule_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-full"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id               = tencentcloud_clb_instance.clb_basic.id
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
  sni_switch = true
}

resource "tencentcloud_clb_listener_rule" "rule_full" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_id                = tencentcloud_clb_listener.listener_basic.listener_id
  domain                     = "abcd.com"
  url                        = "/"
  session_expire_time        = 60
  scheduler                  = "WRR"
  health_check_switch        = true
  health_check_interval_time = 300
  health_check_health_num    = 6
  health_check_unhealth_num  = 6
  health_check_time_out      = 4
  health_check_type          = "TCP"
  health_check_http_path     = "/"
  health_check_http_domain   = "abcd.com"
  health_check_http_code     = "1"
  health_check_http_method   = "HEAD"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
}
`

const testAccClbListenerRule_oauth = `
resource "tencentcloud_clb_listener_rule" "rule_oauth" {
  clb_id              = "lb-az5cm2h7"
  listener_id         = "lbl-egzxfxgj"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
  forward_type        = "HTTPS"
  oauth {
    oauth_enable = true
    oauth_failure_status = "REJECT"
  }
}
`

const testAccClbListenerRule_oauthUpdate = `
resource "tencentcloud_clb_listener_rule" "rule_oauth" {
  clb_id              = "lb-az5cm2h7"
  listener_id         = "lbl-egzxfxgj"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
  forward_type        = "HTTPS"
  oauth {
    oauth_enable = false
    oauth_failure_status = "BYPASS"
  }
}
`
