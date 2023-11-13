package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbTargetGroupListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_target_group_list.target_group_list")),
			},
		},
	})
}

const testAccClbTargetGroupListDataSource = `

data "tencentcloud_clb_target_group_list" "target_group_list" {
  target_group_ids = 
  filters {
		name = ""
		values = 

  }
  }

`
