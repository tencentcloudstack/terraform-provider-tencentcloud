package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbTargetGroupAttachmentsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachments,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "id")),
			},
		},
	})
}

const testAccClbTargetGroupAttachments = `

resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  load_balancer_id = "lb-phbx2420"
  associations {
		listener_id = "lbl-m2q6sp9m"
		target_group_id = "lbtg-5xunivs0"
		location_id = "loc-jjqr0ric"
  }
}

`
