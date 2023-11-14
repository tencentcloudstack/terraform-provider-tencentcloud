package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveTimeshiftTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveTimeshiftTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_timeshift_template.timeshift_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_timeshift_template.timeshift_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveTimeshiftTemplate = `

resource "tencentcloud_live_timeshift_template" "timeshift_template" {
  template_name = ""
  duration = 
  description = ""
  area = ""
  item_duration = 
  remove_watermark = 
  transcode_template_ids = 
}

`
