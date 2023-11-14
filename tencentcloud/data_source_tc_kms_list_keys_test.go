package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_list_keys.list_keys")),
			},
		},
	})
}

const testAccKmsListKeysDataSource = `

data "tencentcloud_kms_list_keys" "list_keys" {
  offset = 0
  limit = 2
  role = 
  hsm_cluster_id = "0"
}

`
