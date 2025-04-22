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
		Steps: []resource.TestStep{{
			Config: testAccWafLogPostCkafkaFlow,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_log_post_ckafka_flow.waf_log_post_ckafka_flow", "id")),
		}, {
			ResourceName:      "tencentcloud_waf_log_post_ckafka_flow.waf_log_post_ckafka_flow",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWafLogPostCkafkaFlow = `

resource "tencentcloud_waf_log_post_ckafka_flow" "waf_log_post_ckafka_flow" {
}
`
