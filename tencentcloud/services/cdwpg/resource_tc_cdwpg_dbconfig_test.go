package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgDbconfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdwpgDbconfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_dbconfig.cdwpg_dbconfig", "id")),
			},
		},
	})
}

const testAccCdwpgDbconfig = `
resource "tencentcloud_cdwpg_dbconfig" "cdwpg_dbconfig" {
  instance_id = "cdwpg-ua8wkqrt"
  node_config_params {
	node_type = "cn"
	parameter_name = "log_min_duration_statement"
	parameter_value = "10001"
  }
}
`
