package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTcmTracingConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmTracingConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcm_tracing_config.tracing_config", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_tracing_config.tracingConfig",
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
  apm {
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
