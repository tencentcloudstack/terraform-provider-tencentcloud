package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwSyncAssetResource_basic -v
func TestAccTencentCloudNeedFixCfwSyncAssetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwCfwSyncAsset,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_policy.example", "id"),
				),
			},
		},
	})
}

const testAccCfwCfwSyncAsset = `
resource "tencentcloud_cfw_sync_asset" "example" {}
`
