package cfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCfsUserQuotaResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  file_system_id = "cfs-4636029bc"
  user_type = "Uid"
  user_id = "2159973417"
  capacity_hard_limit = 10
  file_hard_limit = 10000
}

`
