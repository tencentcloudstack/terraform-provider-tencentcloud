package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseBackupJobsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseBackupJobsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_backup_jobs.backup_jobs")),
			},
		},
	})
}

const testAccClickhouseBackupJobsDataSource = `
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = "cdwch-pcap78rz"
	cos_bucket_name = "keep-export-image-1308726196"
}

data "tencentcloud_clickhouse_backup_jobs" "backup_jobs" {
	instance_id = "cdwch-pcap78rz"
}
`
