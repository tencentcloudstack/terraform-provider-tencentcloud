package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  d_b_instance_id = ""
  backup_type = ""
  backup_id = ""
  u_r_l_expire_time = 
  backup_download_restriction {
		restriction_type = ""
		vpc_restriction_effect = ""
		vpc_id_set = 
		ip_restriction_effect = ""
		ip_set = 

  }
    tags = {
    "createdBy" = "terraform"
  }
}

`
