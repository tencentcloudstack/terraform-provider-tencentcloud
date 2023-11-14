package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApiGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiGroupDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_api_group.api_group")),
			},
		},
	})
}

const testAccTsfApiGroupDataSource = `

data "tencentcloud_tsf_api_group" "api_group" {
  search_word = ""
  group_type = ""
  auth_type = ""
  status = ""
  order_by = ""
  order_type = 
  gateway_instance_id = "gw-ins-lvdypq5k"
  }

`
