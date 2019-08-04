package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudClbServerAttachment_tcp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbServerAttachment_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_server_attachment.server_attachment_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_server_attachment.server_attachment_tcp", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_server_attachment.server_attachment_tcp", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.#", "1"),
					/*resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.0.port", "23"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.0.weight", "10"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.0.instance_id", "ins-1flbqyp8"),*/
				),
			}, {
				Config: testAccClbServerAttachment_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_server_attachment.server_attachment_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_server_attachment.server_attachment_tcp", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_server_attachment.server_attachment_tcp", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.#", "1"),
					/*resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.0.port", "23"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.0.weight", "10"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_tcp", "targets.0.instance_id", "ins-1flbqyp8"),*/
				),
			},
		},
	})
}

func TestAccTencentCloudClbServerAttachment_http(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbServerAttachment_http,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_server_attachment.server_attachment_http"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_server_attachment.server_attachment_http", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_server_attachment.server_attachment_http", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_http", "protocol_type", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_server_attachment.server_attachment_http", "targets.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_server_attachment.server_attachment_http",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClbServerAttachmentDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_server_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) != 3 {
			return fmt.Errorf("id of resource.tencentcloud_clb_server_attachment is wrong")
		}
		locationId := items[0]
		listenerId := items[1]
		clbId := items[2]
		_, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err == nil {
			return fmt.Errorf("clb ServerAttachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbServerAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("clb ServerAttachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("clb ServerAttachment id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) != 3 {
			return fmt.Errorf("id of resource.tencentcloud_clb_server_attachment is wrong")
		}
		locationId := items[0]
		listenerId := items[1]
		clbId := items[2]
		_, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccClbServerAttachment_tcp = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id                     = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = 1
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_clb_server_attachment" "server_attachment_tcp" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  targets {
    instance_id = "ins-1flbqyp8"
    port        = 23
    weight      = 10
  }
}
`

const testAccClbServerAttachment_tcp_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id                     = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = 1
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_clb_server_attachment" "server_attachment_tcp" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  targets {
    instance_id = "ins-1flbqyp8"
    port        = 23
    weight      = 50
  }
}
`

const testAccClbServerAttachment_http = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id               = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VfqcL1ME"

}
resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_basic.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}
resource "tencentcloud_clb_server_attachment" "server_attachment_http" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  location_id = "${tencentcloud_clb_listener_rule.rule_basic.id}"
  targets {
    instance_id = "ins-1flbqyp8"
    port        = 23
    weight      = 10
  }
}
`
