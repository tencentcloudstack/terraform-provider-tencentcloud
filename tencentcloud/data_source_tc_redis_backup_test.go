package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.backup_id", "641186639-8362913-1516672770"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.backup_size", "184"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.backup_type", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.end_time", "2023-04-07 03:57:36"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.expire_time", "2023-04-14 03:57:33"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.file_type", "RDB-Redis 5.0"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.full_backup", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.instance_name", "Keep-terraform"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.instance_type", "8"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.locked", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.remark", "xx"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.start_time", "2023-04-07 03:57:33"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup.backup", "backup_set.0.status", "2"),
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
	instance_id = "crs-jf4ico4v"
	begin_time = "2023-04-07 19:50:40"
	end_time = "2023-04-07 19:50:50"
	status = [2]
	instance_name = "Keep-terraform"
}

`
