package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlBackupDownloadUrlsDataSource_basic(t *testing.T) {
	// t.Parallel()
	// loc, _ := time.LoadLocation("Asia/Chongqing")
	// startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	// endTime := time.Now().AddDate(0, 0, 1).In(loc).Format("2006-01-02 15:04:05")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-guangzhou")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				// Config: fmt.Sprintf(testAccPostgresqlBackupDownloadUrlsDataSource, startTime, endTime),
				Config: testAccPostgresqlBackupDownloadUrlsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "db_instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_type", "LogBackup"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "url_expire_time", "12"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_restriction.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_restriction.0.restriction_type", "NONE"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_restriction.0.vpc_restriction_effect", "ALLOW"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_restriction.0.vpc_id_set.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_restriction.0.ip_restriction_effect", "ALLOW"),
					resource.TestCheckResourceAttr("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_restriction.0.ip_set.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_backup_download_urls.backup_download_urls", "backup_download_url"),
				),
			},
		},
	})
}

const testAccPostgresqlBackupDownloadUrlsDataSource = OperationPresetPGSQL + defaultVpcSubnets + `
data "tencentcloud_postgresql_log_backups" "log_backups" {
	min_finish_time = ""
	max_finish_time = ""
	filters {
		  name = "db-instance-id"
		  values = [local.pgsql_id]
	}
	// order_by = "StartTime"
	order_by_type = "desc"
  }

data "tencentcloud_postgresql_backup_download_urls" "backup_download_urls" {
  db_instance_id = local.pgsql_id
  backup_type = "LogBackup"
  backup_id = data.tencentcloud_postgresql_log_backups.log_backups.log_backup_set.0.id
//   backup_id = "01a57d08-b7f5-584e-b64a-dc2236bb0438"
  url_expire_time = 12
  backup_download_restriction {
		restriction_type = "NONE"
		vpc_restriction_effect = "ALLOW"
		vpc_id_set = [local.vpc_id]
		ip_restriction_effect = "ALLOW"
		ip_set = ["0.0.0.0"]
  }
}

`
