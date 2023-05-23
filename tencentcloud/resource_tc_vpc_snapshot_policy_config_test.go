package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcSnapshotPolicyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSnapshotPolicyConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_snapshot_policy_config.snapshot_policy_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_snapshot_policy_config.snapshot_policy_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcSnapshotPolicyConfig = `

resource "tencentcloud_vpc_snapshot_policy_config" "snapshot_policy_config" {
  snapshot_policy_id = "sspolicy-1t6cobbv"
  enable             = false
}

`
