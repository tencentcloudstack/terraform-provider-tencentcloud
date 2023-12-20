package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseBackupJobDetailDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseBackupJobDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_clickhouse_backup_job_detail.backup_job_detail")),
			},
		},
	})
}

const testAccClickhouseBackupJobDetailDataSource = tcacctest.DefaultClickhouseVariables + `
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = var.instance_id
	cos_bucket_name = "keep-export-image-1308726196"
}

data "tencentcloud_clickhouse_backup_job_detail" "backup_job_detail" {
	instance_id = var.instance_id
	back_up_job_id = 7679
}
`
