package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCiMediaTtsTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaTtsTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_tts_template.media_tts_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_tts_template.media_tts_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaTtsTemplate = `

resource "tencentcloud_ci_media_tts_template" "media_tts_template" {
	bucket = "terraform-ci-1308919341"
	name = "tts_template"
	mode = "Asyc"
	codec = "pcm"
	voice_type = "ruxue"
	volume = "0"
	speed = "100"
  }

`
