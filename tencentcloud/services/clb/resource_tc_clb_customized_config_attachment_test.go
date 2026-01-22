package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClbCustomizedConfigAttachment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbLogsetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbCustomizedConfigAttachmentCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbLogsetExists("tencentcloud_clb_customized_config_attachment.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_customized_config_attachment.example", "config_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_customized_config_attachment.example", "bind_list"),
				),
			},
			{
				Config: testAccClbCustomizedConfigAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbLogsetExists("tencentcloud_clb_customized_config_attachment.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_customized_config_attachment.example", "config_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_customized_config_attachment.example", "bind_list"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_customized_config_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbCustomizedConfigAttachmentCreate = `
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "SERVER"
}

resource "tencentcloud_clb_customized_config_attachment" "example" {
  config_id = tencentcloud_clb_customized_config_v2.example.config_id
  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-9bsa90io"
    domain           = "demo1.com"
  }

  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-qfljudr4"
    domain           = "demo2.com"
  }
}
`

const testAccClbCustomizedConfigAttachmentUpdate = `
resource "tencentcloud_clb_customized_config_v2" "example" {
  config_content = "client_max_body_size 224M;\r\nclient_body_timeout 60s;"
  config_name    = "tf-example"
  config_type    = "SERVER"
}

resource "tencentcloud_clb_customized_config_attachment" "example" {
  config_id = tencentcloud_clb_customized_config_v2.example.config_id
  bind_list {
    load_balancer_id = "lb-g1miv1ok"
    listener_id      = "lbl-9bsa90io"
    domain           = "demo1.com"
  }
}
`
