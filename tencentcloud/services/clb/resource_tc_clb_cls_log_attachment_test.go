package clb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClbClsLogAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccClbClsLogAttachment,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_cls_log_attachment.clb_cls_log_attachment", "id")),
		}, {
			ResourceName:      "tencentcloud_clb_cls_log_attachment.clb_cls_log_attachment",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccClbClsLogAttachment = `

resource "tencentcloud_clb_cls_log_attachment" "clb_cls_log_attachment" {
}
`
