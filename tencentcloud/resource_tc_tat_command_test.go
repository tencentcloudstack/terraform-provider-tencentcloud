package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTatCommandResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatCommand,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tat_command.command", "id")),
			},
			{
				ResourceName:      "tencentcloud_tat_command.command",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTatCommand = `

resource "tencentcloud_tat_command" "command" {
  command_name = &lt;nil&gt;
  content = &lt;nil&gt;
  description = &lt;nil&gt;
  command_type = &lt;nil&gt;
  working_directory = &lt;nil&gt;
  timeout = &lt;nil&gt;
  enable_parameter = &lt;nil&gt;
  default_parameters = &lt;nil&gt;
  tags {
		key = &lt;nil&gt;
		value = &lt;nil&gt;

  }
  username = &lt;nil&gt;
  output_c_o_s_bucket_url = &lt;nil&gt;
  output_c_o_s_key_prefix = &lt;nil&gt;
        }

`
