package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseKeyPairAttachmentResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
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

const testAccLighthouseKeyPairAttachment = DefaultLighthoustVariables + `

resource "tencentcloud_lighthouse_key_pair_attachment" "key_pair_attachment" {
  key_id = "lhkp-d8zf3jmv"
  instance_id = var.lighthouse_id
}

`
