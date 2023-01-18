package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCfsUserQuotaResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfsUserQuota,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cfs_user_quota.user_quota", "id")),
			},
			{
				ResourceName:      "tencentcloud_cfs_user_quota.user_quota",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCfsUserQuota = `

resource "tencentcloud_cfs_user_quota" "user_quota" {
  file_system_id = ""
  user_type = ""
  user_id = ""
  capacity_hard_limit = 
  file_hard_limit = 
}

`
