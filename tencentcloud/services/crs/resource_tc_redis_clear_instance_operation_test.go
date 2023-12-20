package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisClearInstanceOperationResource_basic -v
func TestAccTencentCloudRedisClearInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisClearInstanceOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_clear_instance_operation.clear_instance_operation", "id")),
			},
		},
	})
}

const testAccRedisClearInstanceOperation = testAccRedisInstanceCluster + `

resource "tencentcloud_redis_clear_instance_operation" "clear_instance_operation" {
  instance_id = tencentcloud_redis_instance.redis_cluster.id
  password = "AAA123456BBB"
}

`
