package ses_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixSesBatchSendEmailResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  from_email_address = "aaa@iac-tf.cloud"
  receiver_id        = 1063742
  subject            = "terraform test"
  task_type          = 1
  reply_to_addresses = "reply@mail.qcloud.com"
  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"

  }

  cycle_param {
    begin_time = "2023-09-07 15:10:00"
    interval_time = 1
  }
  timed_param {
    begin_time = "2023-09-07 15:20:00"
  }
  unsubscribe = "0"
  ad_location = 0
}

`
