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
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.attachment_tcp"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_attachments.attachments", "attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_attachments.attachments", "attachment_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_attachments.attachments", "attachment_list.0.listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_attachments.attachments", "attachment_list.0.location_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_attachments.attachments", "attachment_list.0.targets.#", "1"),
				),
			},
		},
	})
}

const testAccClbServerAttachmentsDataSource = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name = "terraform_automation_test_kuruk"
  availability_zone = "ap-guangzhou-3"
  image_id      = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"

  system_disk_type = "CLOUD_PREMIUM"
  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 100
    delete_with_instance = true
  }
  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 100
    delete_with_instance = true
  }
  disable_security_service = true
  disable_monitor_service = true
}
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic1"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id                     = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name              = "listener_tcp"
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

resource "tencentcloud_clb_attachment" "attachment_tcp" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  location_id = "#${tencentcloud_clb_listener.listener_basic.id}"
  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 23
    weight      = 10
  }
}
data "tencentcloud_clb_attachments" "attachments" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  location_id = "${tencentcloud_clb_attachment.attachment_tcp.id}"
}
`
