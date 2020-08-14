package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceCkafkaUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCkafkaUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceCkafkaUser,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCkafkaUserExists("tencentcloud_ckafka_user.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_users.foo", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_users.foo", "user_list.0.account_name", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_users.foo", "user_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_users.foo", "user_list.0.update_time"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceCkafkaUser = `
resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = "ckafka-f9ife4zz"
  account_name = "test"
  password     = "test1234"
}

data "tencentcloud_ckafka_users" "foo" {
	instance_id  = tencentcloud_ckafka_user.foo.instance_id
	account_name = tencentcloud_ckafka_user.foo.account_name
}
`
