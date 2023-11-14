package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisBackupOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_backup_operation.backup_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_backup_operation.backup_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisBackupOperation = `

resource "tencentcloud_redis_backup_operation" "backup_operation" {
  instance_id = "crs-c1nl9rpv"
  remark = &lt;nil&gt;
  storage_days = 7
}

`
