package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseDisksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseDisksDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_disks.disks")),
			},
		},
	})
}

const testAccLighthouseDisksDataSource = tcacctest.DefaultLighthoustVariables + `

data "tencentcloud_lighthouse_disks" "disks" {
  disk_ids = [var.lighthouse_backup_disk_id]
}
`
