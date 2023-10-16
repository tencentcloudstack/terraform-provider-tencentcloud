package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsResource_basic -v
func TestAccTencentCloudKmsOverwriteWhiteBoxDeviceFingerprintsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsOverwriteWhiteBoxDeviceFingerprints,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kms_overwrite_white_box_device_fingerprints.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kms_overwrite_white_box_device_fingerprints.example", "key_id"),
				),
			},
		},
	})
}

const testAccKmsOverwriteWhiteBoxDeviceFingerprints = `
resource "tencentcloud_kms_overwrite_white_box_device_fingerprints" "example" {
  key_id = "8731f440-66c1-11ee-beb0-52540036aed2"
}
`
