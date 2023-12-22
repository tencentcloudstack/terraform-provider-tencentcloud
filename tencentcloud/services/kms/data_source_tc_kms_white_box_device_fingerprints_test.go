package kms_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKmsWhiteBoxDeviceFingerprintsDataSource_basic -v
func TestAccTencentCloudKmsWhiteBoxDeviceFingerprintsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsWhiteBoxDeviceFingerprintsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kms_white_box_device_fingerprints.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_device_fingerprints.example", "key_id"),
				),
			},
		},
	})
}

const testAccKmsWhiteBoxDeviceFingerprintsDataSource = `
data "tencentcloud_kms_white_box_device_fingerprints" "example" {
  key_id = "8731f440-66c1-11ee-beb0-52540036aed2"
}
`
