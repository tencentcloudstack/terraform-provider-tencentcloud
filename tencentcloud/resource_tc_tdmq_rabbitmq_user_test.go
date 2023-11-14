package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRabbitmqUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqUser,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_user.rabbitmq_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqUser = `

resource "tencentcloud_tdmq_rabbitmq_user" "rabbitmq_user" {
  instance_id = ""
  user = ""
  password = ""
  description = ""
  max_connections = 
  max_channels = 
}

`
