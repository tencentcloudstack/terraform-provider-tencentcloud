package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafLogPostCkafkaFlowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafLogPostCkafkaFlow,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "ckafka_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "ckafka_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "brokers"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "compression"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "vip_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "topic"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "kafka_version"),
				),
			},
			{
				Config: testAccWafLogPostCkafkaFlowUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "ckafka_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "ckafka_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "brokers"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "compression"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "vip_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "topic"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.example", "kafka_version"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_log_post_ckafka_flow.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafLogPostCkafkaFlow = `
resource "tencentcloud_waf_log_post_ckafka_flow" "example" {
  ckafka_region = "ap-guangzhou"
  ckafka_id     = "ckafka-k9m5vwar"
  brokers       = "11.135.14.110:18737"
  compression   = "snappy"
  vip_type      = 2
  log_type      = 1
  topic         = "tf-example"
  kafka_version = "2.8.1"

  write_config {
    enable_body    = 0
    enable_bot     = 0
    enable_headers = 0
  }
}
`

const testAccWafLogPostCkafkaFlowUpdate = `
resource "tencentcloud_waf_log_post_ckafka_flow" "example" {
  ckafka_region = "ap-guangzhou"
  ckafka_id     = "ckafka-k9m5vwar"
  brokers       = "11.135.14.110:18737"
  compression   = "snappy"
  vip_type      = 2
  log_type      = 1
  topic         = "tf-example"
  kafka_version = "2.8.1"

  write_config {
    enable_body    = 1
    enable_bot     = 1
    enable_headers = 1
  }
}
`
