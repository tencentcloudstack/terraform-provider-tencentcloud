package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveCallbackTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveCallbackTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_callback_template.callback_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_callback_template.callback_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveCallbackTemplate = `

resource "tencentcloud_live_callback_template" "callback_template" {
  template_name = "demo"
  description = "this is demo"
  stream_begin_notify_url = "http://www.yourdomain.com/api/notify?action=streamBegin"
  stream_end_notify_url = "http://www.yourdomain.com/api/notify?action=streamEnd"
  record_notify_url = "http://www.yourdomain.com/api/notify?action=record"
  snapshot_notify_url = "http://www.yourdomain.com/api/notify?action=snapshot"
  porn_censorship_notify_url = "http://www.yourdomain.com/api/notify?action=porn"
  callback_key = "adasda131312"
  push_exception_notify_url = "http://www.yourdomain.com/api/notify?action=pushException"
}

`
