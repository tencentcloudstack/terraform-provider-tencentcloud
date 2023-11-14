package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamUserAccountGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUserAccountGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_user_account_group.user_account_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_user_account_group.user_account_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamUserAccountGroup = `

resource "tencentcloud_cam_user_account_group" "user_account_group" {
  info {
		group_id = &lt;nil&gt;
		uid = &lt;nil&gt;
		uin = &lt;nil&gt;

  }
}

`
