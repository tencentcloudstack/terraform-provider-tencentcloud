package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTemWorkload_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemWorkload,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tem_workload.workload", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_workload.workload",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemWorkload = `

resource "tencentcloud_tem_workload" "workload" {
  application_id     = "app-3j29aa2p"
  environment_id     = "en-853mggjm"
  deploy_version     = "hello-world"
  deploy_mode        = "IMAGE"
  img_repo           = "tem_demo/tem_demo"
  init_pod_num       = 1
  cpu_spec           = 1
  memory_spec        = 1
  liveness {
    type                  = "HttpGet"
    protocol              = "HTTP"
    path                  = "/"
    port                  = 8080
    initial_delay_seconds = 0
    timeout_seconds       = 1
    period_seconds        = 10

  }
  readiness {
    type                  = "HttpGet"
    protocol              = "HTTP"
    path                  = "/"
    port                  = 8000
    initial_delay_seconds = 0
    timeout_seconds       = 1
    period_seconds        = 10

  }
  startup_probe {
    type                  = "HttpGet"
    protocol              = "HTTP"
    path                  = "/"
    port                  = 36000
    initial_delay_seconds = 0
    timeout_seconds       = 1
    period_seconds        = 10

  }
}

`
