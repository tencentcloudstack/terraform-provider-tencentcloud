package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmModifyImageSharePermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmModifyImageSharePermission,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_image_share_permission.image_share_permission", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_image_share_permission.image_share_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmModifyImageSharePermission = `
resource "tencentcloud_cvm_image_share_permission" "image_share_permission" {
	image_id = "img-b0x811s0"
	account_ids = ["100022975249"]
}
`
