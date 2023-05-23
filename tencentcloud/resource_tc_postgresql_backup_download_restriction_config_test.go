package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestAccPostgresqlBackupDownloadRestrictionObject = "tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config"

func TestAccTencentCloudPostgresqlBackupDownloadRestrictionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig_none,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "restriction_type", "NONE"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_restriction_effect", "ALLOW"),
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_id_set.#"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_restriction_effect", "ALLOW"),
					resource.TestCheckTypeSetElemAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_set.*", "127.0.0.1"),
				),
			},
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig_INTRANET,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "restriction_type", "INTRANET"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_restriction_effect", "DENY"),
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_id_set.#"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_restriction_effect", "DENY"),
					resource.TestCheckTypeSetElemAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_set.*", "0.0.0.0"),
				),
			},
			{
				ResourceName:      TestAccPostgresqlBackupDownloadRestrictionObject,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlBackupDownloadRestrictionConfig_none = defaultVpcVariable + `
resource "tencentcloud_vpc" "pg_vpc" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "NONE"
  vpc_restriction_effect = "ALLOW"
  vpc_id_set = [tencentcloud_vpc.pg_vpc.id]
  ip_restriction_effect = "ALLOW"
  ip_set = [127.0.0.1]
}

`

const testAccPostgresqlBackupDownloadRestrictionConfig_INTRANET = `
resource "tencentcloud_vpc" "pg_vpc" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "INTRANET"
  vpc_restriction_effect = "DENY"
  vpc_id_set = [tencentcloud_vpc.pg_vpc.id]
  ip_restriction_effect = "DENY"
  ip_set = ["0.0.0.0"]
}

`
