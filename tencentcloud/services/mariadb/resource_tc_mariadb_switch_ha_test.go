package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMariadbSwitchHAResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbSwitchHA,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_switch_ha.switch_ha", "id"),
				),
			},
		},
	})
}

const testAccMariadbSwitchHA = `
resource "tencentcloud_mariadb_switch_ha" "switch_ha" {
  instance_id = "tdsql-9vqvls95"
  zone        = "ap-guangzhou-2"
}
`
