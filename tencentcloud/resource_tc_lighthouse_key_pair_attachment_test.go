package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseKeyPairAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseKeyPairAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_key_pair_attachment.key_pair_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_key_pair_attachment.key_pair_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseKeyPairAttachment = `

resource "tencentcloud_lighthouse_key_pair_attachment" "key_pair_attachment" {
  key_ids = 
  instance_ids = 
}

`
