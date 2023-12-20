package as_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixAsExecuteScalingPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsExecuteScalingPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_execute_scaling_policy.execute_scaling_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_as_execute_scaling_policy.execute_scaling_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAsExecuteScalingPolicy = `

resource "tencentcloud_as_execute_scaling_policy" "execute_scaling_policy" {
  auto_scaling_policy_id = "asp-519acdug"
  honor_cooldown = false
  trigger_source = "API"
}

`
