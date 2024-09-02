package cdwdoris

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwdorisWorkloadGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwdorisWorkloadGroup,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_workload_group.cdwdoris_workload_group", "id")),
		}, {
			ResourceName:      "tencentcloud_cdwdoris_workload_group.cdwdoris_workload_group",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccCdwdorisWorkloadGroup = `

resource "tencentcloud_cdwdoris_workload_group" "cdwdoris_workload_group" {
  instance_id = ""
  workload_group = ""
}
`
