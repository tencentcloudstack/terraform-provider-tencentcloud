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

func TestAccTencentCloudAlbServerAttachment_tcp(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbServerAttachment_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbServerAttachmentExists("tencentcloud_alb_server_attachment.attachment_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.attachment_tcp", "loadbalancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.attachment_tcp", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.attachment_tcp", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.attachment_tcp", "backends.#", "1"),
				),
			}, {
				Config: testAccAlbServerAttachment_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbServerAttachmentExists("tencentcloud_alb_server_attachment.attachment_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.attachment_tcp", "loadbalancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.attachment_tcp", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.attachment_tcp", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.attachment_tcp", "backends.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudAlbServerAttachment_http(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbServerAttachment_http,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbServerAttachmentExists("tencentcloud_alb_server_attachment.attachment_http"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.attachment_http", "loadbalancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_server_attachment.attachment_http", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.attachment_http", "protocol_type", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_server_attachment.attachment_http", "backends.#", "1"),
				),
			},
		},
	})
}

func testAccCheckAlbServerAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

const testAccAlbServerAttachment_tcp = `
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
  clb_name     = "tf-clb-attachment-tcp"
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

resource "tencentcloud_alb_server_attachment" "attachment_tcp" {
  loadbalancer_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"

  backends {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 23
    weight      = 10
  }
}
`

const testAccAlbServerAttachment_tcp_update = `
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
  clb_name     = "tf-clb-attachment-tcp"
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

resource "tencentcloud_alb_server_attachment" "attachment_tcp" {
  loadbalancer_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id          = "${tencentcloud_clb_listener.listener_basic.id}"

  backends {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 80
    weight      = 30
  }
}
`

const testAccAlbServerAttachment_http = `
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
  clb_name     = "tf-clb-attachment-http"
  vpc_id       = "${tencentcloud_vpc.foo.id}"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id               = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTP"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_basic.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_alb_server_attachment" "attachment_http" {
  loadbalancer_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id          = "${tencentcloud_clb_listener.listener_basic.id}"
  location_id          = "${tencentcloud_clb_listener_rule.rule_basic.id}"

  backends {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 88
    weight      = 20
  }
}
`
