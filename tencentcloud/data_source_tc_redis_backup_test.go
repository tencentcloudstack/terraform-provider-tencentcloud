package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_backup.backup")),
			},
		},
	})
}

const testAccRedisBackupDataSource = `

data "tencentcloud_redis_backup" "backup" {
  instance_id = "crs-c1nl9rpv"
  begin_time = "2017-02-08 19:09:26"
  end_time = "2017-02-08 19:09:26"
  status = &lt;nil&gt;
  instance_name = &lt;nil&gt;
  }

`
