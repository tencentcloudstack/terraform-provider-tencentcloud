package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorTmpTkeBasicConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpTkeBasicConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_tke_basic_config.tmp_tke_basic_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_tke_basic_config.tmp_tke_basic_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpTkeBasicConfig = `

resource "tencentcloud_monitor_tmp_tke_basic_config" "tmp_tke_basic_config" {
  instance_id = ""
  cluster_type = ""
  cluster_id = ""
  service_monitors {
		name = ""
		config = ""
		template_id = ""

  }
  pod_monitors {
		name = ""
		config = ""
		template_id = ""

  }
  raw_jobs {
		name = ""
		config = ""
		template_id = ""

  }
}

`
