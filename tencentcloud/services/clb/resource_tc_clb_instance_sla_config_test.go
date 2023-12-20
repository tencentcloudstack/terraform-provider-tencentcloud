package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClbInstanceSlaConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  load_balancer_id = "lb-5dnrkgry"
  sla_type         = "SLA"
}

`
