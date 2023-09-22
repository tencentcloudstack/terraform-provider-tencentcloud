package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentNeedFixCloudSesBlackListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-hongkong")
			testAccPreCheckBusiness(t, ACCOUNT_TYPE_SES)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesBlackList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ses_black_list_delete.black_list", "id"),
				),
			},
		},
	})
}

const testAccSesBlackList = `

resource "tencentcloud_ses_send_email" "send_email" {
  from_email_address = "aaa@iac-tf.cloud"
  destination        = ["terraform-tf@gmail.com"]
  subject            = "test subject"
  reply_to_addresses = "aaa@iac-tf.cloud"

  template {
    template_id   = 99629
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  unsubscribe  = "1"
  trigger_type = 1
}

resource "tencentcloud_ses_black_list_delete" "black_list" {
  email_address = "terraform-tf@gmail.com"
  depends_on = [ tencentcloud_ses_send_email.send_email ]
}

`
