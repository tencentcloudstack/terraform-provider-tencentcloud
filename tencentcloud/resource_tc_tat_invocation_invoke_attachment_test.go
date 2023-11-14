package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTatInvocationInvokeAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvocationInvokeAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTatInvocationInvokeAttachment = `

resource "tencentcloud_tat_invocation_invoke_attachment" "invocation_invoke_attachment" {
  instance_ids = 
  working_directory = ""
  timeout = 
  parameters = ""
  username = ""
  output_c_o_s_bucket_url = ""
  output_c_o_s_key_prefix = ""
  command_id = ""
}

`
