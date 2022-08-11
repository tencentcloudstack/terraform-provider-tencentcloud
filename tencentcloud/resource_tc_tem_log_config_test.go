package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTemLogConfig_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemLogConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tem_log_config.logConfig", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_log_config.logConfig",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemLogConfig = `

resource "tencentcloud_tem_log_config" "logConfig" {
  environment_id = "en-853mggjm"
  application_id = "app-3j29aa2p"
  name = "terraform"
  logset_id = "b5824781-8d5b-4029-a2f7-d03c37f72bdf"
  topic_id = "a21a488d-d28f-4ac3-8044-bdf8c91b49f2"
  input_type = "container_stdout"
  log_type = "minimalist_log"
}

`
