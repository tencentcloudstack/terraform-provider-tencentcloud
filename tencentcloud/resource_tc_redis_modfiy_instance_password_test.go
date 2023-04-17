package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisBackupDownloadRestrictionResource_basic -v
func TestAccTencentCloudRedisModfiyInstancePasswordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisModfiyInstancePassword,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_modfiy_instance_password.modfiy_instance_password", "id"),
				),
			},
		},
	})
}

const testAccRedisModfiyInstancePassword = testAccRedisInstanceCluster + `

resource "tencentcloud_redis_modfiy_instance_password" "modfiy_instance_password" {
	instance_id  = tencentcloud_redis_instance.redis_cluster.id
	old_password = "AAA123456BBB"
	password 	 = "AAA123456BBBCCC"
}

`
