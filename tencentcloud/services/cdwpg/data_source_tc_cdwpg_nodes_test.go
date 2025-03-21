package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgNodesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwpgNodesDataSource,
			Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cdwpg_nodes.cdwpg_nodes")),
		}},
	})
}

const testAccCdwpgNodesDataSource = `
data "tencentcloud_cdwpg_nodes" "cdwpg_nodes" {
	instance_id = "cdwpg-gexy9tue"
}
`
