package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisBackupDownloadInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupDownloadInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_backup_download_info.backup_download_info")),
			},
		},
	})
}

const testAccRedisBackupDownloadInfoDataSource = `

data "tencentcloud_redis_backup_download_info" "backup_download_info" {
  instance_id = "crs-c1nl9rpv"
  backup_id = "123456789-123456-1234567812"
  limit_type = "NoLimit"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol = "In"
  limit_vpc {
		region = "ap-guangzhou"
		vpc_list = 

  }
  limit_ip = 
  }

`
