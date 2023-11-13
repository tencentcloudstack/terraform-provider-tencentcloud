package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqEnvironmentRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqEnvironmentRole,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_environment_role.environment_role", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_environment_role.environment_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqEnvironmentRole = `

resource "tencentcloud_tdmq_rocketmq_environment_role" "environment_role" {
  environment_id = &lt;nil&gt;
  role_name = &lt;nil&gt;
  permissions = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
}

`
