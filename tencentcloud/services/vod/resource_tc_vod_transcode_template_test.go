package vod_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

const vodTranscodeTemplateKey = "tencentcloud_vod_transcode_template.transcode_template"

func TestAccTencentCloudVodTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodTranscodeTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vodTranscodeTemplateKey, "id"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "container", "mp4"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "name", "720pTranscodeTemplate"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "comment", "test transcode mp4 720p"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "remove_video", "0"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "remove_audio", "0"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.codec", "libx264"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.fps", "26"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.bitrate", "1000"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.resolution_adaptive", "open"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.width", "0"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.height", "720"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.fill_type", "stretch"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.vcrf", "1"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.gop", "250"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.preserve_hdr_switch", "OFF"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "video_template.0.codec_tag", "hvc1"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "audio_template.0.codec", "libfdk_aac"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "audio_template.0.bitrate", "128"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "audio_template.0.sample_rate", "44100"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "audio_template.0.audio_channel", "2"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "segment_type", "ts"),
				),
			},
			{
				Config: testAccVodTranscodeTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vodTranscodeTemplateKey, "id"),
					resource.TestCheckResourceAttr(vodTranscodeTemplateKey, "comment", "test transcode mp4 720p update"),
				),
			},
			{
				ResourceName:      vodTranscodeTemplateKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVodTranscodeTemplate = `
resource  "tencentcloud_vod_sub_application" "transcode_template_sub_application" {
	name = "transcodeTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container = "mp4"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.transcode_template_sub_application.id)[1])
  name = "720pTranscodeTemplate"
  comment = "test transcode mp4 720p"
  remove_video = 0
  remove_audio = 0
  video_template {
	codec = "libx264"
	fps = 26
	bitrate = 1000
	resolution_adaptive = "open"
	width = 0
	height = 720
	fill_type = "stretch"
	vcrf = 1
	gop = 250
	preserve_hdr_switch = "OFF"
	codec_tag = "hvc1"

  }
  audio_template {
	codec = "libfdk_aac"
	bitrate = 128
	sample_rate = 44100
	audio_channel = 2
  }
  segment_type = "ts"
}
`

const testAccVodTranscodeTemplateUpdate = `
resource  "tencentcloud_vod_sub_application" "transcode_template_sub_application" {
	name = "transcodeTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container = "mp4"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.transcode_template_sub_application.id)[1])
  name = "720pTranscodeTemplate"
  comment = "test transcode mp4 720p update"
  remove_video = 0
  remove_audio = 0
  video_template {
	codec = "libx264"
	fps = 26
	bitrate = 1000
	resolution_adaptive = "open"
	width = 0
	height = 720
	fill_type = "stretch"
	vcrf = 1
	gop = 250
	preserve_hdr_switch = "OFF"
	codec_tag = "hvc1"

  }
  audio_template {
	codec = "libfdk_aac"
	bitrate = 128
	sample_rate = 44100
	audio_channel = 2
  }
  segment_type = "ts"
}
`
