package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisReplicaReadonlyResource_basic -v
func TestAccTencentCloudRedisReplicaReadonlyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicaReadonly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replica_readonly.replica_readonly", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replica_readonly.replica_readonly", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replica_readonly.replica_readonly", "operate", "enable"),
				),
			},
			{
				Config: testAccRedisReplicaReadonlyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replica_readonly.replica_readonly", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_replica_readonly.replica_readonly", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_replica_readonly.replica_readonly", "operate", "disable"),
				),
			},
		},
	})
}

const testAccRedisReplicaReadonly = `

resource "tencentcloud_redis_replica_readonly" "replica_readonly" {
	instance_id = "crs-2yypjrnv"
	readonly_policy = ["master"]
	operate = "enable"
}

`

const testAccRedisReplicaReadonlyUpdate = `

resource "tencentcloud_redis_replica_readonly" "replica_readonly" {
	instance_id = "crs-2yypjrnv"
	readonly_policy = ["master"]
	operate = "disable"
}

`
