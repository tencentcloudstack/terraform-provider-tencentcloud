package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDataSourceWithoutInfoDataSource_basic -v
func TestAccTencentCloudNeedFixWedataDataSourceWithoutInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataSourceWithoutInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_data_source_without_info.example"),
				),
			},
		},
	})
}

const testAccWedataDataSourceWithoutInfoDataSource = `
data "tencentcloud_wedata_data_source_without_info" "example" {
  filters {
    name   = "ownerProjectId"
    values = ["1612982498218618880"]
  }

  order_fields {
    name      = "create_time"
    direction = "DESC"
  }
}
`
