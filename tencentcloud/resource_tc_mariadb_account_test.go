package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_account.account", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_account.account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbAccount = `

resource "tencentcloud_mariadb_account" "account" {
  instance_id = &lt;nil&gt;
  user_name = &lt;nil&gt;
  host = &lt;nil&gt;
  password = &lt;nil&gt;
  read_only = &lt;nil&gt;
  description = &lt;nil&gt;
}

`
