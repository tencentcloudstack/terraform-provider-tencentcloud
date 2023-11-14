package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqRoleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqRole,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_role.role", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_role.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqRole = `

resource "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = &lt;nil&gt;
  remark = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
      }

`
