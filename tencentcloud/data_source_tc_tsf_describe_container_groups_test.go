package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeContainerGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeContainerGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_container_groups.describe_container_groups")),
			},
		},
	})
}

const testAccTsfDescribeContainerGroupsDataSource = `

data "tencentcloud_tsf_describe_container_groups" "describe_container_groups" {
  search_word = ""
  application_id = ""
  order_by = ""
  order_type = 
  cluster_id = ""
  namespace_id = ""
  }

`
