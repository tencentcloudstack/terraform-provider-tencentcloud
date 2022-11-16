package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCSSLiveTranscodeTemplate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveTranscodeTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_template.live_transcode_template", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_live_transcode_template.liveTranscodeTemplate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssLiveTranscodeTemplate = `

resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name = ""
  video_bitrate = ""
  acodec = ""
  audio_bitrate = ""
  vcodec = ""
  description = ""
  need_video = ""
  width = ""
  need_audio = ""
  height = ""
  fps = ""
  gop = ""
  rotate = ""
  profile = ""
  bitrate_to_orig = ""
  height_to_orig = ""
  fps_to_orig = ""
  ai_trans_code = ""
  adapt_bitrate_percent = ""
  short_edge_as_height = ""
  d_r_m_type = ""
  d_r_m_tracks = ""
}

`
