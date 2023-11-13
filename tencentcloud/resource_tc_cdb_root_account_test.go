package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRootAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRootAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_root_account.root_account", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_root_account.root_account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRootAccount = `

resource "tencentcloud_cdb_root_account" "root_account" {
  instance_id = ""
  password = ""
}

`
