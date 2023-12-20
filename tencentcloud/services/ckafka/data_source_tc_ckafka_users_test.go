package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaUsersDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCkafkaUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceCkafkaUser,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCkafkaUserExists("tencentcloud_ckafka_user.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_users.foo", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_ckafka_users.foo", "user_list.0.account_name", "test1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_users.foo", "user_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_users.foo", "user_list.0.update_time"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceCkafkaUser = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = var.instance_id
  account_name = "test1"
  password     = "test1234"
}

data "tencentcloud_ckafka_users" "foo" {
	instance_id  = tencentcloud_ckafka_user.foo.instance_id
	account_name = tencentcloud_ckafka_user.foo.account_name
}
`
