package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsAnimatedGraphicsTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAnimatedGraphicsTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_animated_graphics_templates.animated_graphics_templates")),
			},
		},
	})
}

const testAccMpsAnimatedGraphicsTemplatesDataSource = `

data "tencentcloud_mps_animated_graphics_templates" "animated_graphics_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  animated_graphics_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		format = &lt;nil&gt;
		fps = &lt;nil&gt;
		quality = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

  }
}

`
