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
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_attachments.attachments", "attachment_list.0.rule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_attachments.attachments", "attachment_list.0.targets.#", "1"),
				),
			},
		},
	})
}

const testAccClbServerAttachmentsDataSource = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  name              = "example"
  availability_zone = "${var.availability_zone}"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

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
  instance_name              = "example"
  availability_zone          = "${var.availability_zone}"
  image_id                   = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type              = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type           = "CLOUD_PREMIUM"
  internet_max_bandwidth_out = 0
  vpc_id                     = "${tencentcloud_vpc.foo.id}"
  subnet_id                  = "${tencentcloud_subnet.foo.id}"
}

resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
  vpc_id       = "${tencentcloud_vpc.foo.id}"
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
  rule_id = "#${tencentcloud_clb_listener.listener_basic.id}"
  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 23
    weight      = 10
  }
}

data "tencentcloud_clb_attachments" "attachments" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  rule_id = "${tencentcloud_clb_attachment.attachment_tcp.id}"
}
`
