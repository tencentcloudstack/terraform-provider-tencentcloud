package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDescribeUserTypeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDescribeUserTypeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_describe_user_type.describe_user_type")),
			},
		},
	})
}

const testAccDlcDescribeUserTypeDataSource = `

data "tencentcloud_dlc_describe_user_type" "describe_user_type" {
  user_id = "127382378"
  }

`
