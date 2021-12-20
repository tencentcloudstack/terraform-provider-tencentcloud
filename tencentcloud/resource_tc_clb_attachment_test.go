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

func TestAccTencentCloudClbServerAttachment_tcp(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbServerAttachment_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			}, {
				Config: testAccClbServerAttachment_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "protocol_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbServerAttachment_http(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbServerAttachment_http, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_attachment.foo", "listener_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "protocol_type", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_attachment.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClbServerAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
			return fmt.Errorf("[CHECK][CLB attachment][Destroy] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", rs.Primary.ID)
		}
		locationId := items[0]
		listenerId := items[1]
		clbId := items[2]
		instance, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if (instance != nil && !(len(instance.Targets) == 0 && locationId == "") && !(len(instance.Rules) == 0 && locationId != "")) && err == nil {
			return fmt.Errorf("[CHECK][CLB attachment][Destroy] check: CLB Attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbServerAttachmentTargetGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbServerAttachment_multiple, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccClbServerAttachment_multiple_update, defaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttr("tencentcloud_clb_attachment.foo", "targets.#", "1"),
				),
			},
		},
	})
}

func testAccCheckClbServerAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB attachment][Exists] check: CLB Attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB attachment][Exists] check: CLB Attachment id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) != 3 {
			return fmt.Errorf("[CHECK][CLB attachment][Exists] check: id %s of resource.tencentcloud_clb_attachment is not match loc-xxx#lbl-xxx#lb-xxx", rs.Primary.ID)
		}
		locationId := items[0]
		listenerId := items[1]
		clbId := items[2]
		instance, err := clbService.DescribeAttachmentByPara(ctx, clbId, listenerId, locationId)
		if err != nil {
			return err
		}
		if instance == nil || (len(instance.Targets) == 0 && locationId == "") || (len(instance.Rules) == 0 && locationId != "") {
			return fmt.Errorf("[CHECK][CLB attachment][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbServerAttachment_tcp = instanceCommonTestCase + `
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

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = tencentcloud_clb_instance.foo.id
  listener_id = tencentcloud_clb_listener.foo.id

  targets {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 10
  }
}
`

const testAccClbServerAttachment_tcp_update = instanceCommonTestCase + `
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

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = tencentcloud_clb_instance.foo.id
  listener_id = tencentcloud_clb_listener.foo.id

  targets {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 50
  }
}
`

const testAccClbServerAttachment_http = instanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.vpc_id
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id               = tencentcloud_clb_instance.foo.id
  listener_name        = var.instance_name
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
}

resource "tencentcloud_clb_listener_rule" "foo" {
  clb_id              = tencentcloud_clb_instance.foo.id
  listener_id         = tencentcloud_clb_listener.foo.id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = tencentcloud_clb_instance.foo.id
  listener_id = tencentcloud_clb_listener.foo.id
  rule_id     = tencentcloud_clb_listener_rule.foo.id

  targets {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 10
  }
}
`

const testAccClbServerAttachment_multiple = instanceCommonTestCase + `
resource "tencentcloud_instance" "update" {
  instance_name              = var.instance_name_update
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = var.vpc_id
  subnet_id                  = var.subnet_id
}

resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.vpc_id
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id               = tencentcloud_clb_instance.foo.id
  listener_name        = var.instance_name
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
}

resource "tencentcloud_clb_listener_rule" "foo" {
  clb_id              = tencentcloud_clb_instance.foo.id
  listener_id         = tencentcloud_clb_listener.foo.id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = tencentcloud_clb_instance.foo.id
  listener_id = tencentcloud_clb_listener.foo.id
  rule_id     = tencentcloud_clb_listener_rule.foo.id

  targets {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 10
  }
  targets {
    instance_id = tencentcloud_instance.update.id
    port        = 24
    weight      = 10
  }
}
`

const testAccClbServerAttachment_multiple_update = instanceCommonTestCase + `

resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.vpc_id
}

resource "tencentcloud_clb_listener" "foo" {
  clb_id               = tencentcloud_clb_instance.foo.id
  listener_name        = var.instance_name
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
}

resource "tencentcloud_clb_listener_rule" "foo" {
  clb_id              = tencentcloud_clb_instance.foo.id
  listener_id         = tencentcloud_clb_listener.foo.id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = tencentcloud_clb_instance.foo.id
  listener_id = tencentcloud_clb_listener.foo.id
  rule_id     = tencentcloud_clb_listener_rule.foo.id

  targets {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 10
  }
}
`
