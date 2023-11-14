package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssLiveTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssLiveTranscodeTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_template.live_transcode_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_live_transcode_template.live_transcode_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssLiveTranscodeTemplate = `

resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name = &lt;nil&gt;
  video_bitrate = &lt;nil&gt;
  acodec = &lt;nil&gt;
  audio_bitrate = &lt;nil&gt;
  vcodec = &lt;nil&gt;
  description = &lt;nil&gt;
  need_video = &lt;nil&gt;
  width = &lt;nil&gt;
  need_audio = &lt;nil&gt;
  height = &lt;nil&gt;
  fps = &lt;nil&gt;
  gop = &lt;nil&gt;
  rotate = &lt;nil&gt;
  profile = &lt;nil&gt;
  bitrate_to_orig = &lt;nil&gt;
  height_to_orig = &lt;nil&gt;
  fps_to_orig = &lt;nil&gt;
  ai_trans_code = &lt;nil&gt;
  adapt_bitrate_percent = &lt;nil&gt;
  short_edge_as_height = &lt;nil&gt;
  d_r_m_type = &lt;nil&gt;
  d_r_m_tracks = &lt;nil&gt;
}

`
