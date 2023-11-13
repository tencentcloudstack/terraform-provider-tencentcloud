package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.account", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_account.account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbAccount = `

resource "tencentcloud_dcdb_account" "account" {
  instance_id = &lt;nil&gt;
  user_name = &lt;nil&gt;
  host = &lt;nil&gt;
  password = &lt;nil&gt;
  read_only = &lt;nil&gt;
  description = &lt;nil&gt;
  max_user_connections = &lt;nil&gt;
}

`
