package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmImageSharePermissionResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageSharePermissionResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_image_share_permission.image_share_permission", "id"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_image_share_permission.image_share_permission", "image_id"), resource.TestCheckResourceAttr("tencentcloud_cvm_image_share_permission.image_share_permission", "account_ids.#", "1")),
			},
			{
				ResourceName:      "tencentcloud_cvm_image_share_permission.image_share_permission",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmImageSharePermissionResource_BasicCreate = `

resource "tencentcloud_cvm_image_share_permission" "image_share_permission" {
    image_id = "img-l7uxaine"
    account_ids = [100022975249]
}

`
