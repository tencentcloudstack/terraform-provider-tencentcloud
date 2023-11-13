package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDescribeTopSpaceSchemasDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDescribeTopSpaceSchemasDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_describe_top_space_schemas.describe_top_space_schemas")),
			},
		},
	})
}

const testAccDbbrainDescribeTopSpaceSchemasDataSource = `

data "tencentcloud_dbbrain_describe_top_space_schemas" "describe_top_space_schemas" {
  instance_id = ""
  limit = 
  sort_by = ""
  product = ""
    }

`
