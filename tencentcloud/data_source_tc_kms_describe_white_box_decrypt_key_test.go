package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsDescribeWhiteBoxDecryptKeyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDescribeWhiteBoxDecryptKeyDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_describe_white_box_decrypt_key.describe_white_box_decrypt_key")),
			},
		},
	})
}

const testAccKmsDescribeWhiteBoxDecryptKeyDataSource = `

data "tencentcloud_kms_describe_white_box_decrypt_key" "describe_white_box_decrypt_key" {
  key_id = "244dab8c-6dad-11ea-80c6-5254006d0810"
  }

`
