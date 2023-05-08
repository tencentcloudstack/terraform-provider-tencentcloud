package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisStartupInstanceOperationResource_basic -v
func TestAccTencentCloudNeedFixRedisStartupInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisStartupInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_startup_instance_operation.startup_instance_operation", "id")),
			},
		},
	})
}

const testAccRedisStartupInstanceOperation = testAccRedisInstanceCluster + `

resource "tencentcloud_redis_startup_instance_operation" "startup_instance_operation" {
  instance_id = tencentcloud_redis_instance.redis_cluster.id
}

`
