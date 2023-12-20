package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisUpgradeMultiZoneOperationResource_basic -v
func TestAccTencentCloudRedisUpgradeMultiZoneOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeMultiZoneOperation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_multi_zone_operation.upgrade_multi_zone_operation", "id"),
				),
			},
		},
	})
}

// const testAccRedisUpgradeMultiZoneOperation = tcacctest.DefaultVpcVariable + `

// resource "tencentcloud_redis_instance" "redis_cluster" {
// 	availability_zone = "ap-guangzhou-3"
// 	type_id            = 7
// 	password           = "AAA123456BBB"
// 	mem_size           = 4096
// 	name               = "terraform_cluster"
// 	port               = 6379
// 	redis_shard_num    = 1
// 	redis_replicas_num = 2
// 	replica_zone_ids   = [100003, 100003]
// 	vpc_id 			   = var.vpc_id
// 	subnet_id		   = var.subnet_id
// }

// resource "tencentcloud_redis_upgrade_multi_zone_operation" "upgrade_multi_zone_operation" {
// 	instance_id = tencentcloud_redis_instance.redis_cluster.id
// 	upgrade_proxy_and_redis_server = false
// }

// `

func testAccRedisUpgradeMultiZoneOperation() string {
	return testAccRedisInstanceParamTemplate() + `
resource "tencentcloud_redis_upgrade_multi_zone_operation" "upgrade_multi_zone_operation" {
	instance_id = tencentcloud_redis_instance.redis_instance_test.id
	upgrade_proxy_and_redis_server = false
}`
}
