package tat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvocationCommandAttachmentResource_basic -v
func TestAccTencentCloudTatInvocationCommandAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvocationCommandAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "instance_id", tcacctest.DefaultInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "working_directory", "/root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "timeout", "100"),
					// resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "description", "shell test"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "command_type", "SHELL"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "username", "root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "output_cos_bucket_url", "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_command_attachment.invocation_command_attachment", "output_cos_key_prefix", "log"),
				),
			},
		},
	})
}

const testAccTatInvocationCommandAttachmentVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultInstanceId + `"
}
`

const testAccTatInvocationCommandAttachment = testAccTatInvocationCommandAttachmentVar + `

resource "tencentcloud_tat_invocation_command_attachment" "invocation_command_attachment" {
	content = base64encode("pwd")
	instance_id = var.instance_id
	command_name = "terraform-test"
	description = "shell test"
	command_type = "SHELL"
	working_directory = "/root"
	timeout = 100
	save_command = false
	enable_parameter = false
	# default_parameters = "{\"varA\": \"222\"}"
	# parameters = "{\"varA\": \"222\"}"
	username = "root"
	output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
	output_cos_key_prefix = "log"
  }

`
