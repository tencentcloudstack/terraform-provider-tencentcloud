package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTseGroupListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseGroupListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_group_list.group_list")),
			},
		},
	})
}

const testAccTseGroupListDataSource = `

data "tencentcloud_tse_group_list" "group_list" {
  gateway_id = ""
  filters {
		name = "GroupId"
		values = 

  }
  }

`
