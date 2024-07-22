package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmImageSharePermissionDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageSharePermissionDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_image_share_permission.image_share_permission"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_share_permission.image_share_permission", "image_id", "img-l7uxaine"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_share_permission.image_share_permission", "share_permission_set.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_image_share_permission.image_share_permission", "share_permission_set.0.created_time"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_share_permission.image_share_permission", "share_permission_set.0.account_id", "100022975249")),
			},
		},
	})
}

const testAccCvmImageSharePermissionDataSource_BasicCreate = `

data "tencentcloud_cvm_image_share_permission" "image_share_permission" {
    image_id = tencentcloud_cvm_image_share_permission.image_share_permission.image_id
}
resource "tencentcloud_cvm_image_share_permission" "image_share_permission" {
    image_id = "img-l7uxaine"
    account_ids = [100022975249]
}

`
