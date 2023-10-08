package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKmsListKeysDataSource_basic -v
func TestAccTencentCloudKmsListKeysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsListKeysDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_list_keys.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_list_keys.example", "keys.0.key_id"),
				),
			},
		},
	})
}

const testAccKmsListKeysDataSource = `
data "tencentcloud_kms_list_keys" "example" {
  role = 1
}
`
