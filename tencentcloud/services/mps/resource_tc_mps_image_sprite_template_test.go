package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsImageSpriteTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsImageSpriteTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_image_sprite_template.image_sprite_template", "id")),
			},
			{
				Config: testAccMpsImageSpriteTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_image_sprite_template.image_sprite_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_image_sprite_template.image_sprite_template", "name", "terraform-for-test"),
				),
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
  column_count        = 10
  fill_type           = "stretch"
  format              = "jpg"
  height              = 143
  name                = "terraform-test"
  resolution_adaptive = "open"
  row_count           = 10
  sample_interval     = 10
  sample_type         = "Time"
  width               = 182
}

`

const testAccMpsImageSpriteTemplateUpdate = `

resource "tencentcloud_mps_image_sprite_template" "image_sprite_template" {
  column_count        = 10
  fill_type           = "stretch"
  format              = "jpg"
  height              = 143
  name                = "terraform-for-test"
  resolution_adaptive = "open"
  row_count           = 10
  sample_interval     = 10
  sample_type         = "Time"
  width               = 182
}

`
