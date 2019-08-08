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
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.attachment_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.attachment_tcp", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.attachment_tcp", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.attachment_tcp", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.attachment_tcp", "targets.#", "1"),
				),
			}, {
				Config: testAccClbServerAttachment_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.attachment_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.attachment_tcp", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.attachment_tcp", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.attachment_tcp", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.attachment_tcp", "targets.#", "1"),
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
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.attachment_http"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.attachment_http", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.attachment_http", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.attachment_http", "protocol_type", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.attachment_http", "targets.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_attachment.attachment_http",
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
		if rs.Type != "tencentcloud_clb_attachment" {
			continue
		}
		time.Sleep(5 * time.Second)
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) != 3 {
			return fmt.Errorf("id of resource.tencentcloud_clb_attachment is wrong")
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
			return fmt.Errorf("id of resource.tencentcloud_clb_attachment is wrong")
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

  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 23
    weight      = 10
  }
}
`

const testAccClbServerAttachment_tcp_update = `
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

  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 23
    weight      = 50
  }
}
`

const testAccClbServerAttachment_http = `
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
  clb_id               = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VjANRdz8"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_basic.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "attachment_http" {
  clb_id      = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_id = "${tencentcloud_clb_listener.listener_basic.id}"
  rule_id = "${tencentcloud_clb_listener_rule.rule_basic.id}"

  targets {
    instance_id = "${tencentcloud_instance.foo.id}"
    port        = 23
    weight      = 10
  }
}
`
