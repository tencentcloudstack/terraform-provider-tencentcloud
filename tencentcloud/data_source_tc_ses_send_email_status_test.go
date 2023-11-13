package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesSendEmailStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesSendEmailStatusDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_send_email_status.send_email_status")),
			},
		},
	})
}

const testAccSesSendEmailStatusDataSource = `

data "tencentcloud_ses_send_email_status" "send_email_status" {
  request_date = "2020-09-22"
  message_id = "qcloudses-30-4123414323-date-20210101094334-syNARhMTbKI1"
  to_email_address = "example@cloud.com"
  }

`
