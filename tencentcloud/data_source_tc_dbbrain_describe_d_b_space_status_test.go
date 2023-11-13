package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDescribeDBSpaceStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDescribeDBSpaceStatusDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_describe_d_b_space_status.describe_d_b_space_status")),
			},
		},
	})
}

const testAccDbbrainDescribeDBSpaceStatusDataSource = `

data "tencentcloud_dbbrain_describe_d_b_space_status" "describe_d_b_space_status" {
  instance_id = ""
  range_days = 
  product = ""
        }

`
