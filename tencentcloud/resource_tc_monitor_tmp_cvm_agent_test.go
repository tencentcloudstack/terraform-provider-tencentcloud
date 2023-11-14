package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorTmpCvmAgentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpCvmAgent,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_cvm_agent.tmp_cvm_agent", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_cvm_agent.tmp_cvm_agent",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpCvmAgent = `

resource "tencentcloud_monitor_tmp_cvm_agent" "tmp_cvm_agent" {
  instance_id = "prom-dko9d0nu"
  name = "agent"
}

`
