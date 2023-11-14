package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsGetPublicKeyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsGetPublicKeyDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_get_public_key.get_public_key")),
			},
		},
	})
}

const testAccKmsGetPublicKeyDataSource = `

data "tencentcloud_kms_get_public_key" "get_public_key" {
      }

`
