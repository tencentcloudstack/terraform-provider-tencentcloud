package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudVodAdaptiveDynamicStreamingTemplates(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodAdaptiveDynamicStreamingTemplates,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.format", "HLS"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.name", "tf-adaptive"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.drm_type", "SimpleAES"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.disable_higher_video_bitrate", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.disable_higher_video_resolution", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.comment", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.video.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.video.0.codec", "libx264"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.video.0.fps", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.video.0.bitrate", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.audio.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.audio.0.codec", "libfdk_aac"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.audio.0.bitrate", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.audio.0.sample_rate", "32000"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.0.remove_audio", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.video.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.video.0.codec", "libx264"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.video.0.fps", "4"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.video.0.bitrate", "256"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.audio.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.audio.0.codec", "libfdk_aac"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.audio.0.bitrate", "256"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.audio.0.sample_rate", "44100"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.stream_info.1.remove_audio", "true"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_adaptive_dynamic_streaming_templates.foo", "template_list.0.update_time"),
				),
			},
		},
	})
}

const testAccVodAdaptiveDynamicStreamingTemplates = testAccVodAdaptiveDynamicStreamingTemplate + `
data "tencentcloud_vod_adaptive_dynamic_streaming_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
}
`
