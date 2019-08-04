package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudClbServerAttachmentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbServerAttachmentsDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_server_attachment.server_attachment_tcp"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_server_attachments.attachments", "attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_server_attachments.attachments", "attachment_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_server_attachments.attachments", "attachment_list.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_server_attachments.attachments", "attachment_list.0.location_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_server_attachments.attachments", "attachment_list.0.targets.#", "1"),
				),
			},
		},
	})
}

const testAccClbServerAttachmentsDataSource = `
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
  location_id = "#${tencentcloud_clb_listener.listener_basic.id}"
  targets {
    instance_id = "ins-1flbqyp8"
    port        = 23
    weight      = 10
  }
}
data "tencentcloud_clb_server_attachments" "attachments" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  location_id = "${tencentcloud_clb_server_attachment.server_attachment_tcp.id}"
}
`
