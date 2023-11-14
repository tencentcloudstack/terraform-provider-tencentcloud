package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamAccessKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamAccessKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_access_key.access_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_access_key.access_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamAccessKey = `

resource "tencentcloud_cam_access_key" "access_key" {
  target_uin = &lt;nil&gt;
}

`
