package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcmAccessLogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmAccessLogConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcm_access_log_config.access_log_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcm_access_log_config.access_log_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcmAccessLogConfig = `

resource "tencentcloud_tcm_access_log_config" "access_log_config" {
  mesh_name = "mesh-xxxxxxxx"
  selected_range {
		items {
			namespace = "prod"
			gateways = 
		}
		all = false

  }
  template = "istio"
  enable = true
  c_l_s {
		enable = true
		log_set = "mesh-xxx"
		topic = "accesslog"

  }
  encoding = "TEXT"
  format = "[%START_TIME%]"
  enable_stdout = false
  enable_server = false
  address = "xxx"
}

`
