package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusTreeJobsDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusTreeJobsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusTreeJobsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_tree_jobs.example"),
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
