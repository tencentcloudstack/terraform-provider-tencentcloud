package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseApplyInstanceSnapshotResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseApplyInstanceSnapshot,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_apply_instance_snapshot.apply_instance_snapshot", "id")),
			},
		},
	})
}

const testAccLighthouseApplyInstanceSnapshot = tcacctest.DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_apply_instance_snapshot" "apply_instance_snapshot" {
	instance_id = var.lighthouse_id
	snapshot_id = var.lighthouse_snapshot_id
}
`
