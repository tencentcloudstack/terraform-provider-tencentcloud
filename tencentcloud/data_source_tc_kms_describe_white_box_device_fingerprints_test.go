package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsDescribeWhiteBoxDeviceFingerprintsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDescribeWhiteBoxDeviceFingerprintsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_describe_white_box_device_fingerprints.describe_white_box_device_fingerprints")),
			},
		},
	})
}

const testAccKmsDescribeWhiteBoxDeviceFingerprintsDataSource = `

data "tencentcloud_kms_describe_white_box_device_fingerprints" "describe_white_box_device_fingerprints" {
  key_id = "244dab8c-6dad-11ea-80c6-5254006d0810"
}

`
