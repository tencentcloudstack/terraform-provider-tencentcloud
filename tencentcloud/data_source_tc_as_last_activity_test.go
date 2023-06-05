package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAsLastActivityDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsLastActivityDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_as_last_activity.last_activity")),
			},
		},
	})
}

const testAccAsLastActivityDataSource = `

data "tencentcloud_as_last_activity" "last_activity" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}

`
