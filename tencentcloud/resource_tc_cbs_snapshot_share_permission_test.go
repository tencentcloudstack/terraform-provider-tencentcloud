package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCbsSnapshotSharePermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCbsSnapshotSharePermission,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cbs_snapshot_share_permission.snapshot_share_permission", "id")),
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
  account_ids = 
  permission = "SHARE"
  snapshot_ids = 
}

`
