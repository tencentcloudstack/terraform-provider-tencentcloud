package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmModifyImageSharePermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmModifyImageSharePermission,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_modify_image_share_permission.modify_image_share_permission", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_modify_image_share_permission.modify_image_share_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmModifyImageSharePermission = `

resource "tencentcloud_cvm_modify_image_share_permission" "modify_image_share_permission" {
  image_id = "img-gvbnzy6f"
  account_ids = 
  permission = "SHARE"
}

`
