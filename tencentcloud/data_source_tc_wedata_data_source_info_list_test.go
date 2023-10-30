package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDataSourceInfoListDataSource_basic -v
func TestAccTencentCloudNeedFixWedataDataSourceInfoListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataSourceInfoListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_data_source_info_list.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_wedata_data_source_info_list.example", "project_id"),
				),
			},
		},
	})
}

const testAccWedataDataSourceInfoListDataSource = `
data "tencentcloud_wedata_data_source_info_list" "example" {
  project_id = "1927766435649077248"
  filters {
    name   = "Name"
    values = ["tf_example"]
  }

  order_fields {
    name      = "CreateTime"
    direction = "DESC"
  }
}
`
