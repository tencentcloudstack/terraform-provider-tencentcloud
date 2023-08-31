package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseBackupJobDetailDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseBackupJobDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_backup_job_detail.backup_job_detail")),
			},
		},
	})
}

const testAccClickhouseBackupJobDetailDataSource = `
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = "cdwch-pcap78rz"
	cos_bucket_name = "keep-export-image-1308726196"
}

data "tencentcloud_clickhouse_backup_job_detail" "backup_job_detail" {
	instance_id = "cdwch-pcap78rz"
	back_up_job_id = 7679
}
`
