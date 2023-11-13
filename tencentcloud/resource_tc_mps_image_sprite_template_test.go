package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsImageSpriteTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsImageSpriteTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_image_sprite_template.image_sprite_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_image_sprite_template.image_sprite_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsImageSpriteTemplate = `

resource "tencentcloud_mps_image_sprite_template" "image_sprite_template" {
  sample_type = &lt;nil&gt;
  sample_interval = &lt;nil&gt;
  row_count = &lt;nil&gt;
  column_count = &lt;nil&gt;
  name = &lt;nil&gt;
  width = 0
  height = 0
  resolution_adaptive = "open"
  fill_type = "black"
  comment = &lt;nil&gt;
  format = "jpg"
}

`
