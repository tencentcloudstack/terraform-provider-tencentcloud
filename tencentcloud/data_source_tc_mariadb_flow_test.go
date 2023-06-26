package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMariadbFlowDataSource_basic -v
func TestAccTencentCloudNeedFixMariadbFlowDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbFlowDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_flow.flow")),
			},
		},
	})
}

const testAccMariadbFlowDataSource = `
data "tencentcloud_mariadb_flow" "flow" {
  flow_id = 1307
}
`
