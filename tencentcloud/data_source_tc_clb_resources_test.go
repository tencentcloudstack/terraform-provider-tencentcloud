package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbResourcesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbResourcesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_resources.resources")),
			},
		},
	})
}

const testAccClbResourcesDataSource = `

data "tencentcloud_clb_resources" "resources" {
  filters {
    name = "isp"
    values = ["BGP"]
  }
}
`
