package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusTreeJobsDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusTreeJobsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusTreeJobsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_tree_jobs.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_tree_jobs.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusTreeJobsDataSource = `
data "tencentcloud_oceanus_tree_jobs" "example" {
  work_space_id = "space-2idq8wbr"
}
`
