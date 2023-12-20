package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisBackupOperationResource_basic -v
func TestAccTencentCloudRedisBackupOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_backup_operation.backup_operation", "id"),
				),
			},
		},
	})
}

const testAccRedisBackupOperationVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultCrsInstanceId + `"
}
`

const testAccRedisBackupOperation = testAccRedisBackupOperationVar + `

resource "tencentcloud_redis_backup_operation" "backup_operation" {
	instance_id = var.instance_id
	remark = "backup test"
	storage_days = 7
}

`
