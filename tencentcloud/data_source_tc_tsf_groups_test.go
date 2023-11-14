package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGroupsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_groups.groups")),
			},
		},
	})
}

const testAccTsfGroupsDataSource = `

data "tencentcloud_tsf_groups" "groups" {
  search_word = ""
  application_id = ""
  order_by = ""
  order_type = 
  namespace_id = ""
  cluster_id = ""
  group_resource_type_list = 
  status = ""
  group_id_list = 
  }

`
