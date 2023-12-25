package ses_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudSesSendEmailStatusDataSource_basic -v
func TestAccTencentCloudSesSendEmailStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-hongkong")
			tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesSendEmailStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ses_send_email_status.send_email_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.deliver_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.deliver_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.from_email_address"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.message_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.request_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.send_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.to_email_address"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.user_clicked"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.user_complainted"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.user_opened"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_send_email_status.send_email_status", "email_status_list.0.user_unsubscribed"),
				),
			},
		},
	})
}

const testAccSesSendEmailStatusDataSource = `

data "tencentcloud_ses_send_email_status" "send_email_status" {
  request_date = "2023-09-05"
  message_id = "qcloudses-30-1308919341-date-20230905112131-SGiMWY8csHwn1"
}

`
