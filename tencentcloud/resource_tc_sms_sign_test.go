package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSmsSign_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSmsSign,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sms_sign.sign", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sms_sign.sign",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSmsSign = `

resource "tencentcloud_sms_sign" "sign" {
  sign_name = "SignName"
  sign_type = 0
  document_type = 0
  international = 0
  sign_purpose = 0
  proof_image = ""
  commission_image = ""
  remark = ""
}

`
