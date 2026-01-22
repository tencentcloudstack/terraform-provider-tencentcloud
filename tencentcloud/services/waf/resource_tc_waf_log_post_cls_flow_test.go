package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafLogPostClsFlowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafLogPostClsFlow,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "cls_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "logset_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "log_topic_name"),
				),
			},
			{
				Config: testAccWafLogPostClsFlowUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "cls_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "logset_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_cls_flow.example", "log_topic_name"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_log_post_cls_flow.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafLogPostClsFlow = `
resource "tencentcloud_waf_log_post_cls_flow" "example" {
  cls_region     = "ap-guangzhou"
  logset_name    = "waf_post_logset"
  log_type       = 1
  log_topic_name = "waf_post_logtopic"
}
`

const testAccWafLogPostClsFlowUpdate = `
resource "tencentcloud_waf_log_post_cls_flow" "example" {
  cls_region     = "ap-shanghai"
  logset_name    = "waf_post_logset_update"
  log_type       = 1
  log_topic_name = "waf_post_logtopic_update"
}
`
