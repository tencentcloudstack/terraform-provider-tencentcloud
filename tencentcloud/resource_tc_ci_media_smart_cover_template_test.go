package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCiMediaSmartCoverTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaSmartCoverTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_smart_cover_template.media_smart_cover_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_smart_cover_template.media_smart_cover_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaSmartCoverTemplate = `

resource "tencentcloud_ci_media_smart_cover_template" "media_smart_cover_template" {
	bucket = "terraform-ci-1308919341"
	name = "smart_cover_template"
	smart_cover {
		  format = "jpg"
		  width = "1280"
		  height = "960"
		  count = "10"
		  delete_duplicates = "true"
	}
  }

`
