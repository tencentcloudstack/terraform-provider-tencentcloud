package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVodImageSpriteTemplateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVodImageSpriteTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodImageSpriteTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodImageSpriteTemplateExists("tencentcloud_vod_image_sprite_template.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "sample_type", "Percent"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "sample_interval", "10"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "row_count", "3"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "column_count", "3"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "name", "tf-sprite"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "comment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "fill_type", "stretch"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "width", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "height", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "resolution_adaptive", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_image_sprite_template.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_image_sprite_template.foo", "update_time"),
				),
			},
			{
				Config: testAccVodImageSpriteTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "sample_type", "Time"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "sample_interval", "11"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "row_count", "4"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "column_count", "4"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "name", "tf-sprite-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "comment", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "fill_type", "black"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "width", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "height", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_image_sprite_template.foo", "resolution_adaptive", "true"),
				),
			},
			{
				ResourceName:            "tencentcloud_vod_image_sprite_template.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sub_app_id"},
			},
		},
	})
}

func testAccCheckVodImageSpriteTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	vodService := VodService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_image_sprite_template" {
			continue
		}
		var (
			filter = map[string]interface{}{
				"definitions": []string{rs.Primary.ID},
			}
		)

		templates, err := vodService.DescribeImageSpriteTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 {
			return nil
		}
		return fmt.Errorf("vod image sprite template still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVodImageSpriteTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod image sprite template %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod image sprite template id is not set")
		}
		vodService := VodService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		var (
			filter = map[string]interface{}{
				"definitions": []string{rs.Primary.ID},
			}
		)
		templates, err := vodService.DescribeImageSpriteTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 || len(templates) != 1 {
			return fmt.Errorf("vod image sprite template doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodImageSpriteTemplate = `
resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}
`

const testAccVodImageSpriteTemplateUpdate = `
resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Time"
  sample_interval     = 11
  row_count           = 4
  column_count        = 4
  name                = "tf-sprite-update"
  comment             = "test-update"
  fill_type           = "black"
  width               = 129
  height              = 129
  resolution_adaptive = true
}
`
