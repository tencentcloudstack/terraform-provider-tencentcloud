package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqRabbitmqUserResource_basic -v
func TestAccTencentCloudNeedFixTdmqRabbitmqUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "user"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "max_connections"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "max_channels"),
				),
			},
			{
				Config: testAccTdmqRabbitmqUserUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "user"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "password"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "max_connections"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_user.rabbitmq_user", "max_channels"),
				),
			},
		},
	})
}

const testAccTdmqRabbitmqUser = `
resource "tencentcloud_tdmq_rabbitmq_user" "rabbitmq_user" {
  instance_id     = "amqp-kzbe8p3n"
  user            = "keep-user"
  password        = "asdf1234"
  description     = "test user"
  tags            = ["management", "monitoring", "test"]
  max_connections = 3
  max_channels    = 3
}
`

const testAccTdmqRabbitmqUserUpdate = `
resource "tencentcloud_tdmq_rabbitmq_user" "rabbitmq_user" {
  instance_id     = "amqp-kzbe8p3n"
  user            = "keep-user"
  password        = "asdf1234"
  description     = "test user update"
  tags            = ["management", "monitoring"]
  max_connections = 10
  max_channels    = 10
}
`
