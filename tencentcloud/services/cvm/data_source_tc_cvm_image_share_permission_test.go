package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmImageSharePermissionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageSharePermissionDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_image_share_permission.image_share_permission"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_share_permission.image_share_permission", "image_id", "img-l7uxaine"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_share_permission.image_share_permission", "share_permission_set.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_image_share_permission.image_share_permission", "share_permission_set.0.created_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_share_permission.image_share_permission", "share_permission_set.0.account_id", "100022975249"),
				),
			},
		},
	})
}

const testAccCvmImageSharePermissionDataSource = testAccCvmModifyImageSharePermission + `

data "tencentcloud_cvm_image_share_permission" "image_share_permission" {
  image_id = tencentcloud_cvm_image_share_permission.image_share_permission.image_id
}
`
