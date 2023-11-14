package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDescribeTopSpaceTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDescribeTopSpaceTablesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_describe_top_space_tables.describe_top_space_tables")),
			},
		},
	})
}

const testAccDbbrainDescribeTopSpaceTablesDataSource = `

data "tencentcloud_dbbrain_describe_top_space_tables" "describe_top_space_tables" {
  instance_id = ""
  limit = 
  sort_by = ""
  product = ""
    }

`
