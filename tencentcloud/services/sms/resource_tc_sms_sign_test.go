package sms_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSmsSign_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SMS) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSmsSign,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sms_sign.sign", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sms_sign.sign", "sign_name", "terraform"),
				),
			},
		},
	})
}

const testAccSmsSign = `

resource "tencentcloud_sms_sign" "sign" {
  sign_name     = "terraform"
  sign_type     = 1
  document_type = 4
  international = 0
  sign_purpose  = 0
  proof_image = "dGhpcyBpcyBhIGV4YW1wbGU="
}

`
