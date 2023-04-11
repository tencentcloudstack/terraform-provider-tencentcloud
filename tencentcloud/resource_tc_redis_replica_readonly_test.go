package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRedisReplicaReadonlyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisReplicaReadonly,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_replica_readonly.replica_readonly", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_replica_readonly.replica_readonly",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisReplicaReadonly = `

resource "tencentcloud_redis_replica_readonly" "replica_readonly" {
  instance_id = "crs-c1nl9rpv"
  readonly_policy = 
}

`
