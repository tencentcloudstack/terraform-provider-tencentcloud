package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgUserhbaResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwpgUserhba,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_userhba.cdwpg_userhba", "id")),
		}, {
			ResourceName:      "tencentcloud_cdwpg_userhba.cdwpg_userhba",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccCdwpgUserhba = `

resource "tencentcloud_cdwpg_userhba" "cdwpg_userhba" {
  instance_id = "cdwpg-zpiemnyd"
  hba_configs {
	type = "host"
	database = "all"
	user = "all"
	address = "0.0.0.0/0"
	method = "md5"
  }
}
`
