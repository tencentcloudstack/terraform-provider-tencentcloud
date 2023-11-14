package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTemLogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTemLogConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tem_log_config.log_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_tem_log_config.log_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTemLogConfig = `

resource "tencentcloud_tem_log_config" "log_config" {
  environment_id = "en-xxx"
  application_id = "en-xxx"
  name = "xxx"
  logset_id = "xxx"
  topic_id = "xxx"
  input_type = "container_stdout"
  log_type = "minimalist_log"
  beginning_regex = "**.log"
  log_path = "/xxx"
  file_pattern = "*.log"
}

`
