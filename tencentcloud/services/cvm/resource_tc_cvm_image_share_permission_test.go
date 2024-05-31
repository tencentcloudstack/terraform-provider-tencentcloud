package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmModifyImageSharePermissionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmModifyImageSharePermission,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_image_share_permission.image_share_permission", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_image_share_permission.image_share_permission", "image_id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_image_share_permission.image_share_permission", "account_ids.#", "1"),
				),
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
  image_id    = "img-l7uxaine"
  account_ids = ["100022975249"]
}
`
