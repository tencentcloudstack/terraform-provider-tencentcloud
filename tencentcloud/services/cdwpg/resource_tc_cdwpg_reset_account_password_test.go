package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgResetAccountPasswordResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwpgAccount,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_reset_account_password.cdwpg_account", "id")),
		}, {
			ResourceName:            "tencentcloud_cdwpg_reset_account_password.cdwpg_account",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"new_password"},
		}},
	})
}

const testAccCdwpgAccount = `
resource "tencentcloud_cdwpg_reset_account_password" "cdwpg_account" {
	instance_id = "cdwpg-r3xgk6w3"
	user_name = "dbadmin"
	new_password = "testpassword"
}
`
