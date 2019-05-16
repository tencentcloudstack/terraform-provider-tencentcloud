package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudRedisInstance_import(t *testing.T) {
	resourceName := "tencentcloud_redis_instance.redis_instance_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudRedisInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisInstanceBasic(),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
