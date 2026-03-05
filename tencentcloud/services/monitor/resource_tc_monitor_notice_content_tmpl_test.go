package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMonitorNoticeContentTmplResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorNoticeContentTmpl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_notice_content_tmpl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_notice_content_tmpl.example", "tmpl_name", "terraform-test-template"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_notice_content_tmpl.example", "monitor_type", "MT_QCE"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_notice_content_tmpl.example", "tmpl_language", "zh"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_notice_content_tmpl.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorNoticeContentTmplUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_notice_content_tmpl.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_notice_content_tmpl.example", "tmpl_name", "terraform-test-template-updated"),
				),
			},
		},
	})
}

const testAccMonitorNoticeContentTmpl = `
resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "terraform-test-template"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"
  tmpl_contents = jsonencode({
    "QCloudYehe" : [
      {
        "MatchingStatus" : ["Trigger"],
        "Template" : {
          "Email" : {
            "ContentTmpl" : base64encode("告警通知内容"),
            "TitleTmpl" : base64encode("告警通知标题")
          }
        }
      }
    ]
  })
}
`

const testAccMonitorNoticeContentTmplUpdate = `
resource "tencentcloud_monitor_notice_content_tmpl" "example" {
  tmpl_name     = "terraform-test-template-updated"
  monitor_type  = "MT_QCE"
  tmpl_language = "zh"
  tmpl_contents = jsonencode({
    "QCloudYehe" : [
      {
        "MatchingStatus" : ["Trigger"],
        "Template" : {
          "Email" : {
            "ContentTmpl" : base64encode("更新后的告警通知内容"),
            "TitleTmpl" : base64encode("更新后的告警通知标题")
          }
        }
      }
    ]
  })
}
`
