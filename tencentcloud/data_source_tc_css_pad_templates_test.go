package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssPadTemplatesDataSource_basic -v
func TestAccTencentCloudCssPadTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPadTemplatesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_pad_templates.pad_templates"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.max_duration"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.template_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.template_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_css_pad_templates.pad_templates", "templates.0.wait_duration"),
				),
			},
		},
	})
}

const testAccCssPadTemplatesDataSource = testAccCssPadTemplate + `

data "tencentcloud_css_pad_templates" "pad_templates" {
	depends_on = [ tencentcloud_css_pad_template.pad_template ]
}

`
