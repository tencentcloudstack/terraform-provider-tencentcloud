package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlBackupDownloadUrlsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBackupDownloadUrlsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls")),
			},
		},
	})
}

const testAccPostgresqlBackupDownloadUrlsDataSource = `

data "tencentcloud_postgresql_backup_download_urls" "backup_download_urls" {
  db_instance_id = ""
  backup_type = ""
  backup_id = ""
  url_expire_time = 
  backup_download_restriction {
		restriction_type = ""
		vpc_restriction_effect = ""
		vpc_id_set = 
		ip_restriction_effect = ""
		ip_set = 

  }
}

`
