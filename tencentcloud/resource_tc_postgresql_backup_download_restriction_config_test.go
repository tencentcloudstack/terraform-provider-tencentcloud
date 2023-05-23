package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresqlBackupDownloadRestrictionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlBackupDownloadRestrictionConfig = `

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = ""
  vpc_restriction_effect = ""
  vpc_id_set = 
  ip_restriction_effect = ""
  ip_set = 
}

`
