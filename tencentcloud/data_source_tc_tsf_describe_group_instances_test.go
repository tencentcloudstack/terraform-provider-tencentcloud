package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeGroupInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeGroupInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_group_instances.describe_group_instances")),
			},
		},
	})
}

const testAccTsfDescribeGroupInstancesDataSource = `

data "tencentcloud_tsf_describe_group_instances" "describe_group_instances" {
  group_id = ""
  search_word = ""
  order_by = ""
  order_type = 
  }

`
