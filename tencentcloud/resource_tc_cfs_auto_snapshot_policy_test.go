package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCfsAutoSnapshotPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsAutoSnapshotPolicy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfs_auto_snapshot_policy.auto_snapshot_policy", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfs_auto_snapshot_policy.auto_snapshot_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfsAutoSnapshotPolicy = `

resource "tencentcloud_cfs_auto_snapshot_policy" "auto_snapshot_policy" {
  hour = "&quot;2,3&quot;"
  policy_name = "policy_name"
  day_of_week = "&quot;1,2&quot;"
  alive_days = 7
  day_of_month = "2"
  interval_days = 1
}

`
