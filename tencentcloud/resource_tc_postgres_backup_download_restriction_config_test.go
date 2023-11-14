package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresBackupDownloadRestrictionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresBackupDownloadRestrictionConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_backup_download_restriction_config.backup_download_restriction_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_backup_download_restriction_config.backup_download_restriction_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresBackupDownloadRestrictionConfig = `

resource "tencentcloud_postgres_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = ""
  vpc_restriction_effect = ""
  vpc_id_set = 
  ip_restriction_effect = ""
  ip_set = 
}

`
