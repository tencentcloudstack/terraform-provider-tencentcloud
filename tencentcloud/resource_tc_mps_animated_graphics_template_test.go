package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMpsAnimatedGraphicsTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAnimatedGraphicsTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_animated_graphics_template.animated_graphics_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_animated_graphics_template.animated_graphics_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsAnimatedGraphicsTemplate = `

resource "tencentcloud_mps_animated_graphics_template" "animated_graphics_template" {
  fps = &lt;nil&gt;
  width = 0
  height = 0
  resolution_adaptive = "open"
  format = "gif"
  quality = 
  name = &lt;nil&gt;
  comment = &lt;nil&gt;
}

`
