package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesSendEmailResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesSendEmail,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_send_email.send_email", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_send_email.send_email",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesSendEmail = `

resource "tencentcloud_ses_send_email" "send_email" {
  from_email_address = "noreply@mail.qcloud.com"
  destination = 
  subject = "test subject"
  reply_to_addresses = "reply@mail.qcloud.com"
  cc = 
  bcc = 
  template {
		template_i_d = 7000
		template_data = "{&quot;name&quot;:&quot;xxx&quot;,&quot;age&quot;:&quot;xx&quot;}"

  }
  simple {
		html = ""
		text = ""

  }
  attachments {
		file_name = "doc.zip"
		content = ""

  }
  unsubscribe = "1"
  trigger_type = 1
}

`
