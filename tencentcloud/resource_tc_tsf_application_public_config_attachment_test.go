package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationPublicConfigAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationPublicConfigAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_application_public_config_attachment.application_public_config_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_application_public_config_attachment.application_public_config_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfApplicationPublicConfigAttachment = `

resource "tencentcloud_tsf_application_public_config_attachment" "application_public_config_attachment" {
  config_id = ""
  namespace_id = ""
  release_desc = ""
                    }

`
