package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudVodImageSpriteTemplates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodImageSpriteTemplates,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vod_image_sprite_templates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.sample_type", "Percent"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.sample_interval", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.row_count", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.column_count", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.name", "tf-sprite"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.comment", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.fill_type", "stretch"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.width", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.height", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.resolution_adaptive", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.update_time"),
				),
			},
		},
	})
}

const testAccVodImageSpriteTemplates = testAccVodImageSpriteTemplate + `
data "tencentcloud_vod_image_sprite_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_image_sprite_template.foo.id
}
`
