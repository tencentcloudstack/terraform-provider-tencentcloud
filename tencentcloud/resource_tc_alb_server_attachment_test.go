package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudAlbServerAttachmentTcp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAlbServerAttachmentTcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbServerAttachmentExists("tencentcloud_alb_server_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.foo", "loadbalancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.foo", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.foo", "backends.#", "1"),
				),
			}, {
				Config: testAccTencentCloudAlbServerAttachmentTcpUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbServerAttachmentExists("tencentcloud_alb_server_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.foo", "loadbalancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.foo", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.foo", "backends.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudAlbServerAttachmentHttp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAlbServerAttachmentHttp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbServerAttachmentExists("tencentcloud_alb_server_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.foo", "loadbalancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.foo", "protocol_type", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.foo", "backends.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAlbServerAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_alb_server_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		items := strings.Split(rs.Primary.ID, ":")
		if len(items) != 3 {
			return fmt.Errorf("id %s of resource.tencentcloud_alb_server_attachment is wrong", rs.Primary.ID)
		}
		clbId := items[0]
		listenerId := items[1]
		locationId := items[1]
		_, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err == nil {
			return fmt.Errorf("clb ServerAttachment still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAlbServerAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
		items := strings.Split(rs.Primary.ID, ":")
		if len(items) != 3 {
			return fmt.Errorf("id %s of resource.tencentcloud_alb_server_attachment is wrong", rs.Primary.ID)
		}
		clbId := items[0]
		listenerId := items[1]
		locationId := items[1]
		_, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccTencentCloudAlbServerAttachmentTcp = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.vpc_id
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id                     = tencentcloud_clb_instance.foo.id
  listener_name              = var.instance_name
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_alb_server_attachment" "foo" {
  loadbalancer_id = tencentcloud_clb_instance.foo.id
  listener_id     = tencentcloud_clb_listener.foo.id

  backends {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 10
  }
}
`

const testAccTencentCloudAlbServerAttachmentTcpUpdate = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.vpc_id
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id                     = tencentcloud_clb_instance.foo.id
  listener_name              = var.instance_name
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}

resource "tencentcloud_alb_server_attachment" "foo" {
  loadbalancer_id = tencentcloud_clb_instance.foo.id
  listener_id     = tencentcloud_clb_listener.foo.id

  backends {
    instance_id = tencentcloud_instance.default.id
    port        = 80
    weight      = 30
  }
}
`

const testAccTencentCloudAlbServerAttachmentHttp = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.vpc_id
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id               = tencentcloud_clb_instance.foo.id
  listener_name        = var.instance_name
  port                 = 77
  protocol             = "HTTP"
}

resource "tencentcloud_clb_listener_rule" "foo" {
  clb_id              = tencentcloud_clb_instance.foo.id
  listener_id         = tencentcloud_clb_listener.foo.id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_alb_server_attachment" "foo" {
  loadbalancer_id      = tencentcloud_clb_instance.foo.id
  listener_id          = tencentcloud_clb_listener.foo.id
  location_id          = tencentcloud_clb_listener_rule.foo.id

  backends {
    instance_id = tencentcloud_instance.default.id
    port        = 88
    weight      = 20
  }
}
`
