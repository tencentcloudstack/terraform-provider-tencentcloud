package cam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCamMessageReceiverResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamMessageReceiver,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cam_message_receiver.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_message_receiver.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_message_receiver.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_message_receiver.example", "country_code"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_message_receiver.example", "phone_number"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_message_receiver.example", "email"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_message_receiver.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamMessageReceiver = `
resource "tencentcloud_cam_message_receiver" "example" {
  name         = "tf-example"
  remark       = "remark."
  country_code = "86"
  phone_number = "18123456789"
  email        = "demo@qq.com"

  lifecycle {
    ignore_changes = [ email, phone_number ]
  }
}
`
