package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgRestartInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCdwpgRestartInstance,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_restart_instance.cdwpg_restart_instance", "id")),
		}},
	})
}

const testAccCdwpgRestartInstance = `
resource "tencentcloud_cdwpg_restart_instance" "cdwpg_restart_instance" {
	instance_id = "cdwpg-zpiemnyd"
	node_types = ["gtm"]
}
`
