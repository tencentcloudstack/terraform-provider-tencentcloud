package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationFileConfigAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationFileConfigAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_file_config_attachment.application_file_config_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_file_config_attachment.application_file_config_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationFileConfigAttachment = `

resource "tencentcloud_tsf_application_file_config_attachment" "application_file_config_attachment" {
  config_id = ""
  group_id = ""
  release_desc = ""
                  }

`
