package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmImageSharePermissionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageSharePermissionDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_image_share_permission.image_share_permission")),
			},
		},
	})
}

const testAccCvmImageSharePermissionDataSource = `

data "tencentcloud_cvm_image_share_permission" "image_share_permission" {
  image_id = "img-k4h0m5la"
}
`
