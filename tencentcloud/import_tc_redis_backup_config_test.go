package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudRedisBackupConfig_import(t *testing.T) {
	resourceName := "tencentcloud_redis_backup_config.redis_backup_config"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudRedisBackupConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupConfig(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
