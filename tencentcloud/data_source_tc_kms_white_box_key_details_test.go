package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixKmsWhiteBoxKeyDetailsDataSource_basic -v
func TestAccTencentCloudNeedFixKmsWhiteBoxKeyDetailsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsWhiteBoxKeyDetailsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_white_box_key_details.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.algorithm"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.decrypt_key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.resource_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.key_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.creator_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.alias"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.encrypt_key"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.owner_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_white_box_key_details.example", "key_infos.0.device_fingerprint_bind"),
				),
			},
		},
	})
}

const testAccKmsWhiteBoxKeyDetailsDataSource = `
data "tencentcloud_kms_white_box_key_details" "example" {
  key_status = 0
}
`
