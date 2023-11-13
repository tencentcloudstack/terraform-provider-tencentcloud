package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTatInvocationCommandAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvocationCommandAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tat_invocation_command_attachment.invocation_command_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTatInvocationCommandAttachment = `

resource "tencentcloud_tat_invocation_command_attachment" "invocation_command_attachment" {
  content = ""
  instance_ids = 
  command_name = ""
  description = ""
  command_type = ""
  working_directory = ""
  timeout = 
  save_command = 
  enable_parameter = 
  default_parameters = ""
  parameters = ""
  username = ""
  output_c_o_s_bucket_url = ""
  output_c_o_s_key_prefix = ""
}

`
