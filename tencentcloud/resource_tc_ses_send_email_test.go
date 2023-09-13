package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixSesSendEmailResource_basic(t *testing.T) {
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
  from_email_address = "aaa@iac-tf.cloud"
  destination        = ["1055482519@qq.com"]
  subject            = "test subject"
  reply_to_addresses = "aaa@iac-tf.cloud"

  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  unsubscribe  = "1"
  trigger_type = 1
}

`
