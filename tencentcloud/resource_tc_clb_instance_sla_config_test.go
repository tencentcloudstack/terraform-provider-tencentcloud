package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbInstanceSlaConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceSlaConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_instance_sla_config.instance_sla_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_clb_instance_sla_config.instance_sla_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbInstanceSlaConfig = `

resource "tencentcloud_clb_instance_sla_config" "instance_sla_config" {
  load_balancer_sla {
		load_balancer_id = "lb-xxxxxxxx"
		sla_type = ""

  }
}

`
