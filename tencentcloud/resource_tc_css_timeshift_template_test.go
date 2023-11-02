package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssTimeshiftTemplateResource_basic -v
func TestAccTencentCloudCssTimeshiftTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssTimeshiftTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_timeshift_template.timeshift_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "area", "Mainland"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "description", "timeshift template"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "duration", "604800"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "item_duration", "5"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "remove_watermark", "true"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "template_name", "tf-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_timeshift_template.timeshift_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCssTimeshiftTemplateUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_timeshift_template.timeshift_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "area", "Overseas"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "description", "timeshift template1"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "duration", "259200"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "item_duration", "6"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "remove_watermark", "false"),
					resource.TestCheckResourceAttr("tencentcloud_css_timeshift_template.timeshift_template", "template_name", "tf-test1"),
				),
			},
		},
	})
}

const testAccCssTimeshiftTemplate = `

resource "tencentcloud_css_timeshift_template" "timeshift_template" {
    area                   = "Mainland"
    description            = "timeshift template"
    duration               = 604800
    item_duration          = 5
    remove_watermark       = true
    template_name          = "tf-test"
}

`

const testAccCssTimeshiftTemplateUp = `

resource "tencentcloud_css_timeshift_template" "timeshift_template" {
    area                   = "Overseas"
    description            = "timeshift template1"
    duration               = 259200
    item_duration          = 6
    remove_watermark       = false
    template_name          = "tf-test1"
}

`
