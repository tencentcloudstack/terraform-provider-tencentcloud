package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSesEmail_address_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps:     []resource.TestStep{
			//{
			//	Config: testAccSesEmail_address,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet("tencentcloud_ses_email_address.email_address", "id"),
			//	),
			//},
			//{
			//	ResourceName:      "tencentcloud_ses_email_address.email_address",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

const testAccSesEmail_address = `

resource "tencentcloud_ses_email_address" "email_address" {
  email_address = ""
  email_sender_name = ""
}

`
