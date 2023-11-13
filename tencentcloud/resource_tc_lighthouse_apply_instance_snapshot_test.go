package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseApplyInstanceSnapshotResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseApplyInstanceSnapshot,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_apply_instance_snapshot.apply_instance_snapshot", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_apply_instance_snapshot.apply_instance_snapshot",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseApplyInstanceSnapshot = `

resource "tencentcloud_lighthouse_apply_instance_snapshot" "apply_instance_snapshot" {
  instance_id = "lhins-123456"
  snapshot_id = "lhsnap-123456"
}

`
