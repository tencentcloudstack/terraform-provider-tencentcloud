package tat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvocationInvokeAttachmentResource_basic -v
func TestAccTencentCloudTatInvocationInvokeAttachmentResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvocationInvokeAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "instance_id", tcacctest.DefaultInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "working_directory", "/root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "timeout", "100"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "username", "root"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "output_cos_bucket_url", "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "output_cos_key_prefix", "log"),
					resource.TestCheckResourceAttr("tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment", "command_id", tcacctest.DefaultCommandId),
				),
			},
			{
				ResourceName:      "tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTatInvocationInvokeAttachmentVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultInstanceId + `"
}

variable "command_id" {
	default = "` + tcacctest.DefaultCommandId + `"
}
`
const testAccTatInvocationInvokeAttachment = testAccTatInvocationInvokeAttachmentVar + `

resource "tencentcloud_tat_invocation_invoke_attachment" "invocation_invoke_attachment" {
	instance_id = var.instance_id
	working_directory = "/root"
	timeout = 100
	# parameters = "{\"varA\": \"222\"}"
	username = "root"
	output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
	output_cos_key_prefix = "log"
	command_id = var.command_id
  }

`
