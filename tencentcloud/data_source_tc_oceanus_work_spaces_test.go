package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusWorkSpacesDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusWorkSpacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusWorkSpacesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_work_spaces.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_work_spaces.example", "order_type"),
				),
			},
		},
	})
}

const testAccOceanusWorkSpacesDataSource = `
data "tencentcloud_oceanus_work_spaces" "example" {
  order_type = 1
  filters {
    name   = "WorkSpaceName"
    values = ["keep-work-space"]
  }
}
`
