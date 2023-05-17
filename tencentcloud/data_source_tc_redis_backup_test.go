package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisBackupDataSource_basic -v
func TestAccTencentCloudRedisBackupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_backup.backup"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.backup_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.backup_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.backup_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.file_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.full_backup"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.locked"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup.backup", "backup_set.0.status"),
				),
			},
		},
	})
}

const testAccRedisBackupDataSourceVar = `
variable "instance_id" {
	default = "` + defaultCrsInstanceId + `"
}
`

const testAccRedisBackupDataSource = testAccRedisBackupDataSourceVar + `

data "tencentcloud_redis_backup" "backup" {
	instance_id = var.instance_id
	# begin_time = "2023-04-07 19:50:40"
	# end_time = "2023-04-07 19:50:50"
	status = [2]
	instance_name = "Keep-terraform"
}

`
