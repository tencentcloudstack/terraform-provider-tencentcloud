package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseSnapshotResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseSnapshot,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_snapshot.snapshot", "snapshot_name", "snapshot_test"),
				),
			},
			{
				Config: testAccLighthouseSnapshot_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_snapshot.snapshot", "id"),
					resource.TestCheckResourceAttr("tencentcloud_lighthouse_snapshot.snapshot", "snapshot_name", "snapshot_test_update"),
				),
			},
		},
	})
}

const testAccLighthouseSnapshot = tcacctest.DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_snapshot" "snapshot" {
	instance_id = var.lighthouse_id
	snapshot_name = "snapshot_test"
}
`

const testAccLighthouseSnapshot_update = tcacctest.DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_snapshot" "snapshot" {
	instance_id = var.lighthouse_id
	snapshot_name = "snapshot_test_update"
}
`
