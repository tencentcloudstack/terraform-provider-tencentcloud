package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmImageSharePermissionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
  image_id = "img-6pb6lrmy"
  }

`
