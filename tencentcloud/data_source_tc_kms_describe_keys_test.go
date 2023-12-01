package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKmsDescribeKeysDataSource_basic -v
func TestAccTencentCloudKmsDescribeKeysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyListsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_describe_keys.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.key_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.key_state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.key_usage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.creator_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.key_rotation_enabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.owner"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.next_rotate_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.origin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kms_describe_keys.example", "key_list.0.valid_to"),
				),
			},
		},
	})
}

const testAccKmsKeyListsDataSource = `
data "tencentcloud_kms_describe_keys" "example" {
  key_ids = [
    "72688f39-1fe8-11ee-9f1a-525400cf25a4"
  ]
}
`
