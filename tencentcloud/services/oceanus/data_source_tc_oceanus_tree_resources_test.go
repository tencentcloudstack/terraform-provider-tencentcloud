package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusTreeResourcesDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusTreeResourcesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusTreeResourcesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_tree_resources.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_tree_resources.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusTreeResourcesDataSource = `
data "tencentcloud_oceanus_tree_resources" "example" {
  work_space_id = "space-2idq8wbr"
}
`
