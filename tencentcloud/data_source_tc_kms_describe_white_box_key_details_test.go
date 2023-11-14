package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsDescribeWhiteBoxKeyDetailsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDescribeWhiteBoxKeyDetailsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_describe_white_box_key_details.describe_white_box_key_details")),
			},
		},
	})
}

const testAccKmsDescribeWhiteBoxKeyDetailsDataSource = `

data "tencentcloud_kms_describe_white_box_key_details" "describe_white_box_key_details" {
  key_status = 
  }

`
