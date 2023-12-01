package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssCallbackTemplateResource_basic -v
func TestAccTencentCloudCssCallbackTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssCallbackTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_callback_template.callback_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "callback_key", "adasda131312"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "description", "this is demo"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "porn_censorship_notify_url", "http://www.yourdomain.com/api/notify?action=porn"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "push_exception_notify_url", "http://www.yourdomain.com/api/notify?action=pushException"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "record_notify_url", "http://www.yourdomain.com/api/notify?action=record"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "snapshot_notify_url", "http://www.yourdomain.com/api/notify?action=snapshot"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "stream_begin_notify_url", "http://www.yourdomain.com/api/notify?action=streamBegin"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "stream_end_notify_url", "http://www.yourdomain.com/api/notify?action=streamEnd"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "template_name", "tf-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_callback_template.callback_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCssCallbackTemplateUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_callback_template.callback_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "callback_key", "adasda1313121"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "description", "this is demo1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "porn_censorship_notify_url", "http://www.yourdomain.com/api/notify?action=porn1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "push_exception_notify_url", "http://www.yourdomain.com/api/notify?action=pushException1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "record_notify_url", "http://www.yourdomain.com/api/notify?action=record1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "snapshot_notify_url", "http://www.yourdomain.com/api/notify?action=snapshot1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "stream_begin_notify_url", "http://www.yourdomain.com/api/notify?action=streamBegin1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "stream_end_notify_url", "http://www.yourdomain.com/api/notify?action=streamEnd1"),
					resource.TestCheckResourceAttr("tencentcloud_css_callback_template.callback_template", "template_name", "tf-test1"),
				),
			},
		},
	})
}

const testAccCssCallbackTemplate = `

resource "tencentcloud_css_callback_template" "callback_template" {
  template_name              = "tf-test"
  description                = "this is demo"
  stream_begin_notify_url    = "http://www.yourdomain.com/api/notify?action=streamBegin"
  stream_end_notify_url      = "http://www.yourdomain.com/api/notify?action=streamEnd"
  record_notify_url          = "http://www.yourdomain.com/api/notify?action=record"
  snapshot_notify_url        = "http://www.yourdomain.com/api/notify?action=snapshot"
  porn_censorship_notify_url = "http://www.yourdomain.com/api/notify?action=porn"
  callback_key               = "adasda131312"
  push_exception_notify_url  = "http://www.yourdomain.com/api/notify?action=pushException"
}

`

const testAccCssCallbackTemplateUp = `

resource "tencentcloud_css_callback_template" "callback_template" {
  template_name              = "tf-test1"
  description                = "this is demo1"
  stream_begin_notify_url    = "http://www.yourdomain.com/api/notify?action=streamBegin1"
  stream_end_notify_url      = "http://www.yourdomain.com/api/notify?action=streamEnd1"
  record_notify_url          = "http://www.yourdomain.com/api/notify?action=record1"
  snapshot_notify_url        = "http://www.yourdomain.com/api/notify?action=snapshot1"
  porn_censorship_notify_url = "http://www.yourdomain.com/api/notify?action=porn1"
  callback_key               = "adasda1313121"
  push_exception_notify_url  = "http://www.yourdomain.com/api/notify?action=pushException1"
}

`
