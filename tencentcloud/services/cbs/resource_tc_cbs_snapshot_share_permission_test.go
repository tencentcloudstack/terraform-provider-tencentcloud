package cbs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCbsSnapshotSharePermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshotSharePermission,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot_share_permission.snapshot_share_permission", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cbs_snapshot_share_permission.snapshot_share_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCbsSnapshotSharePermission = `
resource "tencentcloud_cbs_snapshot_share_permission" "snapshot_share_permission" {
	account_ids = ["100022975249"]
	snapshot_id = "snap-6qtrq4fn"
}
`
