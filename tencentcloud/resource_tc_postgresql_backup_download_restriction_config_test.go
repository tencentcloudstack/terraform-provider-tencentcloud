package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const TestAccPostgresqlBackupDownloadRestrictionObject = "tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config"

func TestAccTencentCloudPostgresqlBackupDownloadRestrictionConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig_NONE,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "restriction_type", "NONE"),
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

func TestAccTencentCloudPostgresqlBackupDownloadRestrictionConfigResource_CUSTOMIZE(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig_CUSTOMIZE_ALLOW,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "restriction_type", "CUSTOMIZE"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_restriction_effect", "ALLOW"),
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_id_set.#"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_restriction_effect", "ALLOW"),
					resource.TestCheckTypeSetElemAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_set.*", "192.168.1.1"),
				),
			},
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig_CUSTOMIZE_DENY,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "restriction_type", "CUSTOMIZE"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_restriction_effect", "DENY"),
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "vpc_id_set.#"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_restriction_effect", "DENY"),
					resource.TestCheckTypeSetElemAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "ip_set.*", "192.168.0.0"),
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

func TestAccTencentCloudPostgresqlBackupDownloadRestrictionConfigResource_INTRANET(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-guangzhou")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlBackupDownloadRestrictionConfig_INTRANET,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(TestAccPostgresqlBackupDownloadRestrictionObject, "id"),
					resource.TestCheckResourceAttr(TestAccPostgresqlBackupDownloadRestrictionObject, "restriction_type", "INTRANET"),
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

const testAccPostgresqlBackupDownloadRestrictionConfig_NONE = `

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "NONE"
}

`

const testAccPostgresqlBackupDownloadRestrictionConfig_CUSTOMIZE_ALLOW = defaultVpcVariable + `
resource "tencentcloud_vpc" "pg_vpc1" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_vpc" "pg_vpc2" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "CUSTOMIZE"
  vpc_restriction_effect = "ALLOW"
  vpc_id_set = [tencentcloud_vpc.pg_vpc1.id]
  ip_restriction_effect = "ALLOW"
  ip_set = ["192.168.1.1"]
}

`

const testAccPostgresqlBackupDownloadRestrictionConfig_CUSTOMIZE_DENY = defaultVpcVariable + `
resource "tencentcloud_vpc" "pg_vpc1" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_vpc" "pg_vpc2" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "CUSTOMIZE"
  vpc_restriction_effect = "DENY"
  vpc_id_set = [tencentcloud_vpc.pg_vpc2.id]
  ip_restriction_effect = "DENY"
  ip_set = ["192.168.0.0"]
}

`

const testAccPostgresqlBackupDownloadRestrictionConfig_INTRANET = `

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "INTRANET"
}

`
