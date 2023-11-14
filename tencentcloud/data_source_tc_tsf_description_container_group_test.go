package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescriptionContainerGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescriptionContainerGroupDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_description_container_group.description_container_group")),
			},
		},
	})
}

const testAccTsfDescriptionContainerGroupDataSource = `

data "tencentcloud_tsf_description_container_group" "description_container_group" {
  search_word = ""
  application_id = ""
  order_by = ""
  order_type = 
  cluster_id = ""
  namespace_id = ""
  }

`
