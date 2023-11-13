package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  name = &lt;nil&gt;
  smart_cover {
		format = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		count = &lt;nil&gt;
		delete_duplicates = &lt;nil&gt;

  }
}

`
