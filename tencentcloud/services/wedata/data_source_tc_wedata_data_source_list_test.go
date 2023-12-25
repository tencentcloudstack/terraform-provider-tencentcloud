package wedata_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataDataSourceListDataSource_basic -v
func TestAccTencentCloudNeedFixWedataDataSourceListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataDataSourceListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_wedata_data_source_list.example"),
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
