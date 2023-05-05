package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCfsSignUpCfsServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsSignUpCfsService,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfs_sign_up_cfs_service.sign_up_cfs_service", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfs_sign_up_cfs_service.sign_up_cfs_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfsSignUpCfsService = `

resource "tencentcloud_cfs_sign_up_cfs_service" "sign_up_cfs_service" {
  }

`
