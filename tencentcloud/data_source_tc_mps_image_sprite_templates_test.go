package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsImageSpriteTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsImageSpriteTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_image_sprite_templates.image_sprite_templates")),
			},
		},
	})
}

const testAccMpsImageSpriteTemplatesDataSource = `

data "tencentcloud_mps_image_sprite_templates" "image_sprite_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  image_sprite_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		sample_type = &lt;nil&gt;
		sample_interval = &lt;nil&gt;
		row_count = &lt;nil&gt;
		column_count = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		fill_type = &lt;nil&gt;
		comment = &lt;nil&gt;
		format = &lt;nil&gt;

  }
}

`
