package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTatInvocationTaskDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvocationTaskDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tat_invocation_task.invocation_task")),
			},
		},
	})
}

const testAccTatInvocationTaskDataSource = `

data "tencentcloud_tat_invocation_task" "invocation_task" {
  invocation_task_ids = 
  filters {
		name = ""
		values = 

  }
  hide_output = 
  }

`
