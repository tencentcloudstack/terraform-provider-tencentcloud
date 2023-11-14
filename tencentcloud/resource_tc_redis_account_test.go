package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_account.account", "id")),
			},
			{
				ResourceName:      "tencentcloud_redis_account.account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRedisAccount = `

resource "tencentcloud_redis_account" "account" {
  instance_id = "crs-c1nl9rpv"
  account_name = "user"
  account_password = &lt;nil&gt;
  remark = &lt;nil&gt;
  readonly_policy = 
  privilege = "rw"
}

`
