package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsAdaptiveDynamicStreamingTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAdaptiveDynamicStreamingTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template", "id")),
			},
			{
				Config: testAccMpsAdaptiveDynamicStreamingTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template", "name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsAdaptiveDynamicStreamingTemplate = `

resource "tencentcloud_mps_adaptive_dynamic_streaming_template" "adaptive_dynamic_streaming_template" {
  comment                         = "terrraform test"
  disable_higher_video_bitrate    = 0
  disable_higher_video_resolution = 1
  format                          = "HLS"
  name                            = "terrraform-test"

  stream_infos {
    remove_audio = 0
    remove_video = 0

    audio {
      audio_channel = 1
      bitrate       = 55
      codec         = "libmp3lame"
      sample_rate   = 32000
    }

    video {
      bitrate             = 245
      codec               = "libx264"
      fill_type           = "black"
      fps                 = 30
      gop                 = 0
      height              = 135
      resolution_adaptive = "open"
      vcrf                = 0
      width               = 145
    }
  }
  stream_infos {
    remove_audio = 0
    remove_video = 0

    audio {
      audio_channel = 2
      bitrate       = 60
      codec         = "libfdk_aac"
      sample_rate   = 32000
    }

    video {
      bitrate             = 400
      codec               = "libx264"
      fill_type           = "black"
      fps                 = 40
      gop                 = 0
      height              = 150
      resolution_adaptive = "open"
      vcrf                = 0
      width               = 160
    }
  }
}


`

const testAccMpsAdaptiveDynamicStreamingTemplateUpdate = `

resource "tencentcloud_mps_adaptive_dynamic_streaming_template" "adaptive_dynamic_streaming_template" {
  comment                         = "terrraform test"
  disable_higher_video_bitrate    = 0
  disable_higher_video_resolution = 1
  format                          = "HLS"
  name                            = "terraform-for-test"

  stream_infos {
    remove_audio = 0
    remove_video = 0

    audio {
      audio_channel = 1
      bitrate       = 55
      codec         = "libmp3lame"
      sample_rate   = 32000
    }

    video {
      bitrate             = 245
      codec               = "libx264"
      fill_type           = "black"
      fps                 = 30
      gop                 = 0
      height              = 130
      resolution_adaptive = "open"
      vcrf                = 0
      width               = 140
    }
  }
  stream_infos {
    remove_audio = 0
    remove_video = 0

    audio {
      audio_channel = 2
      bitrate       = 60
      codec         = "libfdk_aac"
      sample_rate   = 32000
    }

    video {
      bitrate             = 400
      codec               = "libx264"
      fill_type           = "black"
      fps                 = 40
      gop                 = 0
      height              = 150
      resolution_adaptive = "open"
      vcrf                = 0
      width               = 160
    }
  }
}


`
