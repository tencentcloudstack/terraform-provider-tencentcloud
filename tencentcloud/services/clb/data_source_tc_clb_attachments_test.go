package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbServerAttachmentsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbServerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceClbServerAttachments,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbServerAttachmentExists("tencentcloud_clb_attachment.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_attachments.foo", "attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_attachments.foo", "attachment_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_attachments.foo", "attachment_list.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_attachments.foo", "attachment_list.0.targets.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceClbServerAttachments = tcacctest.InstanceCommonTestCase + `
resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = var.instance_name
  vpc_id       = var.cvm_vpc_id
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
  listener_id = tencentcloud_clb_listener.foo.listener_id

  targets {
    instance_id = tencentcloud_instance.default.id
    port        = 23
    weight      = 10
  }
}

data "tencentcloud_clb_attachments" "foo" {
  clb_id      = tencentcloud_clb_instance.foo.id
  listener_id = tencentcloud_clb_attachment.foo.listener_id
}
`
