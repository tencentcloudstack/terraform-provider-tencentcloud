package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaSuperResolutionTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaSuperResolutionTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_super_resolution_template.media_super_resolution_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_super_resolution_template.media_super_resolution_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaSuperResolutionTemplate = `

resource "tencentcloud_ci_media_super_resolution_template" "media_super_resolution_template" {
  name = &lt;nil&gt;
  resolution = &lt;nil&gt;
  enable_scale_up = &lt;nil&gt;
  version = &lt;nil&gt;
}

`
