package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  name = &lt;nil&gt;
  mode = &lt;nil&gt;
  codec = &lt;nil&gt;
  voice_type = &lt;nil&gt;
  volume = &lt;nil&gt;
  speed = &lt;nil&gt;
}

`
