package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusWorkSpacesDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusWorkSpacesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusWorkSpacesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_work_spaces.example"),
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
