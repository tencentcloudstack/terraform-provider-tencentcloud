package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsTranscodeTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_transcode_template.transcode_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_transcode_template.transcode_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsTranscodeTemplate = `

resource "tencentcloud_mps_transcode_template" "transcode_template" {
  container    = "mp4"
  name         = "tf_transcode_template"
  remove_audio = 0
  remove_video = 0

  audio_template {
    audio_channel = 2
    bitrate       = 27
    codec         = "libfdk_aac"
    sample_rate   = 32000
  }

  video_template {
    bitrate             = 130
    codec               = "libx264"
    fill_type           = "black"
    fps                 = 20
    gop                 = 0
    height              = 4096
    resolution_adaptive = "close"
    vcrf                = 0
    width               = 128
  }
}

`
