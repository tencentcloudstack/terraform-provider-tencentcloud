package cdwdoris

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwdorisInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwdorisInstance,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.cdwdoris_instance", "id")),
		}, {
			ResourceName:      "tencentcloud_cdwdoris_instance.cdwdoris_instance",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccCdwdorisInstance = `

resource "tencentcloud_cdwdoris_instance" "cdwdoris_instance" {
  fe_spec = {
  }
  be_spec = {
  }
  charge_properties = {
  }
  tags = {
  }
  user_multi_zone_infos = {
  }
}
`
