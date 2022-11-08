package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdmqRabbitmqRole_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_role.rabbitmq_role", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_role.rabbitmqRole",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqRole = `

resource "tencentcloud_tdmq_rabbitmq_role" "rabbitmq_role" {
  role_name = ""
  cluster_id = ""
  remark = ""
}

`
