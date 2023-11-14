package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesEmailAddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesEmailAddress,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_email_address.email_address", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_email_address.email_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesEmailAddress = `

resource "tencentcloud_ses_email_address" "email_address" {
  email_address = &lt;nil&gt;
  email_sender_name = &lt;nil&gt;
}

`
