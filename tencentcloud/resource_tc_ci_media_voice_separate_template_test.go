package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaVoiceSeparateTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaVoiceSeparateTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_voice_separate_template.media_voice_separate_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaVoiceSeparateTemplate = `

resource "tencentcloud_ci_media_voice_separate_template" "media_voice_separate_template" {
  name = &lt;nil&gt;
  audio_mode = &lt;nil&gt;
  audio_config {
		codec = &lt;nil&gt;
		samplerate = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		channels = &lt;nil&gt;

  }
}

`
