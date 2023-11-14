package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesReceiverResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesReceiver,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_receiver.receiver", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_receiver.receiver",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesReceiver = `

resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "Recipient group name"
  desc = "Recipient group description"
}

`
