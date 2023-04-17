package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisUpgradeVersionOperationResource_basic -v
func TestAccTencentCloudNeedFixRedisUpgradeVersionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeVersionOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_version_operation.upgrade_version_operation", "id")),
			},
		},
	})
}

const testAccRedisUpgradeVersionOperation = defaultVpcVariable + testAccRedisDefaultTemplate + `

resource "tencentcloud_redis_instance" "redis_instance_test" {
	availability_zone  = "ap-guangzhou-3"
	type_id            = 6
	password           = "test12345789"
	mem_size           = 8192
	name               = "terraform_test"
	port               = 6379
	redis_shard_num    = 1
	redis_replicas_num = 1
	params_template_id = var.redis_param_template
	vpc_id 			 = var.vpc_id
	subnet_id			 = var.subnet_id
  }

resource "tencentcloud_redis_upgrade_version_operation" "upgrade_version_operation" {
  instance_id = tencentcloud_redis_instance.redis_instance_test.id
  target_instance_type = "7"
  switch_option = 2
}

`
