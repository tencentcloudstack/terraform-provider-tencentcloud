package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDeployVMGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDeployVMGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_deploy_v_m_group.deploy_v_m_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_deploy_v_m_group.deploy_v_m_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfDeployVMGroup = `

resource "tencentcloud_tsf_deploy_v_m_group" "deploy_v_m_group" {
  group_id = ""
  pkg_id = ""
  startup_parameters = ""
  deploy_desc = ""
  force_start = 
  enable_health_check = 
  health_check_settings {
		liveness_probe {
			action_type = ""
			initial_delay_seconds = 
			timeout_seconds = 
			period_seconds = 
			success_threshold = 
			failure_threshold = 
			scheme = ""
			port = 
			path = ""
			command = 
			type = ""
		}
		readiness_probe {
			action_type = ""
			initial_delay_seconds = 
			timeout_seconds = 
			period_seconds = 
			success_threshold = 
			failure_threshold = 
			scheme = ""
			port = 
			path = ""
			command = 
			type = ""
		}

  }
  update_type = 
  deploy_beta_enable = 
  deploy_batch = 
  deploy_exe_mode = ""
  deploy_wait_time = 
  start_script = ""
  stop_script = ""
  incremental_deployment = 
  jdk_name = ""
  jdk_version = ""
  agent_profile_list {
		agent_type = ""
		agent_version = ""

  }
  warmup_setting {
		enabled = 
		warmup_time = 
		curvature = 
		enabled_protection = 

  }
}

`
