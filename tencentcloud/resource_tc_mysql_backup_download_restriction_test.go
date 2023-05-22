package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMysqlBackupDownloadRestrictionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlBackupDownloadRestriction,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_backup_download_restriction.backup_download_restriction", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_backup_download_restriction.backup_download_restriction",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlBackupDownloadRestriction = `

resource "tencentcloud_mysql_backup_download_restriction" "backup_download_restriction" {
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
