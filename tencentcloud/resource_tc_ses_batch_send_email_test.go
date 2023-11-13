package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesBatchSendEmailResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesBatchSendEmail,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_batch_send_email.batch_send_email", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_batch_send_email.batch_send_email",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesBatchSendEmail = `

resource "tencentcloud_ses_batch_send_email" "batch_send_email" {
  from_email_address = "Tencent Cloud team &lt;noreply@mail.qcloud.com&gt;"
  receiver_id = 123
  subject = "test"
  task_type = 1
  reply_to_addresses = "reply@mail.qcloud.com"
  template {
		template_i_d = 5432
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
  cycle_param {
		begin_time = "2021-09-10 11:10:11"
		interval_time = 2
		term_cycle = 0

  }
  timed_param {
		begin_time = "2021-09-11 09:10:11"

  }
  unsubscribe = "1"
  a_d_location = 1
}

`
