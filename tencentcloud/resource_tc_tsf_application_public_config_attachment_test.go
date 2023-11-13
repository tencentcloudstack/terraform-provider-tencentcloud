package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  config_id = "dcfg-p-123456"
  namespace_id = "namespace-123456"
  release_desc = "product version"
}

`
