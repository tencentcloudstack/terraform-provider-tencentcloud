package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesReceiverDetailResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesReceiverDetail,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_receiver_detail.receiver_detail", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_receiver_detail.receiver_detail",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesReceiverDetail = `

resource "tencentcloud_ses_receiver_detail" "receiver_detail" {
  receiver_id = 123
  datas {
		email = "abc@ef.com"
		template_data = "{&quot;name&quot;:&quot;xxx&quot;,&quot;age&quot;:&quot;xx&quot;}"

  }
  emails = 
}

`
