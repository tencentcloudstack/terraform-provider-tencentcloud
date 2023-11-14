package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcmTracingConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmTracingConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcm_tracing_config.tracing_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcm_tracing_config.tracing_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcmTracingConfig = `

resource "tencentcloud_tcm_tracing_config" "tracing_config" {
  mesh_id = "mesh-xxxxxxxx"
  enable = true
  a_p_m {
		enable = true
		region = "ap-shanghai"
		instance_id = "apm-xxx"

  }
  sampling = 
  zipkin {
		address = "10.10.10.10:9411"

  }
}

`
