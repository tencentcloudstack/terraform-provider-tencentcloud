
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsCompareTasksDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDtsCompareTasks,
				Check: resource.ComposeTestCheckFunc(
				  testAccCheckTencentCloudDataSourceID("data.tencentcloud_dts_compare_tasks.compare_tasks"),
				),
			},
		},
	})
}

const testAccDataSourceDtsCompareTasks = `

data "tencentcloud_dts_compare_tasks" "compare_tasks" {
  job_id = ""
  }

`
