package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseBackupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_backup.backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_clickhouse_backup.backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClickhouseBackup = tcacctest.DefaultClickhouseVariables + `
resource "tencentcloud_clickhouse_backup" "backup" {
	instance_id = var.instance_id
	cos_bucket_name = "keep-export-image-1308726196"
}
`
