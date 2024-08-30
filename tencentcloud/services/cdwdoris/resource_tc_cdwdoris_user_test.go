package cdwdoris

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwdorisUserResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwdorisUser,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_user.cdwdoris_user", "id")),
		}, {
			ResourceName:      "tencentcloud_cdwdoris_user.cdwdoris_user",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccCdwdorisUser = `

resource "tencentcloud_cdwdoris_user" "cdwdoris_user" {
  user_info = {
  }
}
`
