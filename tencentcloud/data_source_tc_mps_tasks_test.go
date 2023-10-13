package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsTasksDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_tasks.tasks")),
			},
		},
	})
}

const testAccMpsTasksDataSource = `

data "tencentcloud_mps_tasks" "tasks" {
  status = &lt;nil&gt;
  limit = &lt;nil&gt;
  scroll_token = &lt;nil&gt;
  }

`
