package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClbInstanceMixIpTargetConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceMixIpTargetConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_instance_mix_ip_target_config.instance_mix_ip_target_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_clb_instance_mix_ip_target_config.instance_mix_ip_target_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbInstanceMixIpTargetConfig = `

resource "tencentcloud_clb_instance_mix_ip_target_config" "instance_mix_ip_target_config" {
  load_balancer_id = "lb-5dnrkgry"
  mix_ip_target = false
}

`
