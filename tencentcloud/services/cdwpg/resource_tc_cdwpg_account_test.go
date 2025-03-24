package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwpgAccount,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_account.cdwpg_account", "id")),
		}, {
			ResourceName:            "tencentcloud_cdwpg_account.cdwpg_account",
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"new_password"},
		}},
	})
}

const testAccCdwpgAccount = `
resource "tencentcloud_cdwpg_account" "cdwpg_account" {
	instance_id = "cdwpg-zpiemnyd"
	user_name = "dbadmin"
	new_password = "testpassword"
}
`
