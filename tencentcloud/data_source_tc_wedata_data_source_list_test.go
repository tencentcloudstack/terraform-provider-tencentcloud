package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDataSourceListDataSource_basic -v
func TestAccTencentCloudNeedFixWedataDataSourceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataSourceListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_data_source_list.example"),
				),
			},
		},
	})
}

const testAccWedataDataSourceListDataSource = `
data "tencentcloud_wedata_data_source_list" "example" {
  order_fields {
    name      = "create_time"
    direction = "DESC"
  }

  filters {
    name   = "Name"
    values = ["tf_example"]
  }
}
`
